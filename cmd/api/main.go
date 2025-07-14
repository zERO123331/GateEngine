package main

import (
	"GateEngine/internal/data"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type config struct {
	port    int
	address string
	secret  string
}

type application struct {
	logger  *slog.Logger
	config  *config
	proxies []*data.Proxy
	servers []*data.Server
	client  *http.Client
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := &http.Client{}

	app := &application{
		logger: logger,
		config: &config{
			port:    8080,
			address: "127.0.0.1",
			secret:  "secret",
		},
		proxies: []*data.Proxy{},
		servers: []*data.Server{},
		client:  client,
	}

	srv := &http.Server{
		Addr:         combineAddresses(app.config.address, app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	srv.ListenAndServe()
}
