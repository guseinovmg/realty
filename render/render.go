package render

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"realty/dto"
)

var ResultOK = &dto.Result{Result: "OK"}

func RenderLoginPage(writer http.ResponseWriter, errDto *dto.Err) error {
	return nil
}

func Json(requestId int64, writer http.ResponseWriter, statusCode int, v any) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("X-Content-Type-Options", "nosniff")
	writer.WriteHeader(statusCode)
	err := json.NewEncoder(writer).Encode(v)
	if err != nil {
		slog.Error("response", "resultId", requestId, "httpCode", statusCode, "msg", err.Error())
	} else {
		slog.Debug("response", "resultId", requestId, "httpCode", statusCode)
	}
}
