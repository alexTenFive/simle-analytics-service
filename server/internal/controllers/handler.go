package controllers

import "net/http"

type Handler interface {
	Handle(http.ResponseWriter, *http.Request)
}
