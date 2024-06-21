package router

import (
	"log"
	"net/http"
	"realty/config"
	"realty/handlers"
	mw "realty/middleware"
	"realty/utils"
)

func Initialize() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir(config.GetStaticFilesPath()))))
	mux.Handle("GET /uploaded/", http.StripPrefix("/uploaded/", http.FileServer(http.Dir(config.GetUploadedFilesPath()))))

	mux.Handle("POST /login", utils.Handler(handlers.Login, mw.SetAuthCookie).OnPanic(handlers.TextError))
	mux.Handle("GET /logout/me", utils.Handler(handlers.LogoutMe))
	mux.Handle("GET /logout/all", utils.Handler(mw.Auth, handlers.LogoutAll))
	mux.Handle("POST /registration", utils.Handler(handlers.Registration))
	mux.Handle("PUT /password", utils.Handler(mw.Auth, mw.SetAuthCookie, handlers.UpdatePassword))

	mux.Handle("PUT /user", utils.Handler(mw.Auth, mw.SetAuthCookie, handlers.UpdateUser))

	mux.Handle("GET /adv/{id}", utils.Handler(handlers.GetAdv))
	mux.Handle("GET /adv", utils.Handler(handlers.GetAdvList))

	mux.Handle("GET /my/adv/{id}", utils.Handler(mw.Auth, mw.SetAuthCookie, handlers.GetUsersAdv))
	mux.Handle("GET /my/adv", utils.Handler(mw.Auth, mw.SetAuthCookie, handlers.GetUsersAdvList))

	mux.Handle("POST /adv", utils.Handler(mw.Auth, mw.SetAuthCookie, handlers.CreateAdv))
	mux.Handle("PUT /adv/{id}", utils.Handler(mw.Auth, mw.SetAuthCookie, handlers.UpdateAdv))
	mux.Handle("DELETE /adv/{id}", utils.Handler(mw.Auth, mw.SetAuthCookie, handlers.DeleteAdv))

	go log.Fatal(http.ListenAndServe(config.GetHttpServerPort(), mux))

	return mux

}
