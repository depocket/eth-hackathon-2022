package main

import (
	"depocket.io/app"
	"github.com/joho/godotenv"
	"log"

	"os"
)

var s app.Server

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Print(err)
	}
	if err := s.Run(os.Getenv("SERVER_ENV")); err != nil {
		log.Print(err)
	}
}
