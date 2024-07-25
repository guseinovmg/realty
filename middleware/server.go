package middleware

import (
	"net/http"
	"realty/cache"
	"strings"
)

// RequestData
// можно расширять для передачи данных по цепочке обработчиков
type RequestData struct {
	User *cache.UserCache
	Adv  *cache.AdvCache
}

type HandlerFunction func(rd *RequestData, writer http.ResponseWriter, request *http.Request) (next bool)

type PanicHandlerFunction func(recovered any, rd *RequestData, writer http.ResponseWriter, request *http.Request)

type Chain struct {
	onPanic  PanicHandlerFunction
	handlers []HandlerFunction
}

func (m *Chain) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	rd := &RequestData{}
	defer func() {
		if err := recover(); err != nil {
			//todo логировать паники нужно
			if m.onPanic != nil {
				m.onPanic(err, rd, writer, request)
			} else {
				writer.WriteHeader(500)
				_, _ = writer.Write([]byte("Internal error"))
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
			//todo в любом случае нужно создать оповещение админу(sms или email)
		}
	}()
	for _, f := range m.handlers {
		if !f(rd, writer, request) {
			return
		}
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
