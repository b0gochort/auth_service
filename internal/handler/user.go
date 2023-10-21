package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/b0gochort/internal/model"
	"github.com/valyala/fasthttp"
)

func (h *Handler) SignUp(ctx *fasthttp.RequestCtx) {
	if string(ctx.Method()) != "POST" {
		slog.Info("handler.SingUp: unsaporrted unsupported method")
		ctx.Error("unsaporrted unsupported method", fasthttp.StatusNotFound)

		return
	}
	var user model.User

	if err := json.Unmarshal(ctx.Request.Body(), &user); err != nil {
		slog.Info("hadnler.SignUp.Unmarshal: %s", err.Error())
		ctx.Error(fmt.Sprintf("json.Unmarshal : %s", err.Error()), fasthttp.StatusUnprocessableEntity)

		return
	}

	auth, err := h.services.UserService.SignUp(user)
	if err != nil {
		slog.Info("handler.SignUp.Userservice: %s", err.Error())
		ctx.Error(fmt.Sprintf("userService.SignUp : %s", err.Error()), fasthttp.StatusInternalServerError)

		return
	}

	body, err := json.Marshal(auth)
	if err != nil {
		slog.Info("hadnler.SignUp.Marshal: %s", err.Error())
		ctx.Error(fmt.Sprintf("json.Marshal : %s", err.Error()), fasthttp.StatusInternalServerError)

	}

	ctx.Write(body)
}

func (h *Handler) Login(ctx *fasthttp.RequestCtx) {

}
