package main

import (
	"github.com/joho/godotenv"
	"github.com/kelmik91/weather"
	"log"
	"sync"
	"targetPlus/internal/db"
	"targetPlus/internal/telegram"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	message, _ := weather.Weather(55.6302, 37.6045, "2")

	wg := sync.WaitGroup{}
	users := db.GetDigitalUser()
	for _, userId := range users {
		wg.Add(1)
		go telegram.SendTelegram(userId, message, &wg)
	}

	wg.Wait()
}
