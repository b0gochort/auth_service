package handler

import (
	"fmt"
	"time"

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
	start := time.Now()
	switch string(ctx.Path()) {
	case "/signup":
		h.SignUp(ctx, start)
	case "/login":
		ping(ctx)
	default:
		ping(ctx)
	}
}

func ping(ctx *fasthttp.RequestCtx) {
	fmt.Println("pong")
}
