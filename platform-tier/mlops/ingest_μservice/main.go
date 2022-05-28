package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
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
	fmt.Println("Env var ", envVar, " = ", envVarValue)
	int64Val, _ := strconv.ParseInt(envVarValue, 10, 64)
	return int64Val
}

type RawSensorData struct {
	DeviceID    string  `json:"deviceID"`
	TimeStamp   string  `json:"timeStamp"`
	Current     float32 `json:"current"`
	Temperature float32 `json:"temperature"`
	Vibration   float32 `json:"vibration"`
	Sound       float32 `json:"sound"`
}

func main() {

	data_topic := lookupStringEnv("DATA-TOPIC", "shaded-pole-motor-sensor_data")
	control_topic := lookupStringEnv("CONTROL-TOPIC", "control-message")
	brokerAddress := lookupStringEnv("KAFKA-BROKER", "35.236.22.237:32199")
	upload_url := lookupStringEnv("TRAINING_DATA_UPLOAD_REGISTRY_URL", "http://localhost:8080/upload")

	ctx := context.Background()

	kafka_reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   data_topic,
		// GroupID: lookupStringEnv("STREAM_GRP_ID","work-ld-1"),
		// assign the logger to the reader
		// Logger: l,
	})

	kafka_writer := kafka.Writer{
		Addr:     kafka.TCP(brokerAddress),
		Topic:    control_topic,
		Balancer: &kafka.LeastBytes{},
	}

	counter := 0
	maxRows := int(lookupInt64Env("MAX_ROWS", 10))

	var testDataFile *os.File
	var fileName string

	for {
		fmt.Printf("Waiting for kafka mesages from broker %s on data topic %s, counter = %d \n", brokerAddress, data_topic , counter)

		msg, err := kafka_reader.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		// after receiving the message, log its value
		fmt.Println("received: ", string(msg.Value))

		var rawSensorData RawSensorData
		json.Unmarshal([]byte(string(msg.Value)), &rawSensorData)
		fmt.Println("received: ", rawSensorData.DeviceID)

		t := time.Now()
		if rawSensorData.TimeStamp == "" {
			rawSensorData.TimeStamp = fmt.Sprintf("%02d",t.UnixNano()/int64(time.Millisecond))
		}

		if rawSensorData.DeviceID != "" {
			if counter == 0 { // create new file and write header
				
				ts := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d:%02d",
					t.Year(), t.Month(), t.Day(),
					t.Hour(), t.Minute(), t.Second(), t.UnixNano()/int64(time.Millisecond))

				fileName = fmt.Sprintf("raw_sensor_training_data_%s.csv",ts)

				f, err := os.Create(fileName)
				testDataFile = f
				if err != nil {
					panic(err)
				}
				testDataFile.WriteString("deviceID,timeStamp,current,temperature,vibration,sound\n")
			}

			testDataFile.WriteString(fmt.Sprintf("%s,%s,%.1f,%.1f,%.1f,%.1f\n", rawSensorData.DeviceID, rawSensorData.TimeStamp, rawSensorData.Current, rawSensorData.Temperature, rawSensorData.Vibration, rawSensorData.Sound))

			counter = counter + 1
			if maxRows == counter { // upload the file
				counter = 0 //reset
				testDataFile.Close()
				uploadFileToModelRegistry(fileName, upload_url)
				publishControlMessage(fileName, kafka_writer, ctx)

			}
		}

	}

}

func uploadFileToModelRegistry(filename string, url string) []byte {

	caCert, err := ioutil.ReadFile("/keys/ssh-publickey")
    if err != nil {
        log.Fatal(err)
    }
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

	
	filetype := "file"
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
	// client := &http.Client{}
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

	fmt.Printf("uploading data file %s to url %s complete !!! \n", filename, url)
	return content
}

func publishControlMessage(filename string, kafka_writer kafka.Writer, cxt context.Context) {

	message := fmt.Sprintf("{\"command\": \"extract-data\", \"payload\" : \"%s\"}", filename)

	msg := kafka.Message{
		Value: []byte(message),
	}

	err := kafka_writer.WriteMessages(cxt, msg)

	if err != nil {
		log.Fatal(err)
	}
}
