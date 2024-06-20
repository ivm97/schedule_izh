package models

type SearchData struct {
	Dep      string
	Arr      string
	From     string
	To       string
	DepYear  string
	DepMonth string
	DepDay   string
	ArrYear  string
	ArrMonth string
	ArrDay   string
	Keyboard any
}

type ScheduleFlights struct {
	Flight []Flight
}

type Flight struct {
	OriginportName      string   `json:"originportName"`
	EndDate             string   `json:"end_date"`
	DestinationCity     string   `json:"destinationcity"`
	DestinationCityName string   `json:"destinationcityName"`
	BeginDate           string   `json:"begin_date"`
	Classes             []string `json:"classes"`
	DeparturedayShift   string   `json:"departuredayshift"`
	OriginCity          string   `json:"origincity"`
	VehicleCodeEn       string   `json:"vehicleCodeEn"`
	DepartureTime       string   `json:"departuretime"`
	DestinationPort     string   `json:"destinationport"`
	OrigincityName      string   `json:"origincityName"`
	ArrivalDayShift     string   `json:"arrivaldayshift"`
	Airplane            string   `json:"airplane"`
	ArrivalTime         string   `json:"arrivaltime"`
	OriginPort          string   `json:"originport"`
	Days                string   `json:"days"`
	Comp                Company  `json:"company"`
	FlightTime          string   `json:"flightTime"`
	DestinationPortName string   `json:"destinationportName"`
	RaceNumber          string   `json:"racenumber"`
}

type Company struct {
	Code   string `json:"code"`
	CodeEn string `json:"code_en"`
	Name   string `json:"name"`
	NameEn string `json:"name_en"`
}

type Configuration struct {
	Bot     BotSettings     `json:"bot_settings"`
	Parser  ParserSettings  `json:"parser_settings"`
	Storage StorageSettings `json:"storage_settings"`
}

type BotSettings struct {
	Token string `json:"token"`
}

type ParserSettings struct {
	StartUrl string `json:"start_url"`
}

type StorageSettings struct {
	DName string `json:"dname"`
	DUser string `json:"duser"`
	DPass string `json:"dpass"`
	DPort string `json:"dport"`
	Dhost string `json:"dhost"`
	Try   int    `json:"try"`
}

type User struct {
	Id      int
	TgId    int
	Name    string
	Records []Recs
}

type Recs struct {
	TgId       int
	StartPoint string
	EndPoint   string
	DateStart  string
	DateEnd    string
}
