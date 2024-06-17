package utils

import (
	"net/http"
	"realty/models"
)

// RequestData
// можно расширять для передачи данных по цепочке обработчиков
type RequestData struct {
	stop bool
	User *models.User
}

type HandlerFunction func(rd *RequestData, writer http.ResponseWriter, request *http.Request)

type PanicHandlerFunction func(recovered any, rd *RequestData, writer http.ResponseWriter, request *http.Request)

func (rd *RequestData) Stop() {
	rd.stop = true
}

type Chain struct {
	onPanic  PanicHandlerFunction
	handlers []HandlerFunction
}

func (m *Chain) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rd := &RequestData{}
	defer func() {
		if err := recover(); err != nil { //todo при некоторых паниках нужно действительно дать серверу перазагрузиться
			if m.onPanic != nil {
				m.onPanic(err, rd, w, r)
			} else {
				w.WriteHeader(500)
				_, _ = w.Write([]byte("Internal error"))
			}
			//todo в любом случае нужно создать оповещение админу(sms или email)
		}
	}()
	for _, f := range m.handlers {
		f(rd, w, r)
		if rd.stop {
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
