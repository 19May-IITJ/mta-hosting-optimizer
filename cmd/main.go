package main

import (
	"log"
	"mta2/internal/config"
	"mta2/internal/handler"
	"mta2/pkg/ipconfig"
	"net/http"
)

// TODO - Context handling + Request ID generation

func main() {
	port := "8080"
	result := ipconfig.NewMap()
	list := ipconfig.NewIPConfigList()
	// Load configuration
	x := config.LoadConfigThreshold()
	if err := config.LoadConfigIPConfiguration(result, list); err == nil {
		// Create HTTP server
		log.Println("ip configurations loaded successfully")
		// register handlers to endpoints
		http.HandleFunc("/hostnames", handler.RetrieveHostnames(x, result))
		http.HandleFunc("/refresh", handler.RefreshDataSet(result, list))
		log.Printf("Server listening on port %s\n", port)
		// Starting server
		log.Fatal(http.ListenAndServe(":"+port, nil))
	} else {
		log.Printf("Error %v \n-*-unable to launch application-*-\n", err)
	}
}
