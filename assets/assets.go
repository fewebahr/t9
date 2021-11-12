package assets

import (
	"embed"
)

// Dictionary is a byte array containing the default English dictionary
//go:embed english.txt
var Dictionary []byte

// Cert is a byte array containing the default TLS certificate
//go:embed cert.pem
var Cert []byte

// Key is a byte array containing the default TLS key paired with aforementioned certificate
//go:embed key.pem
var Key []byte

// Frontend contains all files required by the frontend to function
//go:embed frontend
var Frontend embed.FS
