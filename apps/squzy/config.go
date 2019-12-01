package main

import (
	"os"
	"strconv"
	"time"
)

var (
	port           int32
	clientAddress  string
	timeoutStorage time.Duration
)

func ReadConfig() {
	// Read port
	portValue := os.Getenv("PORT")
	port = 8080
	if portValue != "" {
		i, err := strconv.ParseInt(portValue, 10, 32)
		if err == nil {
			port = int32(i)
		}
	}
	// Read storage
	timeoutValue := os.Getenv("STORAGE_TIMEOUT")
	timeoutStorage = time.Second * 5
	if timeoutValue != "" {
		i, err := strconv.ParseInt(timeoutValue, 10, 32)
		if err == nil {
			timeoutStorage = time.Duration(int32(i)) * time.Second
		}
	}
	// Client Adress
	clientAddress = os.Getenv("STORAGE_HOST")
}
