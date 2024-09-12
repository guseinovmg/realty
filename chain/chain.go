package chain

import (
	"fmt"
	"log/slog"
	"net/http"
	"realty/cache"
	"realty/config"
	"realty/metrics"
	"realty/utils"
	"strconv"
	"strings"
	"time"
)

type Result struct {
	StatusCode int
	WriteErr   error
	Body       string
}

var next Result = Result{
	StatusCode: -1,
	WriteErr:   nil,
}

func Next() Result {
	return next
}

// RequestContext
// можно расширять для передачи данных по цепочке обработчиков
type RequestContext struct {
	User      *cache.UserCache
	Adv       *cache.AdvCache
	RequestId int64
}

type HandlerFunction func(rc *RequestContext, writer http.ResponseWriter, request *http.Request) Result

type PanicHandlerFunction func(recovered any, rc *RequestContext, writer http.ResponseWriter, request *http.Request)

type Chain struct {
	onPanic  PanicHandlerFunction
	handlers []HandlerFunction
}

func (m *Chain) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	rc := &RequestContext{}
	rc.RequestId = utils.GenerateId()
	slog.Debug("request", "requestId", rc.RequestId, "method", request.Method, "pattern", request.Pattern, "path", request.URL.Path, "query", request.URL.RawQuery)
	writer.Header().Set("X-Request-ID", strconv.FormatInt(rc.RequestId, 10))
	defer func() {
		if err := recover(); err != nil {
			//todo в любом случае нужно создать оповещение админу(sms или email)
			metrics.IncPanicCounter()
			nanoSec := time.Now().UnixNano() - rc.RequestId
			slog.Error("panic", "requestId", rc.RequestId, "tm", fmt.Sprintf("%dns", nanoSec), "recovered", err)
			if m.onPanic != nil {
				m.onPanic(err, rc, writer, request)
			} else {
				writer.WriteHeader(500)
				_, _ = writer.Write([]byte("Internal error. RequestId=" + strconv.FormatInt(rc.RequestId, 10)))
			}
			//todo при некоторых паниках перезапуск сервиса не решит проблем, поэтому не завершаем работу
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
	var renderResult Result
	for _, f := range m.handlers {
		renderResult = f(rc, writer, request)
		if renderResult != next {
			break
		}
	}
	nanoSec := time.Now().UnixNano() - rc.RequestId
	if renderResult.WriteErr != nil {
		slog.Error("response", "requestId", rc.RequestId, "tm", fmt.Sprintf("%dns", nanoSec), "httpCode", renderResult.StatusCode, "msg", renderResult.WriteErr.Error())
	} else {
		if config.GetLogResponse() {
			slog.Debug("response", "requestId", rc.RequestId, "tm", fmt.Sprintf("%dns", nanoSec), "httpCode", renderResult.StatusCode, "body", renderResult.Body)
		} else {
			slog.Debug("response", "requestId", rc.RequestId, "tm", fmt.Sprintf("%dns", nanoSec), "httpCode", renderResult.StatusCode)
		}
	}
	if renderResult == next {
		slog.Error("unreached writing", "requestId", rc.RequestId, "path", request.URL.Path)
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
