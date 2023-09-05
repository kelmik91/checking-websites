package telegram

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"targetPlus/internal/logger"
)

func SendTelegram(userId int, message string, wgSend *sync.WaitGroup) {
	defer wgSend.Done()

	var token = os.Getenv("telegramToken")
	message = url.QueryEscape(message)
	u := "https://api.telegram.org/bot" + token + "/sendMessage?chat_id=" + strconv.Itoa(userId) + "&text=" + message
	resp, err := http.Get(u)
	if err != nil {
		logger.WriteWorkTelegram(strconv.Itoa(userId) + err.Error())
		log.Println(err)
	}
	//body, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	logger.WriteWorkTelegram(err.Error())
	//} else {
	//	logger.WriteWorkTelegram(string(body) + "\n\n")
	//}

	defer resp.Body.Close()
}
