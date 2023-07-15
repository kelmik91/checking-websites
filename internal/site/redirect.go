package site

import (
	"net/http"
	"strconv"
	"sync"
	"yandex/internal/db"
	"yandex/internal/send"
)

func checkRedirect(host db.Host, resp *http.Response, wg *sync.WaitGroup) {
	defer wg.Done()

	if resp.Request.Response != nil {
		if resp.Request.Response.StatusCode != 301 {
			if host.Redirect.String != "Код редиректа не 301. Код: "+strconv.Itoa(resp.Request.Response.StatusCode) {
				db.SetRedirect(host.Id, "Код редиректа не 301. Код: "+strconv.Itoa(resp.Request.Response.StatusCode))
				sendler.Handler(" 🔀 "+host.Name+" Код редиректа не 301. Код: "+strconv.Itoa(resp.Request.Response.StatusCode)+" 🔀 ", wg)
			}
			return
		} else if resp.Request.URL.String() != ("https://" + host.Name + "/") {
			if resp.Request.URL.String() == ("https://" + host.Name + ":443/") {
				if host.Redirect.String != "" {
					db.SetRedirect(host.Id, "")
					sendler.Handler(" 🔒 "+host.Name+" Редирект исправлен 🔒 ", wg)
				}
				return
			}
			if host.Redirect.String != "Редирект на другой адрес"+resp.Request.URL.String() {
				db.SetRedirect(host.Id, "Редирект на другой адрес"+resp.Request.URL.String())
				sendler.Handler(" 🆘🚫 "+host.Name+" Редирект на другой адрес "+resp.Request.URL.String()+" 🚫🆘 ", wg)
			}
			return
		} else {
			if host.Redirect.String != "" {
				db.SetRedirect(host.Id, "")
				sendler.Handler(" 🔒 "+host.Name+" Редирект исправлен 🔒 ", wg)
			}
			return
		}
	} else {
		if host.Redirect.String != "Нет редиректа на HTTPS" {
			db.SetRedirect(host.Id, "Нет редиректа на HTTPS")
			sendler.Handler(" 🛂 "+host.Name+" Нет редиректа на HTTPS 🛂 ", wg)
		}
	}
}
