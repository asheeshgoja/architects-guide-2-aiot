package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"fmt"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	http.ServeFile(w, r, "index.html")

	log.Printf("indexHandler called!");
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Printf("uploadHandler called!");
	ReceiveFile(w , r)
}


func ReceiveFile(w http.ResponseWriter, r *http.Request) {
    r.ParseMultipartForm(32 << 20) // limit your max input length!
    // var buf bytes.Buffer
    // in your case file would be fileupload
    file, header, err := r.FormFile("file")
    if err != nil {
        log.Println(err)
    }
    defer file.Close()


	// This is path which we want to store the file
    f, err := os.OpenFile("./"+ header.Filename, os.O_WRONLY|os.O_CREATE, 0666)

    if err != nil {
        log.Println(err)
    }

    // Copy the file to the destination path
    io.Copy(f, file)

	wd , _  := os.Getwd()

	fmt.Printf("Saved uploded file %s  to storage at %s", header.Filename, wd )
    // return handler.Filename, nil


    // name := strings.Split(header.Filename, ".")
    // fmt.Printf("File name %s\n", name[0])
    // // Copy the file data to my buffer
    // io.Copy(&buf, file)
    // // do something with the contents...
    // // I normally have a struct defined and unmarshal into a struct, but this will
    // // work as an example
    // contents := buf.String()
    // fmt.Println(contents)
    // // I reset the buffer in case I want to use it again
    // // reduces memory allocations in more intense projects
    // buf.Reset()
    // // do something else
    // // etc write header

	// // copy example
	// f, err := os.OpenFile("./downloaded", os.O_WRONLY|os.O_CREATE, 0666)
	// defer f.Close()
	// io.Copy(f, file)


    // return
}


func main() {

	// port := flag.String("p", "8100", "port to serve on")
	// directory := flag.String("d", ".", "the directory of static file to host")
	// flag.Parse()


	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/uploadmodel", uploadHandler)

	log.Printf("Serving Model Registry on port: 8080\n")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

	
}
