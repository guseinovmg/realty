package router

import (
	"log"
	"net/http"
	"realty/config"
	"realty/handlers"
	mw "realty/middleware"
)

func Initialize() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir(config.GetStaticFilesPath()))))
	mux.Handle("GET /uploaded/", http.StripPrefix("/uploaded/", http.FileServer(http.Dir(config.GetUploadedFilesPath()))))

	mux.Handle("POST /login", mw.Handler(handlers.Login, mw.SetAuthCookie, handlers.JsonOK).OnPanic(handlers.TextError))
	mux.Handle("GET /logout/me", mw.Handler(handlers.LogoutMe))
	mux.Handle("GET /logout/all", mw.Handler(mw.IsNotGracefullyStopped, mw.Auth, handlers.LogoutAll))
	mux.Handle("POST /registration", mw.Handler(mw.IsNotGracefullyStopped, handlers.Registration))
	mux.Handle("PUT /password", mw.Handler(mw.IsNotGracefullyStopped, mw.Auth, handlers.UpdatePassword))

	mux.Handle("PUT /user", mw.Handler(mw.IsNotGracefullyStopped, mw.Auth, mw.SetAuthCookie, handlers.UpdateUser).OnPanic(handlers.JsonError))

	mux.Handle("GET /adv/{advId}", mw.Handler(mw.FindAdv, handlers.GetAdv))
	mux.Handle("GET /adv", mw.Handler(handlers.GetAdvList))
	mux.Handle("POST /adv", mw.Handler(handlers.GetAdvList))

	mux.Handle("GET /user/adv/{advId}", mw.Handler(mw.Auth, mw.FindAdv, mw.CheckAdvOwner, mw.SetAuthCookie, handlers.GetUsersAdv))
	mux.Handle("GET /user/adv", mw.Handler(mw.Auth, mw.SetAuthCookie, handlers.GetUsersAdvList))
	mux.Handle("POST /user/adv", mw.Handler(mw.Auth, mw.SetAuthCookie, handlers.GetUsersAdvList))

	mux.Handle("POST /adv", mw.Handler(mw.IsNotGracefullyStopped, mw.Auth, mw.SetAuthCookie, handlers.CreateAdv))
	mux.Handle("PUT /adv/{advId}", mw.Handler(mw.IsNotGracefullyStopped, mw.Auth, mw.FindAdv, mw.CheckAdvOwner, mw.SetAuthCookie, handlers.UpdateAdv))
	mux.Handle("DELETE /adv/{advId}", mw.Handler(mw.IsNotGracefullyStopped, mw.Auth, mw.FindAdv, mw.CheckAdvOwner, mw.SetAuthCookie, handlers.DeleteAdv))

	mux.Handle("POST /adv/{advId}/photos", mw.Handler(mw.IsNotGracefullyStopped, mw.Auth, mw.FindAdv, mw.CheckAdvOwner, mw.SetAuthCookie, handlers.AddAdvPhoto))
	mux.Handle("DELETE /adv/{advId}/photos/{photoId}", mw.Handler(mw.IsNotGracefullyStopped, mw.Auth, mw.FindAdv, mw.CheckAdvOwner, mw.SetAuthCookie, handlers.DeleteAdvPhoto))

	go log.Fatal(http.ListenAndServe(config.GetHttpServerPort(), mux))

	return mux

}
