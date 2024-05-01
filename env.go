package main

import (
  "github.com/joho/godotenv"
  "log"
)


func loadEnv() {
  envFile, err := godotenv.Read(".env")
  if err != nil {
    log.Fatal("Error loading .env file")
  }
  println(envFile["DB_USER"])
}
