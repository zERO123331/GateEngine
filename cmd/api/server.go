package main

import (
	"GateEngine/internal/data"
	"bytes"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func (app *application) AddServerHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("Server add request received")
	app.mutex.Lock()
	var serverStruct = struct {
		Name     string `json:"name"`
		Address  string `json:"address"`
		Kind     string `json:"kind"`
		Fallback bool   `json:"fallback"`
	}{}
	err := app.readJSON(w, r, &serverStruct)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid JSON")
		app.mutex.Unlock()
		return
	}

	address := strings.Split(serverStruct.Address, ":")
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
	err = app.RegisterServer(*newServer)
	if err != nil {
		app.logger.Error(err.Error())
		app.UnregisterServer(newServer.Name)
		app.mutex.Unlock()
		return
	}
	app.mutex.Unlock()
	app.logger.Info("Server added", "Name", serverStruct.Name, "Kind", serverStruct.Kind, "Fallback", serverStruct.Fallback, "Address", serverStruct.Address)
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
	app.UnregisterServer(serverStruct.Name)
	app.logger.Info("Server removed", "Name", serverStruct.Name)
	w.WriteHeader(200)
}

func (app *application) ListServersHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("Server list request received")
	app.mutex.Lock()
	w.Header().Set("Content-Type", "application/json")
	sort.Slice(app.servers, func(i, j int) bool {
		return app.servers[i].Name < app.servers[j].Name
	})
	body, err := serverList(app.servers)
	if err != nil {
		w.WriteHeader(500)
		app.mutex.Unlock()
		return
	}
	app.mutex.Unlock()
	w.WriteHeader(200)
	_, err = w.Write(body.Bytes())
	if err != nil {
		app.logger.Error(err.Error())
		return
	}
}

func (app *application) RegisterServer(s data.Server) error {
	body, err := serverList([]*data.Server{
		&s,
	})
	if err != nil {
		app.logger.Error(err.Error())
		return err
	}
	app.logger.Info("Registering server", "Name", s.Name)

	for _, p := range app.proxies {
		url := "http://" + p.APIAddress.String() + "/servers/add"

		req, err := http.NewRequest(http.MethodPost, url, body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", p.Secret)
		if err != nil {
			app.logger.Error(err.Error())
			return err
		}
		res, err := app.client.Do(req)
		if err != nil {
			if res != nil && res.StatusCode != 200 {
				app.logger.Error(res.Status)
				return err
			}
			app.logger.Error(err.Error())
			return err
		}

	}
	return nil
}

func (app *application) UnregisterServer(name string) {
	serverStruct := struct {
		Name string `json:"name"`
	}{
		Name: name,
	}
	for _, p := range app.proxies {
		url := "http://" + p.APIAddress.String() + "/servers/remove"

		body, err := json.Marshal(serverStruct)
		if err != nil {
			app.logger.Error(err.Error())
			return
		}
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", p.Secret)
		if err != nil {
			app.logger.Error(err.Error())
			return
		}
		res, err := app.client.Do(req)
		if err != nil {
			if res != nil && res.StatusCode != 200 {
				app.logger.Error(res.Status)
				return
			}
			app.logger.Error(err.Error())
			return
		}
	}
}
