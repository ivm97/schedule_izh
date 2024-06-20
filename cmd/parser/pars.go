package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ivm97/schedule_izh/models"
)

type parsapp struct {
	models.ParserSettings
}

//https://booking.izhavia.su/websky/json/get-schedule-period?departure=IJK&arrival=MOW&fromDate=13.06.2024&toDate=25.07.2024

func get(url string, src *models.SearchData) []byte {
	//Оказывается апи доступно извне, даже токен не требуется :(
	getReq := fmt.Sprintf("departure=%s&arrival=%s&fromDate=%s&toDate=%s", src.Dep, src.Arr, src.From, src.To)
	startUrl := "booking.izhavia.su/websky/json/get-schedule-period"

	resp, err := http.Get(getReq + startUrl)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	bytes, _ := io.ReadAll(resp.Body)
	return bytes
}

func convert(b []byte) *models.ScheduleFlights {
	var recs models.ScheduleFlights
	err := json.Unmarshal(b, &recs)
	if err != nil {
		log.Println(err)
	}

	return &recs

}

/*
func (p *parsapp) Parser(src *models.SearchData) {
	//startUrl := "booking.izhavia.su/websky/json/get-schedule-period"
	b := get(p.StartUrl, src)
	result := convert(b)
}
*/
