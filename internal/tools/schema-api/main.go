package main

import (
	"log"
	"net/http"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/providerjson"
)

/*
/docs/v1/data-sources - list of data sources
/docs/v1/data-sources/{name} - info for a specific data source
/docs/v1/resources - list of resources
/docs/v1/resources/{name} - info for a specific resource
*/

func main() {
	data := providerjson.LoadData()

	mux := http.NewServeMux()
	// paths
	mux.HandleFunc(providerjson.DataSourcesList, data.ListDataSources)
	mux.HandleFunc(providerjson.ResourcesList, data.ListResources)

	mux.HandleFunc(providerjson.DataSourcesPath, data.DataSourcesHandler)
	mux.HandleFunc(providerjson.ResourcesPath, data.ResourcesHandler)

	log.Println("starting api service on localhost:8080") // TODO - Make this configurable
	log.Fatal(http.ListenAndServe(":8080", mux))
}
