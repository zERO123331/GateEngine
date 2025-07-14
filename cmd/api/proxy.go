package main

import (
	"GateEngine/internal/data"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func (app *application) addProxyHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("Proxy add request received")
	var proxyStruct = struct {
		Name    string `json:"name"`
		Address string `json:"address"`
		Kind    string `json:"kind"`
		Secret  string `json:"secret"`
	}{}
	app.readJSON(w, r, &proxyStruct)
	address := strings.Split(proxyStruct.Address, ":")
	port, err := strconv.Atoi(address[1])
	kind := strings.ToLower(proxyStruct.Kind)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid address")
		return
	}

	app.proxies = append(app.proxies, &data.Proxy{
		Name: proxyStruct.Name,
		Address: data.Address{
			IP:   address[0],
			Port: port,
		},
		Kind:   kind,
		Secret: proxyStruct.Secret,
	})
	app.logger.Info("Proxy added", "Name", proxyStruct.Name)
	w.WriteHeader(200)
}

func (app *application) removeProxyHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("Proxy remove request received")
	var proxyStruct = struct {
		Name string `json:"name"`
	}{}
	app.readJSON(w, r, &proxyStruct)
	var proxies []*data.Proxy
	for i, proxy := range app.proxies {
		if proxy.Name != proxyStruct.Name {
			proxies = append(proxies, app.proxies[i])
		}
	}
	app.proxies = proxies
	app.logger.Info("Proxy removed", "Name", proxyStruct.Name)
	w.WriteHeader(200)
}

func (app *application) listProxyHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("Proxy list request received")
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(app.proxies)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	w.Write(data)
}
