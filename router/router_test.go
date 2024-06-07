package router

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"realty/utils"
	"strings"
	"testing"
)

func TestAuth(t *testing.T) {
	handler := Initialize()

	type args struct {
		req *http.Request
	}

	tests := []struct {
		name      string
		request   *http.Request
		wantCode  int
		checkBody func(body string, t *testing.T)
	}{
		{
			name:      "must return http.StatusUnauthorized",
			request:   NewRequest("POST", nil, "/login", nil, nil, ""),
			wantCode:  http.StatusUnauthorized,
			checkBody: nil,
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
			if tt.checkBody != nil {
				tt.checkBody(resp.Body.String(), t)
			}
		})
	}
}

func NewRequest(method string, headers utils.H, url string, pathParams utils.H, queryParams utils.H, body string) *http.Request {
	if pathParams != nil {
		for k, v := range pathParams {
			url = strings.Replace(url, k, "{"+v+"}", 1)
		}
	}

	if strings.Contains(url, "{") || strings.Contains(url, "}") {
		panic("неправильно сформирован путь " + url)
	}
	var buf io.Reader
	if body != "" {
		buf = strings.NewReader(body)
	}

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		panic(fmt.Sprintf("fail to create request: %s", err.Error()))
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	if queryParams != nil {
		q := req.URL.Query()
		for k, v := range queryParams {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
	return req
}
