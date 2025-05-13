package common

import "github.com/gorilla/mux"

type HttpHandler interface {
	RegisterRoutes(router *mux.Router)
}
