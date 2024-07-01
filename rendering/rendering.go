package rendering

import (
	"encoding/json"
	"net/http"
	"realty/dto"
)

func RenderLoginPage(writer http.ResponseWriter, errDto *dto.Err) error {
	return nil
}

func RenderJson(writer http.ResponseWriter, v any) error {
	return json.NewEncoder(writer).Encode(v)
}
