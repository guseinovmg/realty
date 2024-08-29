package render

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"realty/dto"
	"time"
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
	nanoSec := time.Now().UnixNano() - requestId
	if err != nil {
		slog.Error("response", "resultId", requestId, "tm", fmt.Sprintf("%dns", nanoSec), "httpCode", statusCode, "msg", err.Error())
	} else {
		slog.Debug("response", "resultId", requestId, "tm", fmt.Sprintf("%dns", nanoSec), "httpCode", statusCode)
	}
}
