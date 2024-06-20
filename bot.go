package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/ivm97/schedule_izh/cmd/bot"
	"github.com/ivm97/schedule_izh/models"
)

func main() {
	cfg := readConf("settings/config.json")
	bot.Entry(cfg)

}

// Читаем и передаем главный конфиг
func readConf(path string) models.Configuration {
	data, err := os.ReadFile(path)
	if err != nil {
		//В идеале хорошо бы написать кастомный логгер и конкурентно собирать ошибки
		//в файл или же в бд
		log.Println(err)
	}

	var cfg models.Configuration
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Println(err, "data: ", cfg)
	}
	return cfg
}
