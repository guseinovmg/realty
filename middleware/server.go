package middleware

import (
	"log/slog"
	"net/http"
	"realty/cache"
	"realty/utils"
	"strconv"
	"strings"
)

// RequestData
// можно расширять для передачи данных по цепочке обработчиков
type RequestData struct {
	User      *cache.UserCache
	Adv       *cache.AdvCache
	RequestId int64
}

type HandlerFunction func(rd *RequestData, writer http.ResponseWriter, request *http.Request) (next bool)

type PanicHandlerFunction func(recovered any, rd *RequestData, writer http.ResponseWriter, request *http.Request)

type Chain struct {
	onPanic  PanicHandlerFunction
	handlers []HandlerFunction
}

func (m *Chain) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	rd := &RequestData{}
	rd.RequestId = utils.GenerateId()
	writer.Header().Set("X-Request-ID", strconv.FormatInt(rd.RequestId, 10))
	defer func() {
		if err := recover(); err != nil {
			slog.Error("panic", "requestId", rd.RequestId, "recovered", err)
			if m.onPanic != nil {
				m.onPanic(err, rd, writer, request)
			} else {
				writer.WriteHeader(500)
				_, _ = writer.Write([]byte("Internal error"))
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
	next := true
	for _, f := range m.handlers {
		next = f(rd, writer, request)
		if !next {
			break
		}
	}
	if next {
		slog.Error("wrong handler returning value", "requestId", rd.RequestId, "url", request.URL.RawPath)
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
