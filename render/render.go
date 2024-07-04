package render

import (
	"encoding/json"
	"net/http"
	"realty/dto"
)

var resultOkBytes = []byte(`{"result":"OK"}`)

func RenderLoginPage(writer http.ResponseWriter, errDto *dto.Err) error {
	return nil
}

func Json(writer http.ResponseWriter, statusCode int, v any) error {
	writer.WriteHeader(statusCode)
	return json.NewEncoder(writer).Encode(v)
}

func JsonOK(writer http.ResponseWriter, statusCode int) error {
	_, err := writer.Write(resultOkBytes)
	return err
}
