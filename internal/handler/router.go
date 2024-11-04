package handler

import (
	"net"
)

type QuoteRepository interface {
	GetQuote() (string, error)
}

type HandlerFunc func(conn net.Conn, repo QuoteRepository)

const (
	HandlerQuote byte = 0x01 // Constant for the "quote" handler
	// You can define more handler constants here
)

type Router struct {
	routes map[byte]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[byte]HandlerFunc),
	}
}

func (r *Router) AddRoute(handlerID byte, handler HandlerFunc) {
	r.routes[handlerID] = handler
}

func (r *Router) GetRoute(handlerID byte) (HandlerFunc, bool) {
	handler, exists := r.routes[handlerID]
	return handler, exists
}
