package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"fmt"
	"strconv"	
)

func lookupStringEnv(envVar string, defaultValue string) string {
	envVarValue, ok := os.LookupEnv(envVar)
	if !ok { 
		return defaultValue
	}
	fmt.Println("Env var ",envVar , " = ",envVarValue)
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


func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	http.ServeFile(w, r, "index.html")

	log.Printf("indexHandler called!");
}

func uploadTrainingDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Printf("uploadTrainingDataHandler called!");
	ReceiveFile(w , r, "raw_training_data")
}


func ReceiveFile(w http.ResponseWriter, r *http.Request, folder string) {
    r.ParseMultipartForm(32 << 20) // limit your max input length!
    // var buf bytes.Buffer
    // in your case file would be fileupload
    file, header, err := r.FormFile("file")
    if err != nil {
        log.Println(err)
    }
	log.Printf("streamed uploded file");

    defer file.Close()

	// This is path which we want to store the file
    f, err := os.OpenFile("/data/"+ folder + "/" +  header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	log.Printf("processed uploaded file");

    if err != nil {
        log.Println(err)
    }

    // Copy the file to the destination path
    io.Copy(f, file)

	wd , _  := os.Getwd()

	// log.Printf("Saved uploaded file");
	log.Printf("Saved uploaded file %s  to storage at %s", header.Filename, wd )
}

func deleteAllData(w http.ResponseWriter, r *http.Request) {

	err := os.RemoveAll("./raw_training_data")

	if (err != nil){
		log.Println("Failed to delete contents of folder." + err.Error())
	}

	log.Printf("deleteAllData called!");
}


func main() {
	os.Chdir("/data")// PVC from longhorn and k3s
	err := os.Mkdir("raw_training_data", 0755)
    if err != nil {
        log.Print(err)
    }

	mux := http.NewServeMux()


	
	// mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/upload", uploadTrainingDataHandler)
	mux.HandleFunc("/deleteAllData", deleteAllData)

	fileServerHtml := http.FileServer(http.Dir("/data/raw_training_data"))
    mux.Handle("/", fileServerHtml)
	
	port := lookupStringEnv("PORT" , "8081")

	log.Printf("Serving TrainingDataStore service on port: %s\n", port)

	// if err := http.ListenAndServe(":" + port, mux); err != nil {
	if err := http.ListenAndServeTLS(":"+port, "/keys/ssh-publickey", "/keys/ssh-privatekey", mux); err != nil {
		log.Fatal(err)
	}

}

// find which port is in use 
//sudo lsof -i -P -n | grep LISTEN