package tg

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"news/log"
	"os"
	"sync"
	"time"

	"golang.org/x/net/proxy"
	"gopkg.in/yaml.v2"
)

var (
	transport = &http.Transport{}
	c         = Config{}

	//chatID key
	photos     = map[int64][]string{}
	lockPhotos = sync.RWMutex{}

	caption string

	ticker = time.NewTicker(3 * time.Second)
)

//Config s
type Config struct {
	TGToken     *string `yaml:"tgToken"`
	TGChannel   *string `yaml:"tgChannel"`
	TGAdmin     *int64  `yaml:"tgAdmin"`
	TGChatID    *int64  `yaml:"tgChatID"`
	TGSocks     *string `yaml:"tgSocks"`
	TGSocksUser *string `yaml:"tgSocksUser"`
	TGSocksPass *string `yaml:"tgSocksPass"`
}

//Init func
func Init(withSocks bool) {
	//Parse config
	f, err := os.Open(`config.yaml`)
	defer f.Close()
	if err != nil {
		log.LogRequestFile(fmt.Sprintf("[YAML]: %s", err))
		log.Errorf("[YAML configuration] open: %s", err)
	}

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		log.LogRequestFile(fmt.Sprintf("[YAML]: %s", err))
		log.Errorf("[YAML configuration] read: %s", err)
	}

	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		log.LogRequestFile(fmt.Sprintf("[YAML]: %s", err))
		log.Errorf("[YAML configuration] unmarshalling: %s", err)
	}

	//Proxy for tg
	if withSocks {
		auth := &proxy.Auth{User: *c.TGSocksUser, Password: *c.TGSocksPass}
		dialer, err := proxy.SOCKS5("tcp", *c.TGSocks, auth, proxy.Direct)
		if err != nil {
			log.LogRequestFile(fmt.Sprintf("[SOCKS]: %s", err))
			log.Fatalf("[SOCKS] %s", err)
		}

		transport = &http.Transport{
			Dial: dialer.Dial,
		}
	}
}
