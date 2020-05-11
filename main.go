package main

import (
	"flag"
	"fmt"
	"time"

	"news/db"
	"news/hibinform"
	"news/hibiny"
	"news/log"
	"news/model"
	"news/severpost"
	"news/tg"
	"news/weather"
)

var (
	sites = []NewsSite{
		&hibiny.Hibiny{},
		&hibinform.Hibinform{},
		&severpost.Severpost{},
	}

	//Flags
	_ = flag.Bool("bg", false, "background process")
	_ = flag.Bool("socks", false, "socks5")
)

//NewsSite t
type NewsSite interface {
	GetDailyNews() []model.New
	GetLastUpdate() string
	SetLastUpdate(string)
}

//Func for flag bg
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

//PostNews f
func PostNews(news []model.New, s NewsSite) {
	log.Infof("[UPDATE]: amount of news %d", len(news))
	//check news are not empty
	if len(news) > 0 {
		//reverse order
		for i := len(news) - 1; i >= 0; i-- {
			//check new news
			n := news[i]
			if n.URL == s.GetLastUpdate() {
				break
			}
			/*message := fmt.Sprintf(
				"*%s*\n_%s_\n%s\n%s\n",
				n.Header,
				n.Date,
				n.Content,
				n.URL,
			)
			tg.SendMessage(message, n.ImageURL)*/
			message := fmt.Sprintf(
				"*%s*\n%s\n",
				n.Date,
				n.URL,
			)
			tg.SendMessage(message, "")
		}
		//write last update new
		s.SetLastUpdate(news[0].URL)
	}
}

//UpdateNews f
func UpdateNews() {
	ticker := time.NewTicker(time.Minute)
	wPosted := false

	for {
		select {
		case t := <-ticker.C:
			log.Infof("[UPDATE]: Started")

			//Weather
			h, m, _ := t.Clock()
			if h == 6 && m > 30 && !wPosted {
				PostWeather()
				wPosted = true
			}
			//reset
			if h == 0 && m > 30 {
				wPosted = false
			}

			//News
			for _, s := range sites {
				news := s.GetDailyNews()
				PostNews(news, s)
			}
		default:
		}
	}
}

//PostWeather f
func PostWeather() {
	cast := weather.GetWeather()
	message := "*Прогноз погоды на сегодня*\n"
	for _, w := range cast {
		message += fmt.Sprintf(
			"*%s* 🌡 %s 🌧︎ %s 💨 %s 💧 %s ☁ %s P %s\n",
			w.Time,
			w.Temp,
			w.Precipitation,
			w.Wind,
			w.Humidity,
			w.Cloudness,
			w.Pressure,
		)
	}
	tg.SendMessage(message, "")
}

func main() {
	flag.Parse()
	isBackground := isFlagPassed("bg")
	withSocks := isFlagPassed("socks")
	//Init logging
	log.Init(true, isBackground, "log.txt")
	//Init telegram
	tg.Init(withSocks)

	//Init db
	db.Init()
	//Main func
	go UpdateNews()

	tg.GetUpdates()
}
