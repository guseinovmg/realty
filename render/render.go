package render

import (
	"encoding/json"
	"net/http"
	"realty/dto"
)

var ResultOk = `{"result":"OK"}`

var resultOkBytes = []byte(ResultOk)

func RenderLoginPage(writer http.ResponseWriter, errDto *dto.Err) error {
	return nil
}

func Json(writer http.ResponseWriter, statusCode int, v any) error {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("X-Content-Type-Options", "nosniff")
	writer.WriteHeader(statusCode)
	return json.NewEncoder(writer).Encode(v)
}

func JsonOK(writer http.ResponseWriter, statusCode int) error {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("X-Content-Type-Options", "nosniff")
	writer.WriteHeader(statusCode)
	_, err := writer.Write(resultOkBytes)
	return err
}
