package main

import (
	"GateEngine/internal/data"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
)

func (app *application) AddServerHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("Server add request received")
	app.mutex.Lock()
	var serverStruct = struct {
		Name     string `json:"name"`
		Address  string `json:"address"`
		Fallback bool   `json:"fallback"`
	}{}
	err := app.readJSON(w, r, &serverStruct)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid JSON")
		app.mutex.Unlock()
		return
	}

	reg := regexp.MustCompile("[0-9.]:[0-9]")
	address := reg.Split(serverStruct.Address, 2)
	port, err := strconv.Atoi(address[1])
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid address")
		app.mutex.Unlock()
		return
	}

	newServer := &data.Server{
		Name: serverStruct.Name,
		Address: data.Address{
			IP:   address[0],
			Port: port,
		},
		Fallback: serverStruct.Fallback,
	}

	for _, server := range app.servers {
		if server.Name == newServer.Name {
			app.errorResponse(w, r, http.StatusBadRequest, "Server already exists")
			app.mutex.Unlock()
			return
		}
	}

	app.servers = append(app.servers, newServer)
	app.mutex.Unlock()
	app.logger.Info("Server added", "Name", serverStruct.Name)
	w.WriteHeader(200)
}

func (app *application) RemoveServerHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("Server remove request received")
	app.mutex.Lock()
	var serverStruct = struct {
		Name string `json:"name"`
	}{}
	err := app.readJSON(w, r, &serverStruct)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid JSON")
		app.mutex.Unlock()
		return
	}
	var servers []*data.Server
	for i, server := range app.servers {
		if server.Name != serverStruct.Name {
			servers = append(servers, app.servers[i])
		}
	}
	app.servers = servers
	app.mutex.Unlock()
	app.logger.Info("Server removed", "Name", serverStruct.Name)
	w.WriteHeader(200)
}

func (app *application) ListServersHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("Server list request received")
	app.mutex.Lock()
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(app.servers)
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
