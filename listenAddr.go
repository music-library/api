package main

import (
	"fmt"
	"os"
)

func ListenAddr() string {
	host := GetHost()
	port := GetPort()
	return fmt.Sprintf("%s:%s", host, port)
}

func GetHost() string {
	host := os.Getenv("HOST")

	if len(host) == 0 {
		host = "localhost"
	}

	return host
}

func GetPort() string {
	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = "3001"
	}

	return port
}
