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
	app.mutex.Lock()
	var proxyStruct = struct {
		Name    string `json:"name"`
		Address string `json:"address"`
		Kind    string `json:"kind"`
		Secret  string `json:"secret"`
	}{}
	err := app.readJSON(w, r, &proxyStruct)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid JSON")
		app.mutex.Unlock()
		return
	}
	address := strings.Split(proxyStruct.Address, ":")
	port, err := strconv.Atoi(address[1])
	kind := strings.ToLower(proxyStruct.Kind)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid address")
		app.mutex.Unlock()
		return
	}

	newProxy := &data.Proxy{
		Name: proxyStruct.Name,
		Address: data.Address{
			IP:   address[0],
			Port: port,
		},
		Kind:   kind,
		Secret: proxyStruct.Secret,
	}

	for _, proxy := range app.proxies {
		if proxy.Name == newProxy.Name {
			app.errorResponse(w, r, http.StatusBadRequest, "Proxy already exists")
			app.mutex.Unlock()
			return
		}
	}

	app.proxies = append(app.proxies, newProxy)
	app.mutex.Unlock()
	app.logger.Info("Proxy added", "Name", proxyStruct.Name)
	w.WriteHeader(200)
}

func (app *application) removeProxyHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("Proxy remove request received")
	app.mutex.Lock()
	var proxyStruct = struct {
		Name string `json:"name"`
	}{}
	err := app.readJSON(w, r, &proxyStruct)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid JSON")
		app.mutex.Unlock()
		return
	}
	var proxies []*data.Proxy
	for i, proxy := range app.proxies {
		if proxy.Name != proxyStruct.Name {
			proxies = append(proxies, app.proxies[i])
		}
	}
	app.proxies = proxies
	app.mutex.Unlock()
	app.logger.Info("Proxy removed", "Name", proxyStruct.Name)
	w.WriteHeader(200)
}

func (app *application) listProxyHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("Proxy list request received")
	app.mutex.Lock()
	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(app.proxies)
	if err != nil {
		w.WriteHeader(500)
		app.mutex.Unlock()
		return
	}
	app.mutex.Unlock()
	w.WriteHeader(200)
	_, err = w.Write(jsonData)
	if err != nil {
		app.logger.Error(err.Error())
		return
	}
}
