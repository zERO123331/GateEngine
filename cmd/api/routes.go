package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/addserver", app.AddServerHandler)
	router.HandlerFunc(http.MethodPost, "/removeserver", app.RemoveServerHandler)
	router.HandlerFunc(http.MethodGet, "/listservers", app.ListServersHandler)

	return app.checkToken(router)
}
