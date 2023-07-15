package site

import (
	"strings"
	"sync"
	"yandex/internal/db"
	"yandex/internal/send"
)

func checkGtm(host db.Host, body []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	if string(body) == "" {
		return
	}

	GTMs := db.GetGTM()

	for _, gtmValue := range GTMs {
		if strings.Contains(string(body), gtmValue) {
			if host.Gtm.String != gtmValue {
				db.SetGTM(host.Id, gtmValue)
				wg.Add(1)
				sendler.Handler("📜 "+host.Name+" найден "+gtmValue+" 📜", wg)
			}
			return // прерываем проверку если найден GTM из базы
		}
	}

	if strings.Contains(string(body), "GTM-") {
		if host.Gtm.String != "Обнаружен неизвестный GTM" {
			db.SetGTM(host.Id, "Обнаружен неизвестный GTM")
			wg.Add(1)
			sendler.Handler("📜 ⚠️ "+host.Name+" обнаружен неизвестный GTM ⚠️ 📜", wg)
		}
	} else {
		if host.Gtm.String != "GTM не найден" {
			db.SetGTM(host.Id, "GTM не найден")
			wg.Add(1)
			sendler.Handler("🔔 "+host.Name+" GTM не найден 🔔", wg)
			//logger.WriteBody(string(body))
		}
	}
}
