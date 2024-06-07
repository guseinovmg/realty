package handlers

import (
	"net/http"
	"realty/utils"
)

func TextError(recovered any, rd *utils.RequestData, writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(500)
	writer.Write(utils.UnsafeStringToBytes("Internal errore"))
}

func Login(rd *utils.RequestData, writer http.ResponseWriter, request *http.Request) {
	panic("hahaha")
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
