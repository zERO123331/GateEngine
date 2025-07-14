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
	var serverStruct = struct {
		Name     string `json:"name"`
		Address  string `json:"address"`
		Fallback bool   `json:"fallback"`
	}{}
	app.readJSON(w, r, &serverStruct)

	reg := regexp.MustCompile("[0-9.]:[0-9]")
	address := reg.Split(serverStruct.Address, 2)
	port, err := strconv.Atoi(address[1])
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid address")
		return
	}

	app.servers = append(app.servers, &data.Server{
		Name: serverStruct.Name,
		Address: data.Address{
			IP:   address[0],
			Port: port,
		},
		Fallback: serverStruct.Fallback,
	})
	app.logger.Info("Server added", "Name", serverStruct.Name)
	w.WriteHeader(200)
}

func (app *application) RemoveServerHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("Server remove request received")
	var serverStruct = struct {
		Name string `json:"name"`
	}{}
	app.readJSON(w, r, &serverStruct)
	for i, server := range app.servers {
		if server.Name != serverStruct.Name {
			app.servers = append(app.servers, app.servers[i])
			break
		}
	}
	app.logger.Info("Server removed", "Name", serverStruct.Name)

}

func (app *application) ListServersHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("Server list request received")
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(app.servers)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	w.Write(data)
}
