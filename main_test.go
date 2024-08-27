package main

import (
	"encoding/json"
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
	"realty/utils"
	"testing"
)

var mux *http.ServeMux

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	config.Initialize()
	slog.SetLogLoggerLevel(config.GetLogLevel())
	db.Initialize()
	cache.Initialize()
	mux = router.Initialize()
}

func TestStaticFiles(t *testing.T) {
	req, err := utils.NewRequest("GET", nil, "/static/file.txt", nil, nil, nil)
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
	req, err := utils.NewRequest("GET", nil, "/uploaded/file.txt", nil, nil, nil)
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
	req, err := utils.NewRequest("POST", nil, "/registration", nil, nil, &dto.RegisterRequest{
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

	expected := render.ResultOk // Replace with actual expected response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestLogin(t *testing.T) {
	req, err := utils.NewRequest("POST", nil, "/login", nil, nil, &dto.LoginRequest{
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

	expected := render.ResultOk // Replace with actual expected response
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestLogoutMe(t *testing.T) {
	req, err := utils.NewRequest("GET", nil, "/logout/me", nil, nil, nil)
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
	req, err := utils.NewRequest("GET", nil, "/logout/all", nil, nil, nil)
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

func TestUpdatePassword(t *testing.T) {
	req, err := utils.NewRequest("PUT", nil, "/password", nil, nil, nil)
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

func TestUpdateUser(t *testing.T) {
	req, err := utils.NewRequest("PUT", nil, "/user", nil, nil, nil)
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
	req, err := utils.NewRequest("GET", nil, "/adv/123", nil, nil, nil)
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
	req, err := utils.NewRequest("GET", nil, "/adv", nil, nil, nil)
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

func TestCreateAdv(t *testing.T) {
	req, err := utils.NewRequest("POST", nil, "/adv", nil, nil, nil)
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
	req, err := utils.NewRequest("PUT", nil, "/adv/123", nil, nil, nil)
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
	req, err := utils.NewRequest("DELETE", nil, "/adv/123", nil, nil, nil)
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
	req, err := utils.NewRequest("POST", nil, "/adv/123/photos", nil, nil, nil)
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
	req, err := utils.NewRequest("DELETE", nil, "/adv/123/photos/456", nil, nil, nil)
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

func TestStaticFilesNotFound(t *testing.T) {
	req, err := utils.NewRequest("GET", nil, "/static/nonexistentfile.txt", nil, nil, nil)
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
	req, err := utils.NewRequest("GET", nil, "/uploaded/nonexistentfile.txt", nil, nil, nil)
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
	req, err := utils.NewRequest("POST", nil, "/login", nil, nil, nil)
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
	req, err := utils.NewRequest("GET", nil, "/logout/me", nil, nil, nil)
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
	req, err := utils.NewRequest("GET", nil, "/logout/all", nil, nil, nil)
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
	req, err := utils.NewRequest("POST", nil, "/registration", nil, nil, nil)
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
	req, err := utils.NewRequest("PUT", nil, "/password", nil, nil, nil)
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
	req, err := utils.NewRequest("PUT", nil, "/user", nil, nil, nil)
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
	req, err := utils.NewRequest("GET", nil, "/adv/nonexistent", nil, nil, nil)
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
	req, err := utils.NewRequest("GET", nil, "/adv", nil, nil, nil)
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
	req, err := utils.NewRequest("POST", nil, "/adv", nil, nil, nil)
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
	req, err := utils.NewRequest("PUT", nil, "/adv/nonexistent", nil, nil, nil)
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
	req, err := utils.NewRequest("DELETE", nil, "/adv/nonexistent", nil, nil, nil)
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
	req, err := utils.NewRequest("POST", nil, "/adv/nonexistent/photos", nil, nil, nil)
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
	req, err := utils.NewRequest("DELETE", nil, "/adv/nonexistent/photos/nonexistent", nil, nil, nil)
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
