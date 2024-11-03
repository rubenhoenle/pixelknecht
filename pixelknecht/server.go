package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

const (
	heartbeatInterval = 5 * time.Second
)

func readEnvWithFallback(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetCommandererUrl() string {
	return readEnvWithFallback("COMMANDERER_URL", "http://commanderer.hoenle.xyz:9000")
}

func GetCommandererTcp() string {
	return readEnvWithFallback("COMMANDERER_TCP", "commanderer.hoenle.xyz:8999")
}

func SendHeartbeat() {
	conn, err := net.Dial("tcp", GetCommandererTcp())
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connected Commanderer to server")

	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	for range ticker.C {
		_, err := conn.Write([]byte("heartbeat\n"))
		if err != nil {
			fmt.Println("Error sending heartbeat:", err)
			return
		}
		fmt.Println("Sent heartbeat") // debug
	}
}
