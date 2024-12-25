package router

import (
	"net/http"
	"realty/api/handlers"
	mw "realty/api/middleware"
	"realty/chain"
	"realty/config"
)

var serveMux *http.ServeMux

func Initialize() *http.ServeMux {
	if serveMux != nil {
		return serveMux
	}
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir(config.GetStaticFilesPath()))))

	mux.Handle("/metrics", chain.Handler(handlers.GetMetrics))

	mux.Handle("GET /generate/id", chain.Handler(mw.Auth, handlers.GenerateId))

	mux.Handle("POST /login", chain.Handler(mw.Login, mw.SetAuthCookie, handlers.JsonOK).OnPanic(handlers.TextError))
	mux.Handle("GET /logout/me", chain.Handler(handlers.LogoutMe))
	mux.Handle("GET /logout/all", chain.Handler(mw.CheckGracefullyStop, mw.Auth, mw.StopIfUnsavedMoreThan(900), handlers.LogoutAll))
	mux.Handle("POST /registration", chain.Handler(mw.CheckGracefullyStop, mw.StopIfUnsavedMoreThan(900), handlers.Registration))
	mux.Handle("PUT /password", chain.Handler(mw.CheckGracefullyStop, mw.StopIfUnsavedMoreThan(900), mw.Auth, mw.SetAuthCookie, handlers.UpdatePassword))

	mux.Handle("PUT /user", chain.Handler(mw.CheckGracefullyStop, mw.StopIfUnsavedMoreThan(700), mw.Auth, mw.SetAuthCookie, handlers.UpdateUser).OnPanic(handlers.JsonError))

	mux.Handle("GET /adv/{advId}", chain.Handler(mw.FindAdv, handlers.GetAdv))
	mux.Handle("GET /adv", chain.Handler(handlers.GetAdvList))

	mux.Handle("GET /user/adv/{advId}", chain.Handler(mw.Auth, mw.FindAdv, mw.CheckAdvOwner, mw.SetAuthCookie, handlers.GetUsersAdv))
	mux.Handle("GET /user/adv", chain.Handler(mw.Auth, mw.SetAuthCookie, handlers.GetUsersAdvList))
	mux.Handle("POST /user/adv", chain.Handler(mw.Auth, mw.SetAuthCookie, handlers.GetUsersAdvList))

	mux.Handle("POST /adv", chain.Handler(mw.CheckGracefullyStop, mw.StopIfUnsavedMoreThan(500), mw.Auth, mw.SetAuthCookie, handlers.CreateAdv))
	mux.Handle("PUT /adv/{advId}", chain.Handler(mw.CheckGracefullyStop, mw.StopIfUnsavedMoreThan(300), mw.Auth, mw.FindAdv, mw.CheckAdvOwner, mw.SetAuthCookie, handlers.UpdateAdv))
	mux.Handle("DELETE /adv/{advId}", chain.Handler(mw.CheckGracefullyStop, mw.StopIfUnsavedMoreThan(300), mw.Auth, mw.FindAdv, mw.CheckAdvOwner, mw.SetAuthCookie, handlers.DeleteAdv))

	mux.Handle("POST /adv/{advId}/photos", chain.Handler(mw.CheckGracefullyStop, mw.StopIfUnsavedMoreThan(200), mw.Auth, mw.FindAdv, mw.CheckAdvOwner, mw.SetAuthCookie, handlers.AddAdvPhoto))
	mux.Handle("DELETE /adv/{advId}/photos/{photoId}", chain.Handler(mw.CheckGracefullyStop, mw.StopIfUnsavedMoreThan(200), mw.Auth, mw.FindAdv, mw.CheckAdvOwner, mw.SetAuthCookie, handlers.DeleteAdvPhoto))

	serveMux = mux
	return mux

}
