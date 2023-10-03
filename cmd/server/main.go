package main

import (
	"github.com/joho/godotenv"
	"log"
	"targetPlus/internal/serverGo"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = serverGo.RunServer()
	if err != nil {
		log.Panicln(err.Error())
	}

}
