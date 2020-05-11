package hibinform

import (
	"fmt"
	"news/log"
	"strings"
	"time"

	"news/model"

	"github.com/gocolly/colly"
)

//Hibinform s
type Hibinform struct {
	lastUpdate string
}

var (
	months = []string{
		"Янв", "Фев", "Мар", "Апр", "Май", "Июн", "Июл", "Авг", "Сен", "Окт", "Ноя", "Дек",
	}
)

//GetLastUpdate f
func (h *Hibinform) GetLastUpdate() string {
	return h.lastUpdate
}

//SetLastUpdate f
func (h *Hibinform) SetLastUpdate(u string) {
	h.lastUpdate = u
	log.Infof("[HIBINFORM LAST URL] %s", u)
}

//GetDailyNews func
func (h *Hibinform) GetDailyNews() []model.New {
	var news []model.New
	var stop bool
	n := model.New{}
	var url = "https://hibinform.ru/"
	c := colly.NewCollector()

	t := time.Now()
	arrTime := strings.Split(t.Format("02 01 2006"), " ")
	_, m, _ := t.Date()
	date := fmt.Sprintf("%s %s %s", arrTime[0], months[int(m)-1], arrTime[2])

	log.Infof("[HIBINFORM] Scrap started!")

	c.OnHTML(`*`, func(e *colly.HTMLElement) {
		//Header & ImageURL
		if strings.Contains(e.Attr("class"), "entry-header") &&
			!stop {
			n.URL = e.ChildAttr("a", "href")
			if n.URL == h.lastUpdate {
				stop = true
				return
			}
			n.Header = e.ChildText("a")
			n.ImageURL = e.ChildAttr("img", "src")
		}

		//Date
		if strings.Contains(e.Attr("class"), "publish-date") {
			n.Date = strings.TrimSpace(e.Text)
			if n.Date != date {
				stop = true
				return
			}
		}
		//Time
		if strings.Contains(e.Attr("class"), "publish-time") &&
			!stop {
			n.Date += " " + strings.TrimSpace(e.Text)
		}
		//Content
		if strings.Contains(e.Attr("class"), "entry-content") &&
			!stop {
			n.Content = e.ChildText("p")
			news = append(news, n)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		logData := fmt.Sprintf("[HIBINFORM SCRAP]: %s", err)
		log.Errorf(logData)
		log.LogRequestFile(logData)
	})

	err := c.Visit(url)
	if err != nil {
		logData := fmt.Sprintf("[HIBINFORM SCRAP] %s: %s", url, err)
		log.Errorf(logData)
		log.LogRequestFile(logData)
	}
	c.Wait()
	return news
}
