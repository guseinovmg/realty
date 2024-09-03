package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"realty/cache"
	"realty/render"
	"realty/utils"
	"strconv"
	"strings"
	"time"
)

// RequestData
// можно расширять для передачи данных по цепочке обработчиков
type RequestData struct {
	User      *cache.UserCache
	Adv       *cache.AdvCache
	RequestId int64
}

type HandlerFunction func(rd *RequestData, writer http.ResponseWriter, request *http.Request) render.Result

type PanicHandlerFunction func(recovered any, rd *RequestData, writer http.ResponseWriter, request *http.Request)

type Chain struct {
	onPanic  PanicHandlerFunction
	handlers []HandlerFunction
}

func (m *Chain) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	rd := &RequestData{}
	rd.RequestId = utils.GenerateId()
	slog.Debug("request", "requestId", rd.RequestId, "method", request.Method, "path", request.URL.Path, "query", request.URL.RawQuery)
	writer.Header().Set("X-Request-ID", strconv.FormatInt(rd.RequestId, 10))
	defer func() {
		if err := recover(); err != nil {
			nanoSec := time.Now().UnixNano() - rd.RequestId
			slog.Error("panic", "requestId", rd.RequestId, "tm", fmt.Sprintf("%dns", nanoSec), "recovered", err)
			if m.onPanic != nil {
				m.onPanic(err, rd, writer, request)
			} else {
				writer.WriteHeader(500)
				_, _ = writer.Write([]byte("Internal error. RequestId=" + strconv.FormatInt(rd.RequestId, 10)))
			}
			//todo при некоторых паниках перезапуск сервиса не решит проблем, поэтому не завершаем работу
			//todo в любом случае нужно создать оповещение админу(sms или email)
			switch err.(type) {
			case string:
				if strings.Contains(err.(string), "not implemented") {
					return
				}
			case error:
				if strings.Contains(err.(error).Error(), "not implemented") {
					return
				}
			}
			// в остальных случаях красиво завершаем работу
			cache.GracefullyStopAndExitApp()
		}
	}()
	var renderResult render.Result
	for _, f := range m.handlers {
		renderResult = f(rd, writer, request)
		if renderResult != render.Next() {
			break
		}
	}
	nanoSec := time.Now().UnixNano() - rd.RequestId
	if renderResult.WriteErr != nil {
		slog.Error("response", "requestId", rd.RequestId, "tm", fmt.Sprintf("%dns", nanoSec), "httpCode", renderResult.StatusCode, "msg", renderResult.WriteErr.Error())
	} else {
		slog.Debug("response", "requestId", rd.RequestId, "tm", fmt.Sprintf("%dns", nanoSec), "httpCode", renderResult.StatusCode, "body", renderResult.Body)
	}
	if renderResult == render.Next() {
		slog.Error("unreached writing", "requestId", rd.RequestId, "path", request.URL.Path)
	}
}

func (m *Chain) OnPanic(onPanic PanicHandlerFunction) *Chain {
	m.onPanic = onPanic
	return m
}

func Handler(handlerFunc ...HandlerFunction) *Chain {
	return &Chain{
		handlers: handlerFunc,
	}
}
