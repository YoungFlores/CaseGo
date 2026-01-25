package server

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBName     string
	DBUser     string
	DBPassword string
	DBHost     string

	PublicKey  string
}

func LoadConfig() *Config {
	godotenv.Load()

	return &Config{
		DBHost:     os.Getenv("POSTGRES_HOST"),
		DBName:     os.Getenv("POSTGRES_DB"),
		DBUser:     os.Getenv("POSTGRES_USER"),
		DBPassword: os.Getenv("POSTGRES_PASSWORD"),
		PublicKey:  os.Getenv("PUBLIC_KEY"),
	}
}

func ParseRSAPublicKey(pemStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.New("unknown type of public key")
	}
}
