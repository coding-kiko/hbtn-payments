package main

import (
	"log"
	"net/http"
	"os"
)

var (
	directory = os.Getenv("RECEIPTS_FOLDER_PATH")
	port      = os.Getenv("PORT")
)

func main() {
	http.Handle("/", http.FileServer(http.Dir(directory)))

	log.Printf("Serving %s on HTTP port: %s\n", directory, port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
