package sendler

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
	"yandex/internal/db"
	"yandex/internal/logger"
	"yandex/internal/telegram"
)

var nameFileNight = "/nightWatch.txt"
var loc, _ = time.LoadLocation("Europe/Moscow")

func Handler(message string, wg *sync.WaitGroup) {
	defer wg.Done()

	day := time.Now().In(loc).Weekday()
	hour := time.Now().In(loc).Hour()

	wg.Add(1)
	if time.Saturday != day && time.Sunday != day && hour >= 10 && hour < 22 {
		sendAllUser(message, wg)
		//log.Println(message)
		//wg.Done()
	} else {
		writeNight(message, wg)
	}
}

func sendAllUser(message string, wg *sync.WaitGroup) {
	defer wg.Done()

	users := db.GetAllUser()

	for _, user := range users {
		wg.Add(1)
		go telegram.SendTelegram(user, message, wg)
	}

	logger.WriteWork("Уведомления отправлены\n")
}

func writeNight(message string, wg *sync.WaitGroup) {
	defer wg.Done()

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	if _, err := os.Stat(exPath); os.IsNotExist(err) {
		err := os.MkdirAll(exPath, 0777)
		if err != nil {
			log.Println(err)
		}
	}
	file, err := os.OpenFile(exPath+nameFileNight, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		//logger.WriteWork(err.Error())
		log.Println(err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	_, err = writer.WriteString(time.Now().In(loc).Format("02.01.2006 15:04:05") + " " + message + "\n")
	if err != nil {
		//logger.WriteWork(err.Error())
		log.Println(err)
		return
	}
}
