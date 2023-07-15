package site

import (
	"log"
	"strings"
	"sync"
	"yandex/internal/db"
	"yandex/internal/send"
)

func CheckTemplate(host db.Host, body []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	if strings.Contains(string(body), "краин") {
		db.SetTemplateError(host.Id, "найдена украина!")
		wg.Add(1)
		sendler.Handler("🆘🆘🆘 "+host.Name+" найдена украина!🆘🆘🆘 ", wg)
		log.Println(host.Name + " найдена украина!")
		return // прерываем проверку если найдена украина
	}
	if strings.Contains(string(body), "localhost") {
		db.SetTemplateError(host.Id, "найдена ошибка!")
		wg.Add(1)
		sendler.Handler("🆘🆘🆘 "+host.Name+" найдена ошибка 🆘🆘🆘 ", wg)
		log.Println(host.Name + " найдена ошибка")
		return // прерываем проверку если найдена ошибка
	}

	if host.TemplateError.Valid {
		db.SetTemplateError(host.Id, "")
	}
}
