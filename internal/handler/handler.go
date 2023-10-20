package handler

import (
	"fmt"

	"github.com/b0gochort/internal/serivce"
	"github.com/valyala/fasthttp"
)

type Handler struct {
	services *serivce.Service
}

func NewHandler(services *serivce.Service) *Handler {
	return &Handler{
		services: services,
	}

}

func (h *Handler) InitRoutes(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/ping":
		ping(ctx)
	}
}

func ping(ctx *fasthttp.RequestCtx) {
	fmt.Println("pong")
}
