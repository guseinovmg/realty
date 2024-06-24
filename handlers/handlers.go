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
	userCache := cache.FindUserCacheByLogin(login.Email)
	if userCache == nil {
		rd.Stop()
		writer.WriteHeader(401)
		writer.Write(utils.UnsafeStringToBytes("пользователь не найден"))
		return
	}
	if !bytes.Equal(utils.GeneratePasswordHash(login.Password), userCache.CurrentUser.PasswordHash) {
		rd.Stop()
		writer.WriteHeader(401)
		writer.Write(utils.UnsafeStringToBytes("пароль не верен"))
		return
	}

	rd.User = userCache

	writer.Write(utils.UnsafeStringToBytes("OK"))
}

func LogoutMe(rd *utils.RequestData, writer http.ResponseWriter, request *http.Request) {

}

func LogoutAll(rd *utils.RequestData, writer http.ResponseWriter, request *http.Request) {

}

func Registration(rd *utils.RequestData, writer http.ResponseWriter, request *http.Request) {

}

func UpdatePassword(rd *utils.RequestData, writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			//todo
		}
	}(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNoContent)
		return
	}
	requestDto := &dto.UpdatePasswordRequest{}
	err = json.Unmarshal(body, requestDto)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNoContent)
		return
	}
	if !bytes.Equal(rd.User.CurrentUser.PasswordHash, utils.GeneratePasswordHash(requestDto.OldPassword)) {
		http.Error(writer, "неверный пароль", http.StatusBadRequest)
		return
	}
	rd.User.CurrentUser.PasswordHash = utils.GeneratePasswordHash(requestDto.NewPassword)
	rd.User.CurrentUser.SessionSecret = utils.GenerateSessionsSecret()
	writer.Write(utils.UnsafeStringToBytes("ok"))
}

func UpdateUser(rd *utils.RequestData, writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			//todo
		}
	}(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNoContent)
		return
	}
	requestDto := &dto.UpdateUserRequest{}
	err = json.Unmarshal(body, requestDto)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNoContent)
		return
	}
	rd.User.CurrentUser.Name = requestDto.Name
	rd.User.CurrentUser.Description = requestDto.Description
	writer.Write(utils.UnsafeStringToBytes("ok"))
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
