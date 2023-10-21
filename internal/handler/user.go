package handler

import (
	"encoding/json"
	"fmt"

	"github.com/b0gochort/internal/model"
	"github.com/valyala/fasthttp"
)

func (h *Handler) SignUp(ctx *fasthttp.RequestCtx) {

	if string(ctx.Method()) != "POST" {
		ctx.Error("unsaporrted unsupported method", fasthttp.StatusNotFound)
		return
	}
	var user model.User

	if err := json.Unmarshal(ctx.Request.Body(), &user); err != nil {
		ctx.Error(fmt.Sprintf("json.Unmarshal : %s", err.Error()), fasthttp.StatusInternalServerError)
		return
	}

	userId, err := h.services.UserService.SignUp(user)
	if err != nil {
		ctx.Error(fmt.Sprintf("json.Unmarshal : %s", err.Error()), fasthttp.StatusInternalServerError)

	}
}
