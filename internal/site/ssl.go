package site

import (
	"crypto/tls"
	"log"
	"sync"
	"time"
	"yandex/internal/db"
	"yandex/internal/logger"
	"yandex/internal/send"
)

func Ssl(host db.Host, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := tls.Dial("tcp", host.Name+":443", nil)
	if err != nil {

		conn, err = tls.Dial("tcp", host.Name+":443", nil)
		if err != nil {
			logger.WriteWork(err.Error())
			log.Println(err.Error())
			if host.SslNotification.Bool == false {
				wg.Add(1)
				sendler.Handler("‼️☠️‼️ "+host.Name+" не удалось получить SSL ‼️☠️‼️", wg)
				db.SetSslNotification(host.Id, true)
			}
			return
		}
	}
	defer conn.Close()

	state := conn.ConnectionState()
	for _, cert := range state.PeerCertificates {
		if cert.IsCA {
			continue
		}

		loc, _ := time.LoadLocation("Europe/Moscow")
		if time.Now().Add(time.Hour*96).Unix() >= cert.NotAfter.In(loc).Unix() && !host.SslNotification.Bool {
			db.SetSslNotification(host.Id, true)
			//todo если осталось меньше 4 дней, то отправляем в бота уведомление
			//fmt.Println(host, cert.NotAfter.In(loc).Format(time.RFC3339))
			wg.Add(1)
			sendler.Handler("⚠️ "+host.Name+" дата окончания SSL "+cert.NotAfter.In(loc).Format("02.01.2006 15:04")+" ⚠️", wg)
		} else if host.SslNotification.Bool {
			db.SetSslNotification(host.Id, false)
		}

		if host.SslTime.Int64 != cert.NotAfter.In(loc).Unix() {
			db.SetSslTime(host.Id, cert.NotAfter.In(loc).Unix())
			wg.Add(1)
			sendler.Handler("🆕 "+host.Name+" новый SSL "+cert.NotAfter.In(loc).Format("02.01.2006 15:04")+" 🆕", wg)
		}
	}
}
