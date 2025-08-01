package main

import (
	"GateEngine/internal/data"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
)

type config struct {
	address data.Address
	secret  string
}

type application struct {
	logger  *slog.Logger
	config  *config
	mutex   sync.Mutex
	proxies []*data.Proxy
	servers []*data.Server
	client  *http.Client
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := &http.Client{}
	var mu sync.Mutex

	app := &application{
		logger: logger,
		config: &config{
			address: data.Address{
				IP:   "127.0.0.1",
				Port: 8080,
			},
			secret: "secret",
		},
		mutex:   mu,
		proxies: []*data.Proxy{},
		servers: []*data.Server{},
		client:  client,
	}

	srv := &http.Server{
		Addr:         app.config.address.String(),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	srv.ListenAndServe()
}
