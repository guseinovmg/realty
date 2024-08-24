package metrics

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Metrics struct {
	InstanceStartTime           int64   `json:"instanceStartTime,omitempty"`
	FreeRAM                     int64   `json:"freeRAM,omitempty"`
	CPUTemp                     float64 `json:"cpuTemp,omitempty"`
	CPUConsumption              float64 `json:"cpuConsumption,omitempty"`
	UnSavedChangesQueueCount    int64   `json:"unSavedChangesCount,omitempty"`
	DiskUsagePercent            float64 `json:"diskUsagePercent,omitempty"`
	RecoveredPanicsCount        int64   `json:"recoveredPanicsCount,omitempty"`
	MaxRAMConsumptions          int64   `json:"maxRAMConsumptions,omitempty"`
	MaxCPUConsumptions          int64   `json:"maxCPUConsumptions,omitempty"`
	MaxRPS                      int64   `json:"maxRPS,omitempty"`
	MaxUnSavedChangesQueueCount int64   `json:"maxUnSavedChangesCount,omitempty"`
}

// GetInstanceStartTime returns the instance start time
func GetInstanceStartTime() int64 {
	return time.Now().Unix()
}

// GetFreeRAM returns the free RAM in bytes
func getFreeRAM() (int64, error) {
	const memInfoPath = "/proc/meminfo"

	file, err := os.Open(memInfoPath)
	if err != nil {
		return 0, fmt.Errorf("failed to open meminfo file: %v", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, fmt.Errorf("failed to read meminfo file: %v", err)
		}

		if strings.HasPrefix(line, "MemFree:") {
			fields := strings.Fields(line)
			if len(fields) < 2 {
				return 0, fmt.Errorf("invalid meminfo file format")
			}
			freeRAMStr := fields[1]
			freeRAM, err := strconv.ParseInt(freeRAMStr, 10, 64)
			if err != nil {
				return 0, fmt.Errorf("failed to parse free RAM: %v", err)
			}
			// Convert from kilobytes to bytes
			return freeRAM * 1024, nil
		}
	}

	return 0, fmt.Errorf("MemFree not found in meminfo file")
}

// GetCPUTemp returns the CPU temperature in degrees Celsius
func getCPUTemp() (float64, error) {
	const thermalZonePath = "/sys/class/thermal/thermal_zone0/temp"

	file, err := os.Open(thermalZonePath)
	if err != nil {
		return 0, fmt.Errorf("failed to open thermal zone file: %v", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return 0, fmt.Errorf("failed to read thermal zone file: %v", err)
	}

	tempStr := strings.TrimSpace(line)
	temp, err := strconv.ParseInt(tempStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse temperature: %v", err)
	}

	// Convert temperature from millidegrees Celsius to degrees Celsius
	tempCelsius := float64(temp) / 1000.0
	return tempCelsius, nil
}

func GetFreeRAM() int64 {
	freeRAM, err := getFreeRAM()
	if err != nil {
		log.Printf("Error getting free RAM: %v", err)
		return 0
	}
	return freeRAM
}

func GetCPUTemp() float64 {
	cpuTemp, err := getCPUTemp()
	if err != nil {
		log.Printf("Error getting CPU temperature: %v", err)
		return 0
	}
	return cpuTemp
}

// GetCPUConsumption returns the CPU consumption percentage
func GetCPUConsumption() float64 {
	// Mock data: Replace with actual CPU consumption retrieval logic
	return 25.0
}

// GetUnSavedChangesQueueCount returns the count of unsaved changes in the queue
func GetUnSavedChangesQueueCount() int64 {
	// Mock data: Replace with actual unsaved changes count retrieval logic
	return 10
}

// GetDiskUsagePercent returns the disk usage percentage
func GetDiskUsagePercent() float64 {
	// Mock data: Replace with actual disk usage retrieval logic
	return 75.0
}

// GetRecoveredPanicsCount returns the count of recovered panics
func GetRecoveredPanicsCount() int64 {
	// Mock data: Replace with actual recovered panics count retrieval logic
	return 5
}

// GetMaxRAMConsumptions returns the maximum RAM consumptions
func GetMaxRAMConsumptions() int64 {
	// Mock data: Replace with actual maximum RAM consumptions retrieval logic
	return 2048
}

// GetMaxCPUConsumptions returns the maximum CPU consumptions
func GetMaxCPUConsumptions() int64 {
	// Mock data: Replace with actual maximum CPU consumptions retrieval logic
	return 100
}

// GetMaxRPS returns the maximum requests per second
func GetMaxRPS() int64 {
	// Mock data: Replace with actual maximum RPS retrieval logic
	return 1000
}

// GetMaxUnSavedChangesQueueCount returns the maximum count of unsaved changes in the queue
func GetMaxUnSavedChangesQueueCount() int64 {
	// Mock data: Replace with actual maximum unsaved changes count retrieval logic
	return 50
}
