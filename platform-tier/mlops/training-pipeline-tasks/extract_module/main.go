package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	// "cloud.google.com/go/storage"
	// "google.golang.org/api/option"
)

func lookupStringEnv(envVar string, defaultValue string) string {
	envVarValue, ok := os.LookupEnv(envVar)
	if !ok {
		return defaultValue
	}
	fmt.Println("Env var ", envVar, " = ", envVarValue)
	return envVarValue
}

func lookupInt64Env(envVar string, defaultValue int64) int64 {
	envVarValue, ok := os.LookupEnv(envVar)
	if !ok {
		return defaultValue
	}
	int64Val, _ := strconv.ParseInt(envVarValue, 10, 64)
	return int64Val
}

func SendPostRequest(url string, filename string, filetype string) []byte {

	fmt.Printf("uploading data file %s to url %s\n", filename, url)

	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(filetype, filepath.Base(file.Name()))

	if err != nil {
		log.Fatal(err)
	}

	io.Copy(part, file)
	writer.Close()
	request, err := http.NewRequest("POST", url, body)

	if err != nil {
		log.Fatal(err)
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("uploading data file %s to url %s complete !!! \n", filename, url)
	return content
}

func downloadFileFromGCPStorage(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

type Pre struct {
	XMLName xml.Name `xml:"pre"`
	Text    string   `xml:",chardata"`
	A       []struct {
		Text string `xml:",chardata"`
		Href string `xml:"href,attr"`
	} `xml:"a"`
}

// func downloadRawData(url string) [][]string {
func downloadRawData(url string, fileName string) [][]string {

	caCert, err := ioutil.ReadFile("/keys/ssh-publickey")
    if err != nil {
        log.Fatal(err)
    }
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

	client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                RootCAs:      caCertPool,
				InsecureSkipVerify: true,
            },
        },
    }	
	
	httpGet := func(url_param string) (data []byte, status bool) {
		resp, err := client.Get(url_param)
		if err != nil {
			return nil, false
		}
		defer resp.Body.Close()

		// Check server response
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("bad status: %s", resp.Status)
			return nil, false
		}

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
			return nil, false
		}

		return data, true
	}

	readCsvFile := func(data []byte) [][]string {

		csvReader := csv.NewReader(bytes.NewReader(data))
		records, err := csvReader.ReadAll()
		if err != nil {
			fmt.Printf("Unable to parse file as CSV for %s", err)
		}

		return records
	}

	// var p Pre
	// if err := xml.Unmarshal(httpGet(url), &p); err != nil {
	// 	fmt.Printf("Unable to parse file as CSV for %s", err)
	// }

	var rawDataRows [][]string

	// for _, v := range p.A {
	// file := v.Text
	b, status := httpGet(url + fileName)
	// fmt.Printf(string(b))
	if status {
		records := readCsvFile(b)

		for _, e := range records[1:] {
			rawDataRows = append(rawDataRows, e)
		}
	} else {
		return nil
	}
	// }

	return rawDataRows
}

func labelData(rawData []string) []string {
	// just for the demo , these are fake labels.

	r := 0
	n := 1000000
	if rand.Intn(n) > n/2 {
		r = 0
	} else {
		r = 1
	}
	return append(rawData, fmt.Sprintf("%d", r))
}

func normalizeData(rawData [][]string) []string {

	var normalizedData []string

	for _, v := range rawData {
		v = labelData(v)
		row := strings.Join(v, ",")
		normalizedData = append(normalizedData, row)
	}

	return normalizedData
}

func uploadToModelRegistry(url string, data []string) []byte {

	writeToFile := func() string {
		t := time.Now()
		fileName := fmt.Sprintf("normalized_training_data_%d-%02d-%02dT%02d:%02d:%02d:%02d.csv",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second(), t.UnixNano()/int64(time.Millisecond))

		file, err := os.Create(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		for _, v := range data {
			file.WriteString(v + "\n")
		}

		return fileName
	}

	fileName := writeToFile()
	filetype := "file"
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fmt.Printf("uploading data file %s to url %s\n", fileName, url)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(filetype, filepath.Base(file.Name()))

	if err != nil {
		log.Fatal(err)
	}

	io.Copy(part, file)
	writer.Close()
	request, err := http.NewRequest("POST", url, body)

	if err != nil {
		log.Fatal(err)
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
	caCert, err := ioutil.ReadFile("/keys/ssh-publickey")
    if err != nil {
        log.Fatal(err)
    }
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

	client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                RootCAs:      caCertPool,
				InsecureSkipVerify: true,
            },
        },
    }	
	
	response, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("uploading data file %s to url %s complete !!! \n", fileName, url)
	return content
}

type ControlMessage struct {
	Command string `json:"command"`
	Payload string `json:"payload"`
}

func listenForControlMessage(kafka_reader *kafka.Reader, ctx context.Context, topic string) string {

	var controlMsg ControlMessage

	for {

		fmt.Printf("Waiting for kafka mesages 'extract-data' from broker %s on topic %s\n", topic)
		msg, err := kafka_reader.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}

		json.Unmarshal([]byte(msg.Value), &controlMsg)
		if controlMsg.Command == "extract-data" {
			break
		}
	}

	return controlMsg.Payload

}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	topic := lookupStringEnv("CONTROL-TOPIC", "control-message")
	brokerAddress := lookupStringEnv("KAFKA-BROKER", "35.236.22.237:32199")

	kafka_reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
	})
	ctx := context.Background()

	downloadUrl := lookupStringEnv("RAW_TRAINING_DATA_DOWNLOAD_REGISTRY_URL", "http://localhost:8081/")
	uploadUrl := lookupStringEnv("NORMALIZED_DATA_UPLOAD_REGISTRY_URL", "http://localhost:8080/uploadNormalizedData")

	trainingFileName := listenForControlMessage(kafka_reader, ctx, topic)
	fmt.Println(trainingFileName)

	rawDataRows := downloadRawData(downloadUrl, trainingFileName)
	if rawDataRows != nil {
		normalizedData := normalizeData(rawDataRows)
		uploadToModelRegistry(uploadUrl, normalizedData)

		fmt.Println(rawDataRows[0])
	}

}

// find which port is in use
//sudo lsof -i -P -n | grep LISTEN

// fmt.Println("Starting go app v4...")

// dt := time.Now()
// filename := dt.Format("01-02-2006 15:04:05") + ".csv"

// fileToDownload := "http://storage.googleapis.com/architectsguide2aiot-aiot-mlops-demo/agglomeration-tower1-cframe-shaded-pole_solvent_motor.csv"
// e := downloadFileFromGCPStorage(filename, fileToDownload)

// if e != nil {
// 	log.Fatal(e)
// }

// file, err := os.Open(filename)
// if err != nil {
// 	log.Fatal(err)
// }
// defer file.Close()

// scanner := bufio.NewScanner(file)
// for scanner.Scan() {             // internally, it advances token based on sperator
//     fmt.Println(scanner.Text())  // token in unicode-char
//     fmt.Println(scanner.Bytes()) // token in bytes

// }
// SendPostRequest("http://localhost:8080/uploadTrainingData", "/home/agoja/kubecon-2021-aiot-demo/apps/logistic_regression_model/training_data/training_dataset_1", "file")
// SendPostRequest("http://35.236.22.237:30007/uploadTrainingData", filename , "file")
