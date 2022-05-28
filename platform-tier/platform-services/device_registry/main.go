package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func lookupStringEnv(envVar string, defaultValue string) string {
	envVarValue, ok := os.LookupEnv(envVar)
	if !ok {
		return defaultValue
	}
	fmt.Println("Env var ", envVar, " = ", envVarValue)
	return envVarValue
}

func deleteAllData(w http.ResponseWriter, r *http.Request) {

	err := os.RemoveAll("./model_store")

	if err != nil {
		log.Println("Failed to delete contents of folder." + err.Error())
	}

	log.Printf("deleteAllData called!")
}

func confirmActivation(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		log.Printf("ConfirmActivation request from website! r.PostFrom = %v\n", r.PostForm)
		name := r.FormValue("device_id")
		log.Printf("name :%s", name)

		dat, err := ioutil.ReadFile("/data/deviceRegistry/device_registry.dat")
		if err != nil {
			fmt.Fprintf(w, "read err: %v", err)
		}

		validDevices := strings.Split(string(dat), "\n")

		valid := "FALSE"

		for _, v := range validDevices {
			if v == name {
				valid = "TRUE"
				break
			}
		}

		w.Write([]byte(valid))

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func provisionDevice(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		log.Printf("provisionDevice request from website! r.PostFrom = %v\n", r.PostForm)
		name := r.FormValue("device_id")

		f, err := os.OpenFile("/data/deviceRegistry/device_registry.dat", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		if _, err := f.WriteString(fmt.Sprintf("%s\n", name)); err != nil {
			log.Println(err)
		}

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {
	err := os.Chdir("/data") // PVC from longhorn and k3s
	if err != nil {
		log.Print(err)
	}

	err = os.Mkdir("deviceRegistry", 0755)
	if err != nil {
		log.Print(err)
	}

	mux := http.NewServeMux()

	os.MkdirAll("./deviceRegistry", os.ModePerm)

	mux.HandleFunc("/deleteAllData", deleteAllData)
	mux.HandleFunc("/confirmActivation", confirmActivation)
	mux.HandleFunc("/provisionDevice", provisionDevice)

	fileServerHtml := http.FileServer(http.Dir("/data/deviceRegistry"))
	mux.Handle("/", fileServerHtml)
	// http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	// 	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// 	http.ServeFile(w, req, "/data/deviceRegistry" + req.RequestURI)

	// })

	port := lookupStringEnv("PORT", "8082")

	log.Printf("Serving secure device Registry on port: %s\n", port)

	// err = http.ListenAndServe(":"+port, mux)
	err = http.ListenAndServeTLS(":"+port, "/keys/ssh-publickey", "/keys/ssh-privatekey", mux)
	if err != nil {
		log.Fatal(err)
	}

}

// find which port is in use
//sudo lsof -i -P -n | grep LISTEN
// curl -k --cacert localhost.crt -d '{"device_id":"1"}' -H 'Content-Type: application/json' https://localhost:8082/confirmActivation
// curl -k -d "device_id=223" -H 'Content-Type: application/x-www-form-urlencoded' https://10.0.0.30:30006/confirmActivation
// curl -k -d "device_id=14333616" -H 'Content-Type: application/x-www-form-urlencoded' https://10.0.0.30:30006/provisionDevice