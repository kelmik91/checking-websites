package main

import (
	"log"
	"targetPlus/internal/api/timeweb"
)

func main() {
	token := ""

	err, r := timeweb.Finance(token)
	if err != nil {
		log.Println(err.Error())
	}

	log.Println("Balance", r.Finances.Balance)
	log.Println("MonthlyCost", r.Finances.MonthlyCost)
	log.Println("HoursLeft", r.Finances.HoursLeft)
	log.Println("DayLeft", r.Finances.HoursLeft/24)
}
