package router

import (
	"net/http"
	"realty/config"
	"realty/handlers"
	"realty/handlers_chain"
	mw "realty/middleware"
)

var serveMux *http.ServeMux

func Initialize() *http.ServeMux {
	if serveMux != nil {
		return serveMux
	}
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir(config.GetStaticFilesPath()))))

	mux.Handle("/metrics", handlers_chain.Handler(handlers.GetMetrics))

	mux.Handle("GET /generate/id", handlers_chain.Handler(mw.Auth, handlers.GenerateId))

	mux.Handle("POST /login", handlers_chain.Handler(mw.Login, mw.SetAuthCookie, handlers.JsonOK).OnPanic(handlers.TextError))
	mux.Handle("GET /logout/me", handlers_chain.Handler(handlers.LogoutMe))
	mux.Handle("GET /logout/all", handlers_chain.Handler(mw.CheckGracefullyStop, mw.Auth, handlers.LogoutAll))
	mux.Handle("POST /registration", handlers_chain.Handler(mw.CheckGracefullyStop, handlers.Registration))
	mux.Handle("PUT /password", handlers_chain.Handler(mw.CheckGracefullyStop, mw.Auth, mw.SetAuthCookie, handlers.UpdatePassword))

	mux.Handle("PUT /user", handlers_chain.Handler(mw.CheckGracefullyStop, mw.Auth, mw.SetAuthCookie, handlers.UpdateUser).OnPanic(handlers.JsonError))

	mux.Handle("GET /adv/{advId}", handlers_chain.Handler(mw.FindAdv, handlers.GetAdv))
	mux.Handle("GET /adv", handlers_chain.Handler(handlers.GetAdvList))

	mux.Handle("GET /user/adv/{advId}", handlers_chain.Handler(mw.Auth, mw.FindAdv, mw.CheckAdvOwner, mw.SetAuthCookie, handlers.GetUsersAdv))
	mux.Handle("GET /user/adv", handlers_chain.Handler(mw.Auth, mw.SetAuthCookie, handlers.GetUsersAdvList))
	mux.Handle("POST /user/adv", handlers_chain.Handler(mw.Auth, mw.SetAuthCookie, handlers.GetUsersAdvList))

	mux.Handle("POST /adv", handlers_chain.Handler(mw.CheckGracefullyStop, mw.Auth, mw.SetAuthCookie, handlers.CreateAdv))
	mux.Handle("PUT /adv/{advId}", handlers_chain.Handler(mw.CheckGracefullyStop, mw.Auth, mw.FindAdv, mw.CheckAdvOwner, mw.SetAuthCookie, handlers.UpdateAdv))
	mux.Handle("DELETE /adv/{advId}", handlers_chain.Handler(mw.CheckGracefullyStop, mw.Auth, mw.FindAdv, mw.CheckAdvOwner, mw.SetAuthCookie, handlers.DeleteAdv))

	mux.Handle("POST /adv/{advId}/photos", handlers_chain.Handler(mw.CheckGracefullyStop, mw.Auth, mw.FindAdv, mw.CheckAdvOwner, mw.SetAuthCookie, handlers.AddAdvPhoto))
	mux.Handle("DELETE /adv/{advId}/photos/{photoId}", handlers_chain.Handler(mw.CheckGracefullyStop, mw.Auth, mw.FindAdv, mw.CheckAdvOwner, mw.SetAuthCookie, handlers.DeleteAdvPhoto))

	serveMux = mux
	return mux

}
