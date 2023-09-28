package timeweb

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

type RespSuccess struct {
	Finances struct {
		Balance           float64     `json:"balance"`
		Currency          string      `json:"currency"`
		DiscountEndDateAt interface{} `json:"discount_end_date_at"`
		DiscountPercent   int         `json:"discount_percent"`
		HourlyCost        float64     `json:"hourly_cost"`
		HourlyFee         float64     `json:"hourly_fee"`
		MonthlyCost       int         `json:"monthly_cost"`
		MonthlyFee        int         `json:"monthly_fee"`
		TotalPaid         float64     `json:"total_paid"`
		HoursLeft         int         `json:"hours_left"`
		AutopayCardInfo   interface{} `json:"autopay_card_info"`
	} `json:"finances"`
	ResponseId string `json:"response_id"`
}

func Finance(TimewebCloudToken string) (error, RespSuccess) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.timeweb.cloud/api/v1/account/finances", nil)
	if err != nil {
		log.Fatal(err)
		return err, RespSuccess{}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+TimewebCloudToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return err, RespSuccess{}
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		log.Println("StatusCode", resp.StatusCode)
		return errors.New(resp.Status), RespSuccess{}
	}

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%s\n", bodyText)
	var respJ RespSuccess
	err = json.Unmarshal(bodyText, &respJ)
	if err != nil {
		return nil, RespSuccess{}
	}

	return nil, respJ
}
