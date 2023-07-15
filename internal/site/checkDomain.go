package site

import (
	"log"
	"strings"
	"sync"
	"time"
	"yandex/internal/db"
	"yandex/internal/logger"
	"yandex/internal/send"
)

func CheckDomain(host db.Host, wg *sync.WaitGroup) {
	defer wg.Done()

	whois, err := Whois(host.Name)
	if err != nil {
		wg.Add(1)
		sendler.Handler(host.Name+" не удалось получить WHOIS", wg)
		logger.WriteWork(host.Name + " не удалось получить WHOIS\n" + err.Error())
		log.Println(err.Error())
		return
	}

	date := time.Time{}
	if strings.Contains(whois, "paid-till:") {
		startString := "paid-till:"
		date, err = findDate(whois, startString)
		if err != nil {
			panic(err)
		}
	}
	if strings.Contains(whois, "Registry Expiry Date:") {
		startString := "Registry Expiry Date:"
		date, err = findDate(whois, startString)
		if err != nil {
			panic(err)
		}
	}
	if strings.Contains(whois, "Registrar Registration Expiration Date:") {
		startString := "Registrar Registration Expiration Date:"
		date, err = findDate(whois, startString)
		if err != nil {
			panic(err)
		}
	}

	loc, _ := time.LoadLocation("Europe/Moscow")

	if date.Unix() != host.DomainTime.Int64 {

		db.SetDomainTime(host.Id, date.Unix())
		db.SetDomainNotification(host.Id, false)

		dateEnd := date.In(loc).Format("02/01/2006 15:04")
		wg.Add(1)
		sendler.Handler(host.Name+" аренда домена закончится - "+dateEnd, wg)

	} else if time.Now().Add(time.Hour*720).Unix() >= date.In(loc).Unix() && !host.DomainNotification.Bool {

		dateEnd := date.In(loc).Format("02/01/2006 15:04")
		log.Println(host.Name, dateEnd)

		db.SetDomainTime(host.Id, date.Unix())
		db.SetDomainNotification(host.Id, true)

		wg.Add(1)
		sendler.Handler(host.Name+" аренда домена закончится - "+dateEnd, wg)
	}
}

func findDate(s, startString string) (time.Time, error) {
	startIndex := strings.Index(s, startString) + len(startString)
	endIndex := strings.Index(s[startIndex:], "\n") + startIndex
	date := s[startIndex:endIndex]
	date = strings.TrimSpace(date)
	parse, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return parse, err
	}
	return parse, nil
}
