package site

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"yandex/internal/db"
	"yandex/internal/logger"
	"yandex/internal/send"
)

func CheckSite(host db.Host, wg *sync.WaitGroup) {
	defer wg.Done()

	client := http.Client{
		//CheckRedirect: func(req *http.Request, via []*http.Request) error {
		//	//fmt.Println(req.Response.StatusCode) // если не 301 сообщить об отсутствии редиректа на https
		//	//return fmt.Errorf("Not 301 redirect")
		//	return nil
		//},
		//Timeout: time.Second * 30,
	}

	resp, err := client.Get("http://" + host.Name)
	if err != nil {
		if host.Header.Int64 != 0 {
			time.Sleep(time.Second * 30)
			resp, err = client.Get("http://" + host.Name)
			if err != nil {
				if host.Header.Int64 != 0 {
					//db.SetHeader(host.Id, 0)
					//logger.WriteWork(host.Name + " нет ответа сервера")
					//sendler.Handler("🤔 " + host.Name + " нет ответа сервера! 🤔 ")
					logger.WriteWork(fmt.Sprintln(host.Name, err.Error()))
					log.Println(err)
					return
				}
			}
		}
	}
	if resp != nil {
		//defer resp.Body.Close()
	} else {
		return
	}

	if resp.StatusCode == 200 {

		if host.Header.Int64 != int64(resp.StatusCode) {
			db.SetHeader(host.Id, resp.StatusCode)
			wg.Add(1)
			sendler.Handler("✅ "+host.Name+" - "+resp.Status+" ✅ ", wg)
			logger.WriteWork(host.Name + " - " + resp.Status)
		}

		body, _ := io.ReadAll(resp.Body) //TODO перенести в проверку на пустоту

		loc, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			log.Println(err.Error())
			return
		}
		if ((time.Now().In(loc).Format("15:04") == "10:10" || time.Now().In(loc).Format("15:04") == "18:30") && (time.Now().In(loc).Weekday() != 0 && time.Now().In(loc).Weekday() != 6)) || (host.DomainTime.Int64 == 0) {
			wg.Add(1)
			go CheckDomain(host, wg) // проверяем дату аренды домена
		}
		wg.Add(1)
		go checkRedirect(host, resp, wg) // проверяем редирект на https
		wg.Add(1)
		go checkGtm(host, body, wg) // проверяем установленый GTM
		wg.Add(1)
		go CheckTemplate(host, body, wg) // проверяем верстку на ошибки и ключевые слова

		for _, v := range resp.Header.Values("Set-Cookie") {
			if strings.Contains(v, "october_session") {
				wg.Add(1)
				go CheckOctober(host, body, wg) // проверяем meta теги на сайтов на коробке
			}
		}

	} else {
		log.Println(host.Name, resp.StatusCode)

		if host.Header.Int64 != int64(resp.StatusCode) {
			db.SetHeader(host.Id, resp.StatusCode)
			wg.Add(1)
			sendler.Handler("‼️ 🆘 "+host.Name+" - "+strconv.Itoa(resp.StatusCode)+" 🆘 ‼️", wg)
			logger.WriteWork(host.Name + " - " + resp.Status)
		}
		return // прерываем проверку пока хост не станет 200
	}
}
