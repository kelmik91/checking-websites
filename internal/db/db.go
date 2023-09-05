package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strconv"
	"targetPlus/internal/logger"
)

type Host struct {
	Id                 int            `json:"id"`
	Name               string         `json:"name"`
	Header             sql.NullInt64  `json:"header,omitempty"`
	DomainTime         sql.NullInt64  `json:"domain_time,omitempty"`
	DomainNotification sql.NullBool   `json:"domain_notification,omitempty"`
	SslTime            sql.NullInt64  `json:"ssl_time,omitempty"`
	SslNotification    sql.NullBool   `json:"ssl_notification,omitempty"`
	Gtm                sql.NullString `json:"gtm,omitempty"`
	GtmVeryfi          sql.NullBool   `json:"gtm_verify,omitempty"`
	October            sql.NullString `json:"header_october,omitempty"`
	TemplateError      sql.NullString `json:"template_error,omitempty"`
	Redirect           sql.NullString `json:"redirect,omitempty"`
}

func getConn() *sql.DB {

	hostDB := os.Getenv("host")
	loginDB := os.Getenv("loginDB")
	passwordDB := os.Getenv("passwordDB")
	nameDB := os.Getenv("nameDB")
	portDB := os.Getenv("portDB")

	dataSourceName := fmt.Sprint(loginDB, ":", passwordDB, "@tcp(", hostDB, ":", portDB, ")/", nameDB)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		logger.WriteWork(fmt.Sprintf(err.Error()))
		panic(err.Error())
	}
	return db
}

func GetAllUser() []int {
	db := getConn()
	defer db.Close()

	rows, err := db.Query("SELECT id FROM users WHERE is_active = 1")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var users []int
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			panic(err.Error())
		}
		users = append(users, id)
	}

	return users
}

func GetDigitalUser() []int {
	db := getConn()
	defer db.Close()

	rows, err := db.Query("SELECT id FROM users WHERE is_active = 1 AND role = 2")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var users []int
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			panic(err.Error())
		}
		users = append(users, id)
	}

	return users
}

func GetHosts() map[int]Host {
	db := getConn()
	defer db.Close()

	rows, err := db.Query("SELECT id, name, header, ssl_time, ssl_notification, gtm, gtm_verify, header_october, template_error, domain_time, domain_notification, redirect FROM domains WHERE is_active = 1")
	if err != nil {
		//logger.WriteWork(err.Error())
		panic(err.Error())
	}
	defer rows.Close()

	var hosts = make(map[int]Host)
	for rows.Next() {
		host := Host{}
		err = rows.Scan(&host.Id, &host.Name, &host.Header, &host.SslTime, &host.SslNotification, &host.Gtm, &host.GtmVeryfi, &host.October, &host.TemplateError, &host.DomainTime, &host.DomainNotification, &host.Redirect)
		if err != nil {
			//logger.WriteWork(err.Error())
			panic(err.Error())
		}
		hosts[host.Id] = host
	}

	return hosts
}

func GetAllHostsCompany() map[int]Host {
	db := getConn()
	defer db.Close()

	rows, err := db.Query("SELECT id, name, domain_time, domain_notification FROM domainsAll")
	if err != nil {
		//logger.WriteWork(err.Error())
		panic(err.Error())
	}
	defer rows.Close()

	var hosts = make(map[int]Host)
	for rows.Next() {
		host := Host{}
		err = rows.Scan(&host.Id, &host.Name, &host.DomainTime, &host.DomainNotification)
		if err != nil {
			//logger.WriteWork(err.Error())
			panic(err.Error())
		}
		hosts[host.Id] = host
	}

	return hosts
}

func GetDisableHosts() map[int]Host {
	db := getConn()
	defer db.Close()

	rows, err := db.Query("SELECT id, name, header, domain_time, domain_notification FROM domains WHERE is_active = 0")
	if err != nil {
		//logger.WriteWork(err.Error())
		panic(err.Error())
	}
	defer rows.Close()

	var hosts = make(map[int]Host)
	for rows.Next() {
		host := Host{}
		err = rows.Scan(&host.Id, &host.Name, &host.Header, &host.DomainTime, &host.DomainNotification)
		if err != nil {
			//logger.WriteWork(err.Error())
			panic(err.Error())
		}
		hosts[host.Id] = host
	}

	return hosts
}

func GetGTM() map[int]string {
	db := getConn()
	defer db.Close()

	rows, err := db.Query("SELECT id, value FROM gtms")
	if err != nil {
		//logger.WriteWork(err.Error())
		panic(err.Error())
	}
	defer rows.Close()

	var gtm = make(map[int]string)
	for rows.Next() {
		var id int
		var value string
		err = rows.Scan(&id, &value)
		if err != nil {
			//logger.WriteWork(err.Error())
			panic(err.Error())
		}
		gtm[id] = value
	}

	return gtm
}

