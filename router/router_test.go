package router

import (
	"log"
	"net/http"
	"net/http/httptest"
	"realty/dto"
	"realty/utils"
	"testing"
)

func TestAuth(t *testing.T) {
	handler := Initialize()

	tests := []struct {
		name          string
		request       *http.Request
		wantCode      int
		checkResponse func(resp *httptest.ResponseRecorder, t *testing.T)
	}{
		{
			name:          "must return http.StatusUnauthorized",
			request:       utils.NewRequest("POST", nil, "/login", nil, nil, nil),
			wantCode:      http.StatusUnauthorized,
			checkResponse: nil,
		},
		{
			name: "must return http.StatusOK",
			request: utils.NewRequest("POST", nil, "/login", nil, nil, &dto.LoginRequest{
				Email:    "guseinovmg@gmail.com",
				Password: "Password",
			}),
			wantCode:      http.StatusOK,
			checkResponse: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, tt.request)
			log.Println(resp.Body.String())
			if resp.Result().StatusCode != tt.wantCode {
				t.Fatalf("the status code should be [%d] but received [%d]", tt.wantCode, resp.Result().StatusCode)
			}
			if tt.checkResponse != nil {
				tt.checkResponse(resp, t)
			}
		})
	}
}
