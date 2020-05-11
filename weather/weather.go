package weather

import (
	"fmt"
	"news/log"
	"strings"

	"github.com/gocolly/colly"
)

//Weather s
type Weather struct {
	Time          string
	Temp          string
	Precipitation string
	Pressure      string
	Wind          string
	Humidity      string
	Cloudness     string
}

//GetWeather f
func GetWeather() []Weather {
	var stop bool
	weather := Weather{}
	var cast []Weather
	var url = "https://pogoda51.ru/apatity"
	c := colly.NewCollector()

	log.Infof("[POGODA51] Scrap started!")
	c.OnHTML(`.nordui-table-td`, func(e *colly.HTMLElement) {
		val, ok := e.DOM.Parent().Attr("class")
		if val == "nordui-table-subheader" &&
			ok &&
			!stop &&
			!strings.Contains(e.Text, "Сегодня") {
			stop = true
			return
		}

		if strings.Contains(e.Text, ":") &&
			strings.Contains(e.Attr("style"), "center") &&
			!stop {
			weather.Time = e.Text
			precipitation := e.DOM.Next().Next().Next()
			weather.Precipitation = precipitation.Text()
			pressure := precipitation.Next()
			weather.Pressure = pressure.Text()
			wind := pressure.Next()
			weather.Wind = wind.Text()
			humidity := wind.Next()
			weather.Humidity = humidity.Text()
			cloudness := humidity.Next()
			weather.Cloudness = cloudness.Text()
		}

		if e.ChildAttr("img", "src") != "" && !stop {
			imgURL := strings.Split(e.ChildAttr("img", "src"), "/")
			temp := strings.Replace(
				strings.TrimSuffix(
					imgURL[len(imgURL)-1], ".png"),
				"n", "-1", 1) + "°C"
			weather.Temp = temp
			cast = append(cast, weather)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		logData := fmt.Sprintf("[POGODA51 SCRAP]: %s", err)
		log.LogRequestFile(logData)
		log.Errorf(logData)
	})

	err := c.Visit(url)
	if err != nil {
		logData := fmt.Sprintf("[POGODA51 SCRAP] %s: %s", url, err)
		log.LogRequestFile(logData)
		log.Errorf(logData)
	}
	c.Wait()
	return cast
}
