package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func setupLogging() *os.File {
	today := time.Now().Format("20060102")
	logDir := filepath.Join("runlog", today)

	if _, err := os.Stat(logDir); !os.IsNotExist(err) {
		err := os.RemoveAll(logDir)
		if err != nil {
			log.Fatalf("failed to remove existing log directory: %v", err)
		}
	}

	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create log directory: %v", err)
	}

	logFile := filepath.Join(logDir, "output.log")
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, f)

	log.SetOutput(multiWriter)

	return f
}

func main() {
	logFile := setupLogging()
	defer logFile.Close()

	config := LoadEnv()

	fmt.Printf("Loaded configuration: %+v\n", config)

	routeID := config.RouteID

	err := fetchRouteDevices(routeID)
	if err != nil {
		log.Fatalf("failed to fetch route devices: %v", err)
	}

	ProcessDevices(config)

	CreateJSONFiles(config)

	err = runAddActions(routeID, config)
	if err != nil {
		log.Fatalf("Error running add actions: %v", err)
	}

	err = runRemoveActions(routeID, config)
	if err != nil {
		log.Fatalf("Error running remove actions: %v", err)
	}

	time.Sleep(10 * time.Minute)
}
