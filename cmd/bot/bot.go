package bot

import (
	"log"
	"os"
	"time"

	"github.com/ivm97/schedule_izh/cmd/sessions"
	"github.com/ivm97/schedule_izh/models"
)

// Зависимости приложения, здесь же можно разместить кэш с RW мьютексом
type application struct {
	settings *models.Configuration
	iLog     *log.Logger
	eLog     *log.Logger
	sess     *sessions.Session
}

func Entry(cfg models.Configuration) {
	file, err := os.OpenFile("logs/log.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		file = os.Stdout
	}
	iLog := log.New(file, "INFO\t", log.Ldate|log.Ltime)
	eLog := log.New(file, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	session := sessions.New(5*time.Minute, 10*time.Minute)
	app := &application{
		settings: &cfg,
		iLog:     iLog,
		eLog:     eLog,
		sess:     session,
	}

	app.start()
}
