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
	ctx.Response.Header.Set("Content-Type", "application/json")
	switch string(ctx.Path()) {
	case "/signup":
		ping(ctx)
	case "/login":
		ping(ctx)
	default:
		ping(ctx)
	}
}

func ping(ctx *fasthttp.RequestCtx) {
	fmt.Println("pong")
}
