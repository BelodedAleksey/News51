package hibiny

import (
	"fmt"
	"news/db"
	"news/log"
	"news/model"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

//Hibiny s
type Hibiny struct {
	lastUpdate string
}

//Init f
func (h *Hibiny) Init() {
	var err error
	h.lastUpdate, err = db.GetLastNew("hibiny")
	if err != nil {
		log.Errorf("[DB NEWS]: %s", err)
		log.LogRequestFile(fmt.Sprintf("[DB NEWS]: %s", err))
	}
}

//GetLastUpdate f
func (h *Hibiny) GetLastUpdate() string {
	return h.lastUpdate
}

//SetLastUpdate f
func (h *Hibiny) SetLastUpdate(u string) {
	err := db.UpdateLastNew("hibiny", u)
	if err != nil {
		log.Errorf("[DB NEWS]: %s", err)
		log.LogRequestFile(fmt.Sprintf("[DB NEWS]: %s", err))
	} else {
		h.lastUpdate = u
		log.Infof("[HIBINY LAST URL] %s", u)
	}
}

//GetDailyNews func
func (h *Hibiny) GetDailyNews() []model.New {
	year, month, day := time.Now().Date()
	date := fmt.Sprintf("?year=%d&month=%d&day=%d&", year, int(month), day)

	var news []model.New
	var stop bool
	n := model.New{}
	var url = "https://www.hibiny.com/news" + date
	c := colly.NewCollector()
	log.Infof("[HIBINY] Scrap started!")
	c.OnHTML(`*`, func(e *colly.HTMLElement) {
		/*if strings.Contains(e.Attr(`src`), `images/news`) && !stop {
			n.ImageURL = `https://www.hibiny.com` + e.Attr(`src`)
		}*/
		if strings.Contains(e.Attr(`href`), `/news/archive`) && e.Text != "" && !stop {
			n.URL = "https://www.hibiny.com" + e.Attr(`href`)
			if n.URL == h.lastUpdate {
				stop = true
				return
			}

			//n.Header = e.Text
			p := e.DOM.Parent()
			for i := 0; i < 7; i++ {
				p = p.Parent()
			}
			n.Date = p.Find(`td.p10`).Text()
			//n.Content = p.Find(`td.p`).Text()
			news = append(news, n)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		logData := fmt.Sprintf("[HIBINY SCRAP]: %s", err)
		log.Errorf(logData)
		log.LogRequestFile(logData)
	})

	err := c.Visit(url)
	if err != nil {
		logData := fmt.Sprintf("[HIBINY SCRAP] %s: %s", url, err)
		log.Errorf("[HIBINY SCRAP] %s: %s", url, err)
		log.LogRequestFile(logData)
	}
	c.Wait()
	return news
}
