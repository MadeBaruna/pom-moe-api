package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var ProxyFile string

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file!")
	}

	ProxyFile = os.Getenv("PROXY_FILE")
}
