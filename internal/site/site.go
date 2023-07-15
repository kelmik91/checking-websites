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
			time.Sleep(time.Second * 15)
			resp, err = client.Get("https://" + host.Name)
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

		//if time.Now().Minute()%15 == 0 || host.SslNotification.Bool == true || host.SslTime.Int64 == 0 {
		//	wg.Add(1)
		//	go Ssl(host, wg)
		//}

		if time.Now().Minute() == 45 || host.DomainTime.Int64 == 0 {
			wg.Add(1)
			go CheckDomain(host, wg)
		}
		wg.Add(1)
		go checkRedirect(host, resp, wg)
		wg.Add(1)
		go checkGtm(host, body, wg)
		wg.Add(1)
		go CheckTemplate(host, body, wg)

		for _, v := range resp.Header.Values("Set-Cookie") {
			if strings.Contains(v, "october_session") {
				wg.Add(1)
				go CheckOctober(host, body, wg)
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
