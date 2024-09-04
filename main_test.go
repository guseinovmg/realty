package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"realty/cache"
	"realty/config"
	"realty/db"
	"realty/dto"
	"realty/metrics"
	"realty/render"
	"realty/router"
	"strings"
	"testing"
	"time"
)

var mux *http.ServeMux
var cookie string
var advId int64
var photoId int64 = 1720360451151465000
var resultOKStr string

const timeSleepMs = 50

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	config.Initialize()
	slog.SetLogLoggerLevel(config.GetLogLevel())
	db.Initialize()
	cache.Initialize()
	mux = router.Initialize()
	resultOKBytes, _ := json.Marshal(render.ResultOK)
	resultOKStr = string(resultOKBytes)
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

	expected := metrics.Metrics{
		InstanceStartTime:           metrics.GetInstanceStartTime(),
		FreeRAM:                     metrics.GetFreeRAM(),
		CPUTemp:                     metrics.GetCPUTemp(),
		CPUConsumption:              metrics.GetCPUConsumption(),
		UnSavedChangesQueueCount:    metrics.GetUnSavedChangesQueueCount(),
		DiskUsagePercent:            metrics.GetDiskUsagePercent(),
		RecoveredPanicsCount:        metrics.GetRecoveredPanicsCount(),
		MaxRAMConsumptions:          metrics.GetMaxRAMConsumptions(),
		MaxCPUConsumptions:          metrics.GetMaxCPUConsumptions(),
		MaxRPS:                      metrics.GetMaxRPS(),
		MaxUnSavedChangesQueueCount: metrics.GetMaxUnSavedChangesQueueCount(),
	}

	var response metrics.Metrics
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.InstanceStartTime != expected.InstanceStartTime {
		t.Errorf("InstanceStartTime mismatch: got %v want %v", response.InstanceStartTime, expected.InstanceStartTime)
	}

	if response.FreeRAM != expected.FreeRAM {
		t.Errorf("FreeRAM mismatch: got %v want %v", response.FreeRAM, expected.FreeRAM)
	}

	if response.CPUTemp != expected.CPUTemp {
		t.Errorf("CPUTemp mismatch: got %v want %v", response.CPUTemp, expected.CPUTemp)
	}

	if response.CPUConsumption != expected.CPUConsumption {
		t.Errorf("CPUConsumption mismatch: got %v want %v", response.CPUConsumption, expected.CPUConsumption)
	}

	if response.UnSavedChangesQueueCount != expected.UnSavedChangesQueueCount {
		t.Errorf("UnSavedChangesQueueCount mismatch: got %v want %v", response.UnSavedChangesQueueCount, expected.UnSavedChangesQueueCount)
	}

	if response.DiskUsagePercent != expected.DiskUsagePercent {
		t.Errorf("DiskUsagePercent mismatch: got %v want %v", response.DiskUsagePercent, expected.DiskUsagePercent)
	}

	if response.RecoveredPanicsCount != expected.RecoveredPanicsCount {
		t.Errorf("RecoveredPanicsCount mismatch: got %v want %v", response.RecoveredPanicsCount, expected.RecoveredPanicsCount)
	}

	if response.MaxRAMConsumptions != expected.MaxRAMConsumptions {
		t.Errorf("MaxRAMConsumptions mismatch: got %v want %v", response.MaxRAMConsumptions, expected.MaxRAMConsumptions)
	}

	if response.MaxCPUConsumptions != expected.MaxCPUConsumptions {
		t.Errorf("MaxCPUConsumptions mismatch: got %v want %v", response.MaxCPUConsumptions, expected.MaxCPUConsumptions)
	}

	if response.MaxRPS != expected.MaxRPS {
		t.Errorf("MaxRPS mismatch: got %v want %v", response.MaxRPS, expected.MaxRPS)
	}

	if response.MaxUnSavedChangesQueueCount != expected.MaxUnSavedChangesQueueCount {
		t.Errorf("MaxUnSavedChangesQueueCount mismatch: got %v want %v", response.MaxUnSavedChangesQueueCount, expected.MaxUnSavedChangesQueueCount)
	}
	time.Sleep(timeSleepMs * time.Millisecond)
}

func TestRegistration(t *testing.T) {
	req, err := NewRequest("POST", nil, "/registration", nil, nil, &dto.RegisterRequest{
		Email:    "guseinovmg@gmail.com",
		Name:     "Murad",
		Password: "12345678",
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
		Email:    "guseinovmg@gmail.com",
		Password: "12345678",
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

/*func TestUpdatePassword(t *testing.T) {
	req, err := NewRequest("PUT", H{"Cookie": cookie}, "/password", nil, nil, &dto.UpdatePasswordRequest{
		OldPassword: "12345678",
		NewPassword: "123456789",
	})	if err != nil {
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
}*/

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

type H map[string]string

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
