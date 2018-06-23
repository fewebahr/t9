package client

import "time"

const (
	serverAddressDefault     = `127.0.0.1:4239`
	connectionTimeoutDefault = 6 * time.Second
	requestTimeoutDefault    = 3 * time.Second
)
