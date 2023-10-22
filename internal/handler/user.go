package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/b0gochort/internal/model"
	"github.com/golang-jwt/jwt"
	"github.com/valyala/fasthttp"
)

func (h *Handler) SignUp(ctx *fasthttp.RequestCtx, start time.Time) {
	if string(ctx.Method()) != "POST" {
		slog.Info("handler.SingUp: unsupported method")
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

	user.IP = string(ctx.Request.Header.Peek("x-forwarded-for"))
	fmt.Println(user.IP)
	geo, err := getGeo(user.IP)
	if err != nil {
		slog.Info("handler.services.GetGeo:", err.Error())
		response := model.ResponseError{
			Code:        fasthttp.StatusInternalServerError,
			Description: "handler.services.GetGeo",
			Error:       err,
		}
		body, err := json.Marshal(&response)
		if err != nil {
			ctx.Error("json.Marshal", fasthttp.StatusInternalServerError)
		}
		ctx.Write(body)
	}

	position := [2]float64{geo.Latitude, geo.Longitude}

	user.Position[0] = position[0]
	user.Position[1] = position[1]

	user.Date.Create = time.Now().Unix()
	user.Date.Update = time.Now().Unix()

	auth, err := h.services.UserService.SignUp(user)
	if err != nil {
		slog.Info("handler.services.UserService:", err.Error())
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

	if err := h.services.CheckUserAuth(user.Email); err != nil {
		auth := model.Auth{
			Auht: false,
		}
		user.IP = string(ctx.Request.Header.Peek("x-forwarded-for"))

		geo, err := getGeo(user.IP)
		if err != nil {
			slog.Info("handler.services.GetGeo:", err.Error())
			response := model.ResponseError{
				Code:        fasthttp.StatusInternalServerError,
				Description: "handler.services.GetGeo",
				Error:       err,
			}
			body, err := json.Marshal(&response)
			if err != nil {
				ctx.Error("json.Marshal", fasthttp.StatusInternalServerError)
			}
			ctx.Write(body)
		}

		position := [2]float64{geo.Latitude, geo.Longitude}

		user.Position[0] = position[0]
		user.Position[1] = position[1]

		user.Date.Create = time.Now().Unix()
		user.Date.Update = time.Now().Unix()

		auth, err = h.services.FindUser(user)
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

		return
	}
	auth := model.Auth{
		Auht: true,
	}
	user.IP = string(ctx.Request.Header.Peek("x-forwarded-for"))

	geo, err := getGeo(user.IP)
	if err != nil {
		slog.Info("handler.services.GetGeo:", err.Error())
		response := model.ResponseError{
			Code:        fasthttp.StatusInternalServerError,
			Description: "handler.services.GetGeo",
			Error:       err,
		}
		body, err := json.Marshal(&response)
		if err != nil {
			ctx.Error("json.Marshal", fasthttp.StatusInternalServerError)
		}
		ctx.Write(body)
	}

	position := [2]float64{geo.Latitude, geo.Longitude}

	user.Position[0] = position[0]
	user.Position[1] = position[1]

	user.Date.Create = time.Now().Unix()
	user.Date.Update = time.Now().Unix()

	auth, err = h.services.FindUser(user)
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

func (h *Handler) AuthMiddleware(ctx *fasthttp.RequestCtx, start time.Time) {
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
			return []byte("salt"), nil
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
		slog.Info("handler.AuthMiddleware: ", err.Error())
		ctx.Error(fmt.Sprintf("handler.AuthMiddleware.UserExists: %s", err.Error()), fasthttp.StatusInternalServerError)
	}

	response := model.ResponseSuccess{
		Code:   fasthttp.StatusOK,
		Result: "success",
		Time:   time.Since(start).Nanoseconds(),
	}

	body, err := json.Marshal(response)
	if err != nil {
		slog.Info("hadnler.AuthMiddleware.Marshal:", err.Error())
		ctx.Error(fmt.Sprintf("json.AuthMiddleware : %s", err.Error()), fasthttp.StatusInternalServerError)

		return

	}

	ctx.Write(body)
}

func (h *Handler) ActivateAuthByEmail(ctx *fasthttp.RequestCtx, start time.Time) {
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

	idCode, err := h.services.SendConfirmationEmail(user.Email)
	if err != nil {
		slog.Info("handler.SendConfirmationEmail: %s", err.Error())
		response := model.ResponseError{
			Code:        fasthttp.StatusInternalServerError,
			Description: "handler.SendConfirmationEmail",
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
		Result: idCode,
		Time:   time.Since(start).Nanoseconds(),
	}

	body, err := json.Marshal(response)
	if err != nil {
		slog.Info("hadnler.AuthMiddleware.Marshal:", err.Error())
		ctx.Error(fmt.Sprintf("json.AuthMiddleware : %s", err.Error()), fasthttp.StatusInternalServerError)

		return

	}

	ctx.Write(body)
}

func (h *Handler) VerificationCode(ctx *fasthttp.RequestCtx, start time.Time) {
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

	var user model.Verification

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

	err := h.services.UserService.VerificateEmailCode(user.Code, user.Email)
	if err != nil {
		slog.Info("handler.VerificateEmailCode: %s", err.Error())
		response := model.ResponseError{
			Code:        fasthttp.StatusInternalServerError,
			Description: "handler.VerificateEmailCode",
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
		Result: "success",
		Time:   time.Since(start).Nanoseconds(),
	}

	body, err := json.Marshal(response)
	if err != nil {
		slog.Info("hadnler.AuthMiddleware.Marshal:", err.Error())
		ctx.Error(fmt.Sprintf("json.AuthMiddleware : %s", err.Error()), fasthttp.StatusInternalServerError)

		return

	}

	ctx.Write(body)

}

// func (h *Handler) UpdateUser(ctx *fasthttp.RequestCtx, start time.Time) {
// 	if string(ctx.Method()) != "POST" {
// 		slog.Info("handler.UpdateUser: unsupported method")
// 		response := model.ResponseError{
// 			Code:        fasthttp.StatusMethodNotAllowed,
// 			Description: "unsupported method",
// 			Error:       fmt.Errorf("method not allowed"),
// 		}
// 		body, err := json.Marshal(&response)
// 		if err != nil {
// 			ctx.Error("json.Marshal", fasthttp.StatusInternalServerError)
// 		}
// 		ctx.Write(body)

// 		return
// 	}
// }

func getGeo(ip string) (model.GeoResponse, error) {
	var result model.GeoResponse
	c := &fasthttp.Client{
		Dial: func(addr string) (net.Conn, error) {
			return fasthttp.DialTimeout(addr, time.Second*10)
		},
		MaxConnsPerHost: 1,
	}
	code, body, err := c.Get(nil, fmt.Sprintf("http://free.ipwhois.io/json/%s", ip))
	if err != nil {
		return model.GeoResponse{}, err
	}
	if code != 200 {
		return model.GeoResponse{}, fmt.Errorf("something wrong, status code: %d", code)
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return model.GeoResponse{}, err
	}

	return result, nil
}
