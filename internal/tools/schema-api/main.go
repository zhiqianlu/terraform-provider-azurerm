package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/providerjson"
)

var port = "8080"

func main() {
	data := providerjson.LoadData()

	if userPort := os.Getenv("ARM_API_SERVER_PORT"); userPort != "" {
		port = userPort
	}

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sig
		log.Printf("%s signal received, closing provider API server on port %s", sig, port)
		os.Exit(0)
	}()

	mux := http.NewServeMux()
	// paths
	mux.HandleFunc(providerjson.DataSourcesList, data.ListDataSources)
	mux.HandleFunc(providerjson.ResourcesList, data.ListResources)

	mux.HandleFunc(providerjson.DataSourcesPath, data.DataSourcesHandler)
	mux.HandleFunc(providerjson.ResourcesPath, data.ResourcesHandler)

	log.Printf("starting api service on localhost:%s", port)
	log.Println(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}
