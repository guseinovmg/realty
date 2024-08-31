package render

import (
	"encoding/json"
	"net/http"
	"realty/dto"
)

var ResultOK = &dto.Result{Result: "OK"}

func RenderLoginPage(writer http.ResponseWriter, errDto *dto.Err) error {
	return nil
}

type RenderResult struct {
	StatusCode int
	WriteErr   error
}

var Next RenderResult = RenderResult{
	StatusCode: -1,
	WriteErr:   nil,
}

func Json(writer http.ResponseWriter, statusCode int, v any) RenderResult {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("X-Content-Type-Options", "nosniff")
	writer.WriteHeader(statusCode)
	err := json.NewEncoder(writer).Encode(v)
	return RenderResult{
		StatusCode: statusCode,
		WriteErr:   err,
	}
}
