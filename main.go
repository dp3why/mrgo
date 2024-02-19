package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dp3why/mrgo/backend"

	"github.com/dp3why/mrgo/handler"
)


func main() {
	fmt.Println("Starting the application...")
 
	backend.InitElasticsearchBackend()
	log.Fatal(http.ListenAndServe(":8080", handler.InitRouter()))
}
