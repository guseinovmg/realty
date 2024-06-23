package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"realty/cache"
	"realty/dto"
	"realty/utils"
)

func TextError(recovered any, rd *utils.RequestData, writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(500)
	writer.Write(utils.UnsafeStringToBytes("Internal error"))
}

func Login(rd *utils.RequestData, writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		rd.Stop()
		return
	}
	login := dto.LoginRequest{}
	err = json.Unmarshal(body, &login)
	if err != nil {
		rd.Stop()
		return
	}
	user := cache.FindUserByLogin(login.Email)
	if user == nil {
		rd.Stop()
		http.Error(writer, "", 401)
		writer.WriteHeader(401)
		writer.Write(utils.UnsafeStringToBytes("пользователь не найден"))
		return
	}
	if !bytes.Equal(utils.GeneratePasswordHash(login.Password), user.PasswordHash) {
		rd.Stop()
		writer.WriteHeader(401)
		writer.Write(utils.UnsafeStringToBytes("пароль не верен"))
		return
	}

	rd.User = user

	writer.Write(utils.UnsafeStringToBytes("OK"))
}

func LogoutMe(rd *utils.RequestData, writer http.ResponseWriter, request *http.Request) {

}

func LogoutAll(rd *utils.RequestData, writer http.ResponseWriter, request *http.Request) {

}

func Registration(rd *utils.RequestData, writer http.ResponseWriter, request *http.Request) {

}

func UpdatePassword(rd *utils.RequestData, writer http.ResponseWriter, request *http.Request) {

}

func UpdateUser(rd *utils.RequestData, writer http.ResponseWriter, request *http.Request) {

}

func CreateAdv(rd *utils.RequestData, writer http.ResponseWriter, request *http.Request) {

}

func GetAdv(rd *utils.RequestData, writer http.ResponseWriter, request *http.Request) {

}

func GetAdvList(rd *utils.RequestData, writer http.ResponseWriter, request *http.Request) {

}

func GetUsersAdv(rd *utils.RequestData, writer http.ResponseWriter, request *http.Request) {

}

func GetUsersAdvList(rd *utils.RequestData, writer http.ResponseWriter, request *http.Request) {

}

func UpdateAdv(rd *utils.RequestData, writer http.ResponseWriter, request *http.Request) {

}

func DeleteAdv(rd *utils.RequestData, writer http.ResponseWriter, request *http.Request) {

}
