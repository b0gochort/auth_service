package pkg

import (
	"encoding/json"
	"time"

	"github.com/b0gochort/internal/model"
	"github.com/valyala/fasthttp"
)

func HTTPResponseSuccess(code int, res any, start time.Time) interface{} {
	response := model.ResponseSuccess{
		Code:   code,
		Result: res,
		Time:   time.Since(start).Nanoseconds(),
	}

	body, err := json.Marshal(response)
	if err != nil {
		return model.ResponseError{
			Code:        fasthttp.StatusInternalServerError,
			Description: "HTTPResponse.json.Marshal",
			Error:       err,
		}
	}

	return body
}

func HTTPResponseError(code int, err error, description string) interface{} {
	response := model.ResponseError{
		Code:        code,
		Error:       err,
		Description: description,
	}

	body, err := json.Marshal(response)
	if err != nil {
		return model.ResponseError{
			Code:        fasthttp.StatusInternalServerError,
			Description: "HTTPResponse.json.Marshal",
			Error:       err,
		}
	}

	return body
}
