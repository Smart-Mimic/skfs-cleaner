package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	actionAdd    = "Add"
	actionRemove = "Remove"
)

func formatHex(b []byte) string {
	return fmt.Sprintf("%x", b)
}

func makeUniqueKey(devAddr []byte, sessionKey []byte) string {
	return strings.ToUpper(formatHex(devAddr)) + "_" + formatHex(sessionKey)
}

func CreateJSONFiles(config Config) {
	mu.Lock()
	defer mu.Unlock()

	today := time.Now().Format("20060102")
	baseDir := filepath.Join("updates", today)
	addDir := filepath.Join(baseDir, "added")
	removeDir := filepath.Join(baseDir, "removed")

	if _, err := os.Stat(baseDir); !os.IsNotExist(err) {
		log.Printf("Removing existing directory: %s\n", baseDir)
		err := os.RemoveAll(baseDir)
		if err != nil {
			log.Fatalf("failed to remove existing directory: %v", err)
		}
	}

	err := os.MkdirAll(addDir, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create 'added' directories: %v", err)
	}

	err = os.MkdirAll(removeDir, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create 'removed' directories: %v", err)
	}

	addedDevices := []DeviceUpdate{}
	removedDevices := []DeviceUpdate{}

	// Create a map with DevAddr + SessionKey as the unique key for devices from the database
	deviceMap := make(map[string]DeviceInfo)
	for _, device := range deviceInfoStore {
		uniqueKey := makeUniqueKey(device.DevAddr, device.NwkSEncKey)
		deviceMap[uniqueKey] = device
	}

	// Create a map with DevAddr + SessionKey as the unique key for devices from the route
	routeDeviceMap := make(map[string]RouteDevice)
	for _, routeDevice := range routeDevices {
		uniqueKey := strings.ToUpper(routeDevice.DevAddr) + "_" + strings.ToLower(routeDevice.SessionKey)
		routeDeviceMap[uniqueKey] = routeDevice
	}

	// Compare the devices from the database and the route devices
	for _, device := range deviceInfoStore {
		// Only process devices that have a valid DevAddr (non-empty)
		// log.Printf("DevEUI: %v DevAddr: %v SessionKey: %v", formatHex(device.DevEUI), formatHex(device.DevAddr), formatHex(device.NwkSEncKey))
		if len(device.DevAddr) == 0 {
			continue
		}

		devAddrStr := formatHex(device.DevAddr)
		uniqueKey := makeUniqueKey(device.DevAddr, device.NwkSEncKey)

		if _, exists := routeDeviceMap[uniqueKey]; !exists {
			addedDevices = append(addedDevices, DeviceUpdate{
				RouteID:    config.RouteID,
				DevAddr:    devAddrStr,
				SessionKey: formatHex(device.NwkSEncKey),
				MaxCopies:  config.MaxCopies,
				Action:     actionAdd,
			})
		}
	}

	for _, routeDevice := range routeDevices {
		uniqueKey := routeDevice.DevAddr + "_" + routeDevice.SessionKey
		if _, exists := deviceMap[uniqueKey]; !exists {
			removedDevices = append(removedDevices, DeviceUpdate{
				RouteID:    config.RouteID,
				DevAddr:    routeDevice.DevAddr,
				SessionKey: routeDevice.SessionKey,
				MaxCopies:  routeDevice.MaxCopies,
				Action:     actionRemove,
			})
		}
	}

	if len(addedDevices) > 0 {
		fmt.Println("Processing added devices:")
		writeJSONFiles(addDir, addedDevices)
	}

	if len(removedDevices) > 0 {
		fmt.Println("Processing removed devices:")
		writeJSONFiles(removeDir, removedDevices)
	}
}

func writeJSONFiles(baseDir string, updates []DeviceUpdate) {
	batchSize := 100
	totalBatches := (len(updates) + batchSize - 1) / batchSize
	for i := 0; i < len(updates); i += batchSize {
		end := i + batchSize
		if end > len(updates) {
			end = len(updates)
		}

		fileName := fmt.Sprintf("part_%d_%d.json", i, end-1)
		filePath := filepath.Join(baseDir, fileName)

		file, err := os.Create(filePath)
		if err != nil {
			log.Fatalf("failed to create JSON file: %v", err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		err = encoder.Encode(updates[i:end])
		if err != nil {
			log.Fatalf("failed to encode JSON: %v", err)
		}

		batchNumber := (i / batchSize) + 1
		log.Printf("Batch %d/%d written to %s\n", batchNumber, totalBatches, fileName)
	}
}