func GetGtmByDomain(domainId int) string {
	db := getConn()
	defer db.Close()

	rows, err := db.Query("select gtms.value from domains_gtms_rels, gtms where domains_gtms_rels.gtm_id = gtms.id and domains_gtms_rels.domain_id = ?", domainId)
	if err != nil {
		//logger.WriteWork(err.Error())
		panic(err.Error())
	}
	defer rows.Close()

	var gtm string
	for rows.Next() {
		err = rows.Scan(&gtm)
		if err != nil {
			//logger.WriteWork(err.Error())
			panic(err.Error())
		}
	}

	return gtm
}

func SetHeader(id int, statusCode int) {
	db := getConn()
	defer db.Close()

	err := db.QueryRow("UPDATE domains SET header = ? WHERE id = ?", strconv.Itoa(int(statusCode)), id)
	if err.Err() != nil {
		//logger.WriteWork(err.Err().Error())
		panic(err.Err().Error())
	}
}

func SetGTM(id int, gtm string) {
	db := getConn()
	defer db.Close()

	err := db.QueryRow("UPDATE domains SET gtm = ? WHERE id = ?", gtm, id)
	if err.Err() != nil {
		//logger.WriteWork(err.Err().Error())
		panic(err.Err().Error())
	}
}

func SetGtmVerify(id int, verify bool) {
	db := getConn()
	defer db.Close()

	err := db.QueryRow("UPDATE domains SET gtm_verify = ? WHERE id = ?", verify, id)
	if err.Err() != nil {
		//logger.WriteWork(err.Err().Error())
		panic(err.Err().Error())
	}
}

func SetSslTime(id int, sslTime int64) {
	db := getConn()
	defer db.Close()

	err := db.QueryRow("UPDATE domains SET ssl_time = ? WHERE id = ?", sslTime, id)
	if err.Err() != nil {
		//logger.WriteWork(err.Err().Error())
		panic(err.Err().Error())
	}
}

func SetSslNotification(id int, sslNotification bool) {
	db := getConn()
	defer db.Close()

	err := db.QueryRow("UPDATE domains SET ssl_notification = ? WHERE id = ?", sslNotification, id)
	if err.Err() != nil {
		//logger.WriteWork(err.Err().Error())
		panic(err.Err().Error())
	}
}

func SetTemplateError(id int, templateError string) {
	db := getConn()
	defer db.Close()

	var err *sql.Row
	if templateError == "" {
		err = db.QueryRow("UPDATE domains SET template_error = ? WHERE id = ?", nil, id)
	} else {
		err = db.QueryRow("UPDATE domains SET template_error = ? WHERE id = ?", templateError, id)
	}

	if err.Err() != nil {
		//logger.WriteWork(err.Err().Error())
		panic(err.Err().Error())
	}
}

func SetOctober(id int, value string) {
	db := getConn()
	defer db.Close()

	var err *sql.Row
	if value == "" {
		err = db.QueryRow("UPDATE domains SET header_october = ? WHERE id = ?", nil, id)
	} else {
		err = db.QueryRow("UPDATE domains SET header_october = ? WHERE id = ?", value, id)
	}
	if err.Err() != nil {
		//logger.WriteWork(err.Err().Error())
		panic(err.Err().Error())
	}
}

func SetDomainTime(id int, domainTime int64) {
	db := getConn()
	defer db.Close()

	err := db.QueryRow("UPDATE domains SET domain_time = ? WHERE id = ?", domainTime, id)
	if err.Err() != nil {
		//logger.WriteWork(err.Err().Error())
		panic(err.Err().Error())
	}
}

func SetDomainsCompanyTime(id int, domainTime int64) {
	db := getConn()
	defer db.Close()

	err := db.QueryRow("UPDATE domainsAll SET domain_time = ? WHERE id = ?", domainTime, id)
	if err.Err() != nil {
		//logger.WriteWork(err.Err().Error())
		panic(err.Err().Error())
	}
}

func SetDomainNotification(id int, domainNotification bool) {
	db := getConn()
	defer db.Close()

	err := db.QueryRow("UPDATE domains SET domain_notification = ? WHERE id = ?", domainNotification, id)
	if err.Err() != nil {
		//logger.WriteWork(err.Err().Error())
		panic(err.Err().Error())
	}
}

func SetRedirect(id int, status string) {
	db := getConn()
	defer db.Close()

	var err *sql.Row
	if status == "" {
		err = db.QueryRow("UPDATE domains SET redirect = ? WHERE id = ?", nil, id)
	} else {
		err = db.QueryRow("UPDATE domains SET redirect = ? WHERE id = ?", status, id)
	}
	if err.Err() != nil {
		//logger.WriteWork(err.Err().Error())
		panic(err.Err().Error())
	}
}
