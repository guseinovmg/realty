package render

import (
	"encoding/json"
	"net/http"
	"realty/chain"
	"realty/config"
	"realty/dto"
	"realty/utils"
)

var ResultOK = &dto.Result{Result: "OK"}

func Json(writer http.ResponseWriter, statusCode int, v any) chain.Result {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("X-Content-Type-Options", "nosniff")
	writer.WriteHeader(statusCode)
	result := chain.Result{
		StatusCode: statusCode,
	}
	if config.GetLogResponse() {
		bytes, err := json.Marshal(v)
		if err != nil {
			result.WriteErr = err
		} else {
			result.Body = utils.UnsafeBytesToString(bytes)
			_, err = writer.Write(bytes)
			if err != nil {
				result.WriteErr = err
			}
		}
	} else {
		result.WriteErr = json.NewEncoder(writer).Encode(v)
	}

	return result
}
