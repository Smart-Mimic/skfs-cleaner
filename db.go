package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"skfs-cleaner/internal"

	"github.com/jackc/pgx/v4"
	"google.golang.org/protobuf/proto"
)

var (
	deviceInfoStore []DeviceInfo
	mu              sync.Mutex
)

func ConnectDB(config Config) (*pgx.Conn, error) {
	dbHost := config.DBHost
	dbPort := config.DBPort
	dbUser := config.DBUsername
	dbPass := config.DBPassword

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/chirpstack", dbUser, dbPass, dbHost, dbPort)
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}
	return conn, nil
}

func ProcessDevices(config Config) {
	db, err := ConnectDB(config)
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close(context.Background())

	query := `SELECT dev_eui, device_session FROM device`
	rows, err := db.Query(context.Background(), query)
	if err != nil {
		log.Fatalf("failed to select devices: %v", err)
	}
	defer rows.Close()

	var devices []DeviceInfo

	for rows.Next() {
		var devEUI []byte
		var deviceSessionBytes []byte

		err := rows.Scan(&devEUI, &deviceSessionBytes)
		if err != nil {
			log.Printf("failed to scan device: %v", err)
			continue
		}

		var deviceSession internal.DeviceSession
		err = proto.Unmarshal(deviceSessionBytes, &deviceSession)
		if err != nil {
			log.Printf("failed to unmarshal device session for DevEUI %x: %v", devEUI, err)
			continue
		}

		devices = append(devices, DeviceInfo{
			DevEUI:     devEUI,
			DevAddr:    deviceSession.DevAddr,
			NwkSEncKey: deviceSession.NwkSEncKey,
		})
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error processing rows: %v", err)
	}

	mu.Lock()
	deviceInfoStore = devices
	mu.Unlock()

	log.Printf("Total devices processed: %d", len(deviceInfoStore))
}
