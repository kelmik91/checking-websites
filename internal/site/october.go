package site

import (
	"strings"
	"sync"
	"yandex/internal/db"
	"yandex/internal/logger"
	"yandex/internal/send"
)

func CheckOctober(host db.Host, body []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	bodyString := string(body)

	if strings.Contains(bodyString, "<meta name=\"site\"") {
		metaStartIndex := strings.Index(bodyString, "<meta name=\"site\"")
		metaEndIndex := strings.Index(bodyString[metaStartIndex:], ">") + metaStartIndex
		metaTag := bodyString[metaStartIndex:metaEndIndex]
		nameStartIndex := strings.Index(metaTag, "content=\"") + 9
		nameEndIndex := strings.Index(metaTag[nameStartIndex:], "\"") + nameStartIndex

		if host.Name != metaTag[nameStartIndex:nameEndIndex] && host.October.String != "meta не совпадает!" {
			db.SetOctober(host.Id, "meta не совпадает!")
			wg.Add(1)
			sendler.Handler(" 🔠‼️ "+host.Name+" - "+"meta не совпадает! ‼️🔠 \nОбнаружен: "+metaTag[nameStartIndex:nameEndIndex], wg)
			logger.WriteWork(host.Name + " - " + "meta не совпадает!")
		} else if host.Name == metaTag[nameStartIndex:nameEndIndex] && host.October.String != "" {
			db.SetOctober(host.Id, "")
		}

	} else if strings.Contains(bodyString, "<meta name='site'") {
		metaStartIndex := strings.Index(bodyString, "<meta name='site'")
		metaEndIndex := strings.Index(bodyString[metaStartIndex:], ">") + metaStartIndex
		metaTag := bodyString[metaStartIndex:metaEndIndex]
		nameStartIndex := strings.Index(metaTag, "content='") + 9
		nameEndIndex := strings.Index(metaTag[nameStartIndex:], "'") + nameStartIndex

		if host.Name != metaTag[nameStartIndex:nameEndIndex] && host.October.String != "meta не совпадает!" {
			db.SetOctober(host.Id, "meta не совпадает!")
			wg.Add(1)
			sendler.Handler(" 🔠‼️ "+host.Name+" - "+"meta не совпадает! ‼️🔠 \nОбнаружен: "+metaTag[nameStartIndex:nameEndIndex], wg)
			logger.WriteWork(host.Name + " - " + "meta не совпадает!")
		} else if host.Name == metaTag[nameStartIndex:nameEndIndex] && host.October.String != "" {
			db.SetOctober(host.Id, "")
		}

	} else if host.October.String != "meta не установлен!" {
		db.SetOctober(host.Id, "meta не установлен!")
		wg.Add(1)
		sendler.Handler(" 🔠 "+host.Name+" - "+"meta не установлен! 🔠 ", wg)
		logger.WriteWork(host.Name + " - " + "meta не установлен!")
	}
}
