package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"realty/auth_token"
	"realty/cache"
	"realty/config"
	"realty/db"
	"realty/dto"
	"realty/metrics"
	"realty/render"
	"realty/router"
	"realty/validator"
	"strings"
	"testing"
	"time"
)

type H map[string]string

var mux *http.ServeMux
var cookie string
var advId int64
var photoId int64 = 1720360451151465000
var resultOKStr string

const timeSleepMs = 50
const userEmail = "guseinovmg@gmail.com"
const password = "12345678"
const newPassword = "123456789"

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	slog.Info("start", "time", time.Now().Format("2006/01/02 15:04:05"))
	config.Initialize()
	slog.SetLogLoggerLevel(config.GetLogLevel())
	db.Initialize()
	cache.Initialize()
	mux = router.Initialize()
	resultOKBytes, _ := json.Marshal(render.ResultOK)
	resultOKStr = string(resultOKBytes)
}

func TestToken(t *testing.T) {
	var sessionSecret = [24]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	var userId int64 = 1657567
	var expireTime = time.Now().Add(24 * time.Hour).UnixNano()
	token := auth_token.CreateToken(userId, expireTime, sessionSecret)
	unpackedUserId, unpackedExpireTime := auth_token.UnpackToken(token)
	if unpackedUserId != userId {
		t.Fatal("userId!=1")
	}
	if expireTime != unpackedExpireTime {
		fmt.Println(expireTime)
		fmt.Println(unpackedExpireTime)
		t.Fatal("expireTime")
	}
	if !auth_token.IsValidToken(token, sessionSecret) {
		t.Fatal("IsValidToken")
	}
}

func TestShuffle(t *testing.T) {
	var arr = [36]byte([]byte("fsdfsbjfhsjhfvsgefyefiw73wg72rgwehjgrtyu"))
	shuffled := auth_token.Shuffle(arr)
	fmt.Println(string(shuffled[:]))
	unShuffled := auth_token.UnShuffle(shuffled)
	fmt.Println(string(unShuffled[:]))
	if !bytes.Equal(arr[:], unShuffled[:]) {
		t.Fatal("Shuffle error")
	}
}

