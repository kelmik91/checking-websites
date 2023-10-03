package telegram

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"targetPlus/internal/logger"
)

type Message struct {
	Ok     bool `json:"ok"`
	Result []struct {
		UpdateId     int `json:"update_id"`
		MyChatMember struct {
			Chat struct {
				Id        int    `json:"id"`
				FirstName string `json:"first_name"`
				Username  string `json:"username"`
				Type      string `json:"type"`
			} `json:"chat"`
			From struct {
				Id           int    `json:"id"`
				IsBot        bool   `json:"is_bot"`
				FirstName    string `json:"first_name"`
				Username     string `json:"username"`
				LanguageCode string `json:"language_code"`
			} `json:"from"`
			Date          int `json:"date"`
			OldChatMember struct {
				User struct {
					Id        int64  `json:"id"`
					IsBot     bool   `json:"is_bot"`
					FirstName string `json:"first_name"`
					Username  string `json:"username"`
				} `json:"user"`
				Status    string `json:"status"`
				UntilDate int    `json:"until_date,omitempty"`
			} `json:"old_chat_member"`
			NewChatMember struct {
				User struct {
					Id        int64  `json:"id"`
					IsBot     bool   `json:"is_bot"`
					FirstName string `json:"first_name"`
					Username  string `json:"username"`
				} `json:"user"`
				Status    string `json:"status"`
				UntilDate int    `json:"until_date,omitempty"`
			} `json:"new_chat_member"`
		} `json:"my_chat_member,omitempty"`
		Message struct {
			MessageId int `json:"message_id"`
			From      struct {
				Id           int    `json:"id"`
				IsBot        bool   `json:"is_bot"`
				FirstName    string `json:"first_name"`
				Username     string `json:"username"`
				LanguageCode string `json:"language_code"`
			} `json:"from"`
			Chat struct {
				Id        int    `json:"id"`
				FirstName string `json:"first_name"`
				Username  string `json:"username"`
				Type      string `json:"type"`
			} `json:"chat"`
			Date     int    `json:"date"`
			Text     string `json:"text"`
			Entities []struct {
				Offset int    `json:"offset"`
				Length int    `json:"length"`
				Type   string `json:"type"`
			} `json:"entities,omitempty"`
		} `json:"message,omitempty"`
		CallbackQuery struct {
			Id   string `json:"id"`
			From struct {
				Id           int    `json:"id"`
				IsBot        bool   `json:"is_bot"`
				FirstName    string `json:"first_name"`
				Username     string `json:"username"`
				LanguageCode string `json:"language_code"`
			} `json:"from"`
			Message struct {
				MessageId int `json:"message_id"`
				From      struct {
					Id        int64  `json:"id"`
					IsBot     bool   `json:"is_bot"`
					FirstName string `json:"first_name"`
					Username  string `json:"username"`
				} `json:"from"`
				Chat struct {
					Id        int    `json:"id"`
					FirstName string `json:"first_name"`
					Username  string `json:"username"`
					Type      string `json:"type"`
				} `json:"chat"`
				Date     int    `json:"date"`
				Text     string `json:"text"`
				Entities []struct {
					Offset int    `json:"offset"`
					Length int    `json:"length"`
					Type   string `json:"type"`
				} `json:"entities"`
				ReplyMarkup struct {
					InlineKeyboard [][]struct {
						Text         string `json:"text"`
						CallbackData string `json:"callback_data"`
					} `json:"inline_keyboard"`
				} `json:"reply_markup"`
			} `json:"message"`
			ChatInstance string `json:"chat_instance"`
			Data         string `json:"data"`
		} `json:"callback_query,omitempty"`
	} `json:"result"`
}

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

func GetUpdates() error {
	urlTel := "https://api.telegram.org/bot" + os.Getenv("telegramToken") + "/getUpdates"

	get, err := http.DefaultClient.Get(urlTel)
	if err != nil {
		return err
	}
	defer get.Body.Close()

	if get.StatusCode != 200 {
		return errors.New("status not 200")
	}

	var jsonRes Message
	body, err := io.ReadAll(get.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &jsonRes)
	if err != nil {
		return err
	}

	fmt.Println(jsonRes)

	return nil
}
