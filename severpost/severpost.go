package severpost

import (
	"fmt"
	"news/log"
	"strings"
	"time"

	"news/model"

	"github.com/gocolly/colly"
)

//Severpost s
type Severpost struct {
	lastUpdate string
}

//GetLastUpdate f
func (s *Severpost) GetLastUpdate() string {
	return s.lastUpdate
}

//SetLastUpdate f
func (s *Severpost) SetLastUpdate(u string) {
	s.lastUpdate = u
	log.Infof("[SEVERPOST LAST URL] %s", u)
}

//GetDailyNews func
func (s *Severpost) GetDailyNews() []model.New {
	var news []model.New
	var stop bool
	n := model.New{}
	var url = "https://severpost.ru/archive/"
	c := colly.NewCollector()

	t := time.Now()
	date := t.Format("02 01 2006")

	log.Infof("[SEVERPOST] Scrap started!")
	c.OnHTML(`*`, func(e *colly.HTMLElement) {
		//URL & HEADER
		if e.Name == "h2" && e.Attr("class") != "m-3" {
			if strings.Contains(e.ChildAttr("a", "href"), "/read/") && e.Text != "" && !stop {
				n.URL = "https://severpost.ru" + e.ChildAttr("a", "href")
				if n.URL == s.lastUpdate {
					stop = true
					return
				}
				//n.Header = e.Text
			}
		}

		//Date
		if e.Name == "span" && e.Attr("class") == "e-datetime" {
			text := strings.Split(e.Text, "|")
			if len(text) > 1 {
				n.Date = strings.TrimSpace(text[1])
				postDate := strings.Split(n.Date, ",")[0]
				if date != postDate {
					stop = true
					return
				}
				news = append(news, n)
			}
		}

		//Image URL
		/*if e.Name == "a" &&
			strings.HasPrefix(e.Attr("href"), "/read/") &&
			e.Attr("target") != "_blank" &&
			e.Attr("class") == "" &&
			!stop {
			if strings.HasPrefix(e.ChildAttr("img", "src"), "/docs/upload/cache/") {
				n.ImageURL = "https://severpost.ru" + e.ChildAttr("img", "src")
			}
			//Content next <p>
			p := e.DOM.Next()
			if p.Is("p") {
				n.Content = p.Text()
				news = append(news, n)
			}
		}*/
	})

	c.OnError(func(r *colly.Response, err error) {
		logData := fmt.Sprintf("[SEVERPOST SCRAP]: %s", err)
		log.Errorf(logData)
		log.LogRequestFile(logData)
	})

	err := c.Visit(url)
	if err != nil {
		logData := fmt.Sprintf("[SEVERPOST SCRAP] %s: %s", url, err)
		log.Errorf("[SEVERPOST SCRAP] %s: %s", url, err)
		log.LogRequestFile(logData)
	}
	c.Wait()
	return news
}
