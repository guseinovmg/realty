package chain

import (
	"fmt"
	"log/slog"
	"net/http"
	"realty/application"
	"realty/cache"
	"realty/config"
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

// RequestData
// можно расширять для передачи данных по цепочке обработчиков
type RequestData struct {
	User      *cache.UserCache
	Adv       *cache.AdvCache
	RequestId int64 //также используется в качестве времени старта запроса в ns
	chain     *Chain
}

func (rd *RequestData) Logger() *slog.Logger {
	return rd.chain.logger
}

func (rd *RequestData) Timeout() bool {
	return time.Now().UnixNano()-rd.RequestId > rd.chain.timeoutNs
}

func (rd *RequestData) GetOnTimeout() HandlerFunction {
	return rd.chain.onTimeout
}

type HandlerFunction func(rd *RequestData, writer http.ResponseWriter, request *http.Request) Result

type PanicHandlerFunction func(recovered any, rc *RequestData, writer http.ResponseWriter, request *http.Request)

type Chain struct {
	onPanic   PanicHandlerFunction
	handlers  []HandlerFunction
	logger    *slog.Logger
	timeoutNs int64
	onTimeout HandlerFunction
}

func (chain *Chain) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	rd := &RequestData{}
	rd.chain = chain
	var renderResult Result
	rd.RequestId = utils.GenerateId()
	chain.logger.Debug("request", "requestId", rd.RequestId, "method", request.Method, "pattern", request.Pattern, "path", request.URL.Path, "query", request.URL.RawQuery)
	writer.Header().Set("X-Request-ID", strconv.FormatInt(rd.RequestId, 10))
	defer func() {
		nanoSec := time.Now().UnixNano() - rd.RequestId
		application.Hit(request.Pattern, renderResult.StatusCode, nanoSec)
		if err := recover(); err != nil {
			//todo в любом случае нужно создать оповещение админу(sms или email)
			application.IncPanicCounter()
			chain.logger.Error("panic", "requestId", rd.RequestId, "tm", fmt.Sprintf("%dns", nanoSec), "recovered", err)
			if chain.onPanic != nil {
				chain.onPanic(err, rd, writer, request)
			} else {
				writer.WriteHeader(500)
				_, _ = writer.Write([]byte("Internal error. RequestId=" + strconv.FormatInt(rd.RequestId, 10)))
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
			application.GracefullyStopAndExitApp()
		}
	}()

	for _, f := range chain.handlers {
		renderResult = f(rd, writer, request)
		if renderResult != next {
			break
		}
		if rd.Timeout() {
			renderResult = rd.GetOnTimeout()(rd, writer, request)
			break
		}
	}
	nanoSec := time.Now().UnixNano() - rd.RequestId
	if renderResult.WriteErr != nil {
		chain.logger.Error("response", "requestId", rd.RequestId, "tm", fmt.Sprintf("%dns", nanoSec), "httpCode", renderResult.StatusCode, "msg", renderResult.WriteErr.Error())
	} else {
		if config.GetLogResponse() {
			chain.logger.Debug("response", "requestId", rd.RequestId, "tm", fmt.Sprintf("%dns", nanoSec), "httpCode", renderResult.StatusCode, "body", renderResult.Body)
		} else {
			chain.logger.Debug("response", "requestId", rd.RequestId, "tm", fmt.Sprintf("%dns", nanoSec), "httpCode", renderResult.StatusCode) //httpCode=-1 значит что ответ не был отправлен клиенту
		}
	}
}

func (chain *Chain) OnPanic(handler PanicHandlerFunction) *Chain {
	chain.onPanic = handler
	return chain
}

func (chain *Chain) SetLogger(logger *slog.Logger) *Chain {
	chain.logger = logger
	return chain
}

func (chain *Chain) SetTimeout(timeoutNs int64) *Chain {
	chain.timeoutNs = timeoutNs
	return chain
}

func (chain *Chain) OnTimeout(handler HandlerFunction) *Chain {
	chain.onTimeout = handler
	return chain
}

func DefaultOnTimeout(rd *RequestData, writer http.ResponseWriter, request *http.Request) Result {
	writer.WriteHeader(http.StatusRequestTimeout)
	body := "Request timeout. RequestId=" + strconv.FormatInt(rd.RequestId, 10)
	_, err := writer.Write([]byte(body))
	return Result{
		StatusCode: http.StatusRequestTimeout,
		WriteErr:   err,
		Body:       body,
	}
}

func Handler(handlerFunc ...HandlerFunction) *Chain {
	return &Chain{
		handlers:  handlerFunc,
		logger:    slog.Default(),
		timeoutNs: 30 * 1000 * 1000 * 1000, //30 сек
		onTimeout: DefaultOnTimeout,
	}
}
