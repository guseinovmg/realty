package handlers

import (
	"encoding/json"
	"io"
	"net/http"
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
		panic(err)
	}
	login := dto.LoginRequestDTO{}
	err = json.Unmarshal(body, &login)
	if err != nil {
		panic(err)
	}
	if login.Username != "Murad" {
		panic("not murad")
	}

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
