package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/b0gochort/internal/model"
	"github.com/golang-jwt/jwt"
	"github.com/valyala/fasthttp"
)

func (h *Handler) SignUp(ctx *fasthttp.RequestCtx, start time.Time) {

	if string(ctx.Method()) != "POST" {
		slog.Info("handler.SingUp: unsaporrted unsupported method")
		response := model.ResponseError{
			Code:        fasthttp.StatusMethodNotAllowed,
			Description: "unsaporrted unsupported method",
			Error:       fmt.Errorf("method not allowed"),
		}
		body, err := json.Marshal(&response)
		if err != nil {
			ctx.Error("json.Marshal", fasthttp.StatusInternalServerError)
		}
		ctx.Write(body)
		return
	}
	var user model.User

	if err := json.Unmarshal(ctx.Request.Body(), &user); err != nil {
		slog.Info("handler.unmarshal: %s", err.Error())
		response := model.ResponseError{
			Code:        fasthttp.StatusUnprocessableEntity,
			Description: "handler.unmarshal",
			Error:       err,
		}
		body, err := json.Marshal(&response)
		if err != nil {
			ctx.Error("json.Marshal", fasthttp.StatusInternalServerError)
		}
		ctx.Write(body)
		return
	}

	auth, err := h.services.UserService.SignUp(user)
	if err != nil {
		slog.Info("handler.services.UserService: %s", err.Error())
		response := model.ResponseError{
			Code:        fasthttp.StatusInternalServerError,
			Description: "handler.services.UserService",
			Error:       err,
		}
		body, err := json.Marshal(&response)
		if err != nil {
			ctx.Error("json.Marshal", fasthttp.StatusInternalServerError)
		}
		ctx.Write(body)
		return
	}

	response := model.ResponseSuccess{
		Code:   fasthttp.StatusOK,
		Result: auth,
		Time:   time.Since(start).Nanoseconds(),
	}

	body, err := json.Marshal(response)
	if err != nil {
		slog.Info("hadnler.SignUp.Marshal: %s", err.Error())
		ctx.Error(fmt.Sprintf("json.Marshal : %s", err.Error()), fasthttp.StatusInternalServerError)

		return
	}

	ctx.Write(body)
}

func (h *Handler) Login(ctx *fasthttp.RequestCtx, start time.Time) {
	if string(ctx.Method()) != "POST" {
		slog.Info("handler.Login: unsaporrted unsupported method")
		response := model.ResponseError{
			Code:        fasthttp.StatusMethodNotAllowed,
			Description: "unsaporrted unsupported method",
			Error:       fmt.Errorf("method not allowed"),
		}
		body, err := json.Marshal(&response)
		if err != nil {
			ctx.Error("json.Marshal", fasthttp.StatusInternalServerError)
		}
		ctx.Write(body)
		return
	}

	var user model.User

	if err := json.Unmarshal(ctx.Request.Body(), &user); err != nil {
		slog.Info("handler.unmarshal: %s", err.Error())
		response := model.ResponseError{
			Code:        fasthttp.StatusUnprocessableEntity,
			Description: "handler.unmarshal",
			Error:       err,
		}
		body, err := json.Marshal(&response)
		if err != nil {
			ctx.Error("json.Marshal", fasthttp.StatusInternalServerError)
		}
		ctx.Write(body)
		return
	}

	auth, err := h.services.FindUser(user)
	if err != nil {
		slog.Info("handler.services.Login: %s", err.Error())
		response := model.ResponseError{
			Code:        fasthttp.StatusInternalServerError,
			Description: "handler.services.UserService",
			Error:       err,
		}
		body, err := json.Marshal(&response)
		if err != nil {
			ctx.Error("json.Marshal", fasthttp.StatusInternalServerError)
		}
		ctx.Write(body)
		return
	}

	response := model.ResponseSuccess{
		Code:   fasthttp.StatusOK,
		Result: auth,
		Time:   time.Since(start).Nanoseconds(),
	}

	body, err := json.Marshal(response)
	if err != nil {
		slog.Info("hadnler.Login.Marshal: %s", err.Error())
		ctx.Error(fmt.Sprintf("json.Login : %s", err.Error()), fasthttp.StatusInternalServerError)

		return

	}

	ctx.Write(body)

}

func (h *Handler) AuthMiddleware(ctx *fasthttp.RequestCtx) {
	var accsessToken model.Token

	err := json.Unmarshal(ctx.Request.Body(), &accsessToken)

	if err := json.Unmarshal(ctx.Request.Body(), &accsessToken); err != nil {
		slog.Info("handler.unmarshal: %s", err.Error())
		response := model.ResponseError{
			Code:        fasthttp.StatusMethodNotAllowed,
			Description: "handler.unmarshal",
			Error:       err,
		}
		body, err := json.Marshal(&response)
		if err != nil {
			ctx.Error("json.Marshal", fasthttp.StatusInternalServerError)
		}
		ctx.Write(body)
		return
	}

	claims := &model.JwtCustomClaims{}

	token, err := jwt.ParseWithClaims(accsessToken.AccessToken,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return "salt", nil
		})
	if err != nil || !token.Valid {
		slog.Info("handler.AuthMidleware validation token: %s", err.Error())
		response := model.ResponseError{
			Code:        fasthttp.StatusForbidden,
			Description: "handler.AuthMidleware.validate",
			Error:       err,
		}
		body, err := json.Marshal(&response)
		if err != nil {
			ctx.Error("json.Marshal", fasthttp.StatusInternalServerError)
		}
		ctx.Write(body)
		return
	}

	if err = h.services.UserService.UserExists(claims.Id, claims.Login); err != nil {
		slog.Info("handler.AuthMiddleware: %s", err.Error())
		ctx.Error(fmt.Sprintf("handler.AuthMiddleware.UserExists: %s", err.Error()), fasthttp.StatusInternalServerError)
	}

}
