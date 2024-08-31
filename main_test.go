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
)

var mux *http.ServeMux
var cookie string

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	config.Initialize()
	slog.SetLogLoggerLevel(config.GetLogLevel())
	db.Initialize()
	cache.Initialize()
	mux = router.Initialize()
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
}

func TestUploadedFiles(t *testing.T) {
	req, err := NewRequest("GET", nil, "/uploaded/file.txt", nil, nil, nil)
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
	var response dto.Result
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	expected := render.ResultOK
	expectedBytes, _ := json.Marshal(expected)
	if response != *render.ResultOK {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), string(expectedBytes))
	}
	t.Log(response)
	t.Log(rr.Header().Get("X-Request-ID"))
	t.Log(rr.Header().Get("Set-Cookie"))

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

	var response dto.Result
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	expected := render.ResultOK
	expectedBytes, _ := json.Marshal(expected)
	if response != *render.ResultOK {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), string(expectedBytes))
	}
	t.Log(response)
	t.Log(rr.Header().Get("Set-Cookie"))
	cookie = rr.Header().Get("Set-Cookie")
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

	var response dto.Result
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	expected := render.ResultOK
	expectedBytes, _ := json.Marshal(expected)
	if response != *render.ResultOK {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), string(expectedBytes))
	}
	t.Log(response)
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

	var response dto.Result
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	expected := render.ResultOK
	expectedBytes, _ := json.Marshal(expected)
	if response != *render.ResultOK {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), string(expectedBytes))
	}
	t.Log(response)
}

func TestCreateAdv(t *testing.T) {
	req, err := NewRequest("POST", H{"Cookie": cookie}, "/adv", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "expected response" // Replace with actual expected response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestUpdateAdv(t *testing.T) {
	req, err := NewRequest("PUT", H{"Cookie": cookie}, "/adv/123", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "expected response" // Replace with actual expected response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestGetAdv(t *testing.T) {
	req, err := NewRequest("GET", nil, "/adv/123", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "expected response" // Replace with actual expected response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestGetAdvList(t *testing.T) {
	req, err := NewRequest("GET", nil, "/adv", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "expected response" // Replace with actual expected response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestAddAdvPhoto(t *testing.T) {
	req, err := NewRequest("POST", H{"Cookie": cookie}, "/adv/123/photos", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "expected response" // Replace with actual expected response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestDeleteAdvPhoto(t *testing.T) {
	req, err := NewRequest("DELETE", H{"Cookie": cookie}, "/adv/123/photos/456", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "expected response" // Replace with actual expected response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestDeleteAdv(t *testing.T) {
	req, err := NewRequest("DELETE", H{"Cookie": cookie}, "/adv/123", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "expected response" // Replace with actual expected response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestLogoutMe(t *testing.T) {
	req, err := NewRequest("GET", H{"Cookie": cookie}, "/logout/me", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "expected response" // Replace with actual expected response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
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

	var response dto.Result
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	expected := render.ResultOK
	expectedBytes, _ := json.Marshal(expected)
	if response != *render.ResultOK {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), string(expectedBytes))
	}
	t.Log(rr.Body.String())
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
}

func TestUploadedFilesNotFound(t *testing.T) {
	req, err := NewRequest("GET", nil, "/uploaded/nonexistentfile.txt", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestLoginError(t *testing.T) {
	req, err := NewRequest("POST", nil, "/login", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	expected := "Unauthorized" // Replace with actual expected error response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestLogoutMeError(t *testing.T) {
	req, err := NewRequest("GET", H{"Cookie": cookie}, "/logout/me", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	expected := "Unauthorized" // Replace with actual expected error response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestLogoutAllError(t *testing.T) {
	req, err := NewRequest("GET", H{"Cookie": cookie}, "/logout/all", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	expected := "Unauthorized" // Replace with actual expected error response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestRegistrationError(t *testing.T) {
	req, err := NewRequest("POST", nil, "/registration", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := "Bad Request" // Replace with actual expected error response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestUpdatePasswordError(t *testing.T) {
	req, err := NewRequest("PUT", H{"Cookie": cookie}, "/password", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	expected := "Unauthorized" // Replace with actual expected error response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestUpdateUserError(t *testing.T) {
	req, err := NewRequest("PUT", H{"Cookie": cookie}, "/user", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	expected := "Unauthorized" // Replace with actual expected error response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestGetAdvNotFound(t *testing.T) {
	req, err := NewRequest("GET", H{"Cookie": cookie}, "/adv/nonexistent", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	expected := "Not Found" // Replace with actual expected error response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestGetAdvListError(t *testing.T) {
	req, err := NewRequest("GET", nil, "/adv", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expected := "Internal Server Error" // Replace with actual expected error response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestCreateAdvError(t *testing.T) {
	req, err := NewRequest("POST", H{"Cookie": cookie}, "/adv", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	expected := "Unauthorized" // Replace with actual expected error response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestUpdateAdvError(t *testing.T) {
	req, err := NewRequest("PUT", H{"Cookie": cookie}, "/adv/nonexistent", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	expected := "Not Found" // Replace with actual expected error response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestDeleteAdvError(t *testing.T) {
	req, err := NewRequest("DELETE", H{"Cookie": cookie}, "/adv/nonexistent", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	expected := "Not Found" // Replace with actual expected error response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestAddAdvPhotoError(t *testing.T) {
	req, err := NewRequest("POST", H{"Cookie": cookie}, "/adv/nonexistent/photos", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	expected := "Not Found" // Replace with actual expected error response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestDeleteAdvPhotoError(t *testing.T) {
	req, err := NewRequest("DELETE", H{"Cookie": cookie}, "/adv/nonexistent/photos/nonexistent", nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	expected := "Not Found" // Replace with actual expected error response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
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
