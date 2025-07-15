package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/servers/add", app.AddServerHandler)
	router.HandlerFunc(http.MethodPost, "/servers/remove", app.RemoveServerHandler)
	router.HandlerFunc(http.MethodGet, "/servers/list", app.ListServersHandler)

	router.HandlerFunc(http.MethodPost, "/proxies/add", app.addProxyHandler)
	router.HandlerFunc(http.MethodPost, "/proxies/remove", app.removeProxyHandler)
	router.HandlerFunc(http.MethodGet, "/proxies/list", app.listProxyHandler)

	return app.checkToken(router)
}