func TestStaticFiles(t *testing.T) {
	req, err := NewRequest("GET", nil, "/static/file.txt", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "bla bla bla" // Replace with actual expected content
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
	time.Sleep(timeSleepMs * time.Millisecond)
}

func TestStaticFilesNotFound(t *testing.T) {
	req, err := NewRequest("GET", nil, "/static/nonexistentfile.txt", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
	time.Sleep(timeSleepMs * time.Millisecond)
}

func TestMetricsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := dto.Metrics{
		InstanceStartTime:           metrics.GetInstanceStartTime().Format("2006/01/02 15:04:05"),
		UnSavedChangesQueueCount:    cache.GetToSaveCount(),
		RecoveredPanicsCount:        metrics.GetRecoveredPanicsCount(),
		MaxUnSavedChangesQueueCount: metrics.GetMaxUnSavedChangesQueueCount(),
	}

	var response dto.Metrics
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.InstanceStartTime != expected.InstanceStartTime {
		t.Errorf("instanceStartTime mismatch: got %v want %v", response.InstanceStartTime, expected.InstanceStartTime)
	}

	if response.UnSavedChangesQueueCount != expected.UnSavedChangesQueueCount {
		t.Errorf("unSavedChangesQueueCount mismatch: got %v want %v", response.UnSavedChangesQueueCount, expected.UnSavedChangesQueueCount)
	}

	if response.RecoveredPanicsCount != expected.RecoveredPanicsCount {
		t.Errorf("recoveredPanicsCount mismatch: got %v want %v", response.RecoveredPanicsCount, expected.RecoveredPanicsCount)
	}

	if response.MaxUnSavedChangesQueueCount != expected.MaxUnSavedChangesQueueCount {
		t.Errorf("maxUnSavedChangesQueueCount mismatch: got %v want %v", response.MaxUnSavedChangesQueueCount, expected.MaxUnSavedChangesQueueCount)
	}
	time.Sleep(timeSleepMs * time.Millisecond)
}

func TestRegistration(t *testing.T) {
	req, err := NewRequest("POST", nil, "/registration", nil, nil, &dto.RegisterRequest{
		Email:    userEmail,
		Name:     "Murad",
		Password: password,
		InviteId: "",
	})
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := resultOKStr // Replace with actual expected response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
	time.Sleep(timeSleepMs * time.Millisecond)
}

func TestLogin(t *testing.T) {
	req, err := NewRequest("POST", nil, "/login", nil, nil, &dto.LoginRequest{
		Email:    userEmail,
		Password: password,
	})
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := resultOKStr // Replace with actual expected response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
	cookie = rr.Header().Get("Set-Cookie")
	time.Sleep(timeSleepMs * time.Millisecond)
}

func TestUpdatePassword(t *testing.T) {
	req, err := NewRequest("PUT", H{"Cookie": cookie}, "/password", nil, nil, &dto.UpdatePasswordRequest{
		OldPassword: password,
		NewPassword: newPassword,
	})
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := resultOKStr // Replace with actual expected response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
	time.Sleep(timeSleepMs * time.Millisecond)
}

func TestLoginWithNewPassword(t *testing.T) {
	req, err := NewRequest("POST", nil, "/login", nil, nil, &dto.LoginRequest{
		Email:    userEmail,
		Password: newPassword,
	})
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := resultOKStr // Replace with actual expected response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
	cookie = rr.Header().Get("Set-Cookie")
	time.Sleep(timeSleepMs * time.Millisecond)
}

func TestUpdateUser(t *testing.T) {
	req, err := NewRequest("PUT", H{"Cookie": cookie}, "/user", nil, nil, &dto.UpdateUserRequest{
		Name:        "Mamluk",
		Description: "hah",
	})
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Log(rr.Body.String())
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := resultOKStr // Replace with actual expected response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
	time.Sleep(timeSleepMs * time.Millisecond)
}

func TestCreateAdv(t *testing.T) {
	req, err := NewRequest("POST", H{"Cookie": cookie}, "/adv", nil, nil, &dto.CreateAdvRequest{
		OriginLang:   1,
		TranslatedBy: 1,
		TranslatedTo: "ru",
		Title:        "Пентхаус",
		Description:  "Описание пентхауса",
		Photos:       "",
		Price:        22220,
		Currency:     "rub",
		Country:      "Russia",
		City:         "Москва",
		Address:      "ул. Пушкина, дом Кукушника",
		Latitude:     2,
		Longitude:    34,
		UserComment:  "",
	})
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response dto.CreateAdvResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response.AdvId == 0 {
		t.Fatalf("advId: %v", response.AdvId)
	}
	advId = response.AdvId
	time.Sleep(timeSleepMs * time.Millisecond)
}

func TestUpdateAdv(t *testing.T) {
	req, err := NewRequest("PUT", H{"Cookie": cookie}, fmt.Sprintf("/adv/%d", advId), nil, nil, &dto.UpdateAdvRequest{
		OriginLang:   1,
		TranslatedBy: 1,
		TranslatedTo: "ru",
		Title:        "Пентхаус",
		Description:  "Описание пентхауса",
		Photos:       "",
		Price:        22220,
		Currency:     "rub",
		Country:      "Russia",
		City:         "Москва",
		Address:      "ул. Пушкина, дом Кукушкина",
		Latitude:     2,
		Longitude:    34,
		UserComment:  "",
	})
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := resultOKStr // Replace with actual expected response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
	time.Sleep(timeSleepMs * time.Millisecond)
}

func TestGetAdv(t *testing.T) {
	req, err := NewRequest("GET", nil, fmt.Sprintf("/adv/%d", advId), nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response dto.GetAdvResponseItem

	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	time.Sleep(timeSleepMs * time.Millisecond)
}

func TestGetAdvList(t *testing.T) {
	req, err := NewRequest("GET", nil, "/adv", nil, H{"currency": "rub", "page": "1"}, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response dto.GetAdvListResponse

	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response.Count != 1 {
		t.Errorf("Handler returned wrong count value: got %v want %v", response.Count, 1)
	}
	time.Sleep(timeSleepMs * time.Millisecond)
}

func TestGenerateId(t *testing.T) {
	req, err := NewRequest("GET", H{"Cookie": cookie}, "/generate/id", nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response dto.GenerateIdResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !validator.IsValidUnixNanoId(response.Id) {
		t.Fatalf("is not valid id: %v", response.Id)
	}

	time.Sleep(timeSleepMs * time.Millisecond)
}

func TestAddAdvPhoto(t *testing.T) {
	req, err := NewRequest("POST", H{"Cookie": cookie}, fmt.Sprintf("/adv/%d/photos", advId), nil, nil, &dto.AddPhotoRequest{Filename: fmt.Sprintf("%d.png", photoId)})
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := resultOKStr // Replace with actual expected response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
	time.Sleep(timeSleepMs * time.Millisecond)
}

func TestDeleteAdvPhoto(t *testing.T) {
	req, err := NewRequest("DELETE", H{"Cookie": cookie}, fmt.Sprintf("/adv/%d/photos/%d", advId, photoId), nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := resultOKStr // Replace with actual expected response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
	time.Sleep(timeSleepMs * time.Millisecond)
}

func TestDeleteAdv(t *testing.T) {
	req, err := NewRequest("DELETE", H{"Cookie": cookie}, fmt.Sprintf("/adv/%d", advId), nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := resultOKStr // Replace with actual expected response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
	time.Sleep(timeSleepMs * time.Millisecond)
}

func TestLogoutMe(t *testing.T) {
	req, err := NewRequest("GET", H{"Cookie": cookie}, "/logout/me", nil, nil, nil) //todo
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := resultOKStr // Replace with actual expected response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
	time.Sleep(timeSleepMs * time.Millisecond)
}

func TestLogoutAll(t *testing.T) {
	req, err := NewRequest("GET", H{"Cookie": cookie}, "/logout/all", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := resultOKStr // Replace with actual expected response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
	time.Sleep(timeSleepMs * time.Millisecond)
}

func NewRequest(method string, headers H, url string, pathParams H, queryParams H, body any) (*http.Request, error) {
	if pathParams != nil {
		for k, v := range pathParams {
			url = strings.Replace(url, k, "{"+v+"}", 1)
		}
	}

	if strings.Contains(url, "{") || strings.Contains(url, "}") {
		return nil, errors.New("неправильно сформирован путь " + url)
	}
	var buf io.Reader
	if body != nil {
		bytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		buf = strings.NewReader(string(bytes))
	}

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("fail to create request: %s", err.Error()))
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
	return req, nil
}

func parseAuthToken(cookieString string) string {
	// Split the cookie string by semicolon and space to get individual key-value pairs
	cookieParts := strings.Split(cookieString, "; ")

	// Iterate through the parts to find the auth_token
	for _, part := range cookieParts {
		if strings.HasPrefix(part, "auth_token=") {
			// Split the part by '=' to get the value of auth_token
			authToken := strings.SplitN(part, "=", 2)[1]
			return authToken
		}
	}

	// If auth_token is not found, return an empty string
	return ""
}
