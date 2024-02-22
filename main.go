package main

import (
	"log"
	"net/http"

	"github.com/dp3why/mrgo/backend"

	"github.com/dp3why/mrgo/handler"
)


func main() {
	log.Default().Println("Starting the application...")
	backend.InitGCSBackend()
	backend.InitElasticsearchBackend()
	log.Fatal(http.ListenAndServe(":8080", handler.InitRouter()))
}
