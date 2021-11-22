package providerjson

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (p ProviderJSON) DataSourcesHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	dsRaw := strings.Split(req.URL.RequestURI(), DataSourcesPath)
	ds := strings.Split(dsRaw[1], "/")[0]
	if len(ds) > 0 {
		if err := json.NewEncoder(w).Encode(ResourceJSON{
			Resource: p.DataSourcesMap[ds],
		}); err != nil {
			w.Write([]byte(fmt.Sprintf("Marshall error: %+v", err)))
		}
	}
}

func (p ProviderJSON) ResourcesHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	dsRaw := strings.Split(req.URL.RequestURI(), ResourcesPath)
	ds := strings.Split(dsRaw[1], "/")[0]
	if len(ds) > 0 {
		if err := json.NewEncoder(w).Encode(ResourceJSON{
			Resource: p.ResourcesMap[ds],
		}); err != nil {
			w.Write([]byte(fmt.Sprintf("Marshall error: %+v", err)))
		}
	}
}

func (p *ProviderJSON) ListResources(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(p.Provider.Resources()); err != nil {
		panic(err)
	}
}

func (p *ProviderJSON) ListDataSources(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(p.Provider.DataSources()); err != nil {
		panic(err)
	}
}
