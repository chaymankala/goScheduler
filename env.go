package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	envFile = flag.String("envFile", ".env", "The environment File")
)

func loadEnv() {
	// fmt.Print(*envFile)

	err := godotenv.Load(*envFile)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func getEnvVariable(envVar string) string {
	return os.Getenv(envVar)
}
