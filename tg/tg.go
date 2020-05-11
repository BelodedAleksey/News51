package tg

import (
	"fmt"
	"net/http"
	"news/db"
	"news/log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//GetUpdates f
func GetUpdates() {
	http.DefaultTransport = transport

	bot, err := tgbotapi.NewBotAPI(*c.TGToken)
	if err != nil {
		log.LogRequestFile(fmt.Sprintf("[TELEGRAM]: %s", err))
		log.Fatalf("[TELEGRAM]: %s", err)
	}
	bot.Debug = false
	log.Infof("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.LogRequestFile(fmt.Sprintf("[TELEGRAM]: %s", err))
		log.Errorf("[TELEGRAM]: %s", err)
	}
	//routine check photos
	go checkPhotos(bot)

	//main cycle
	for update := range updates {
		//process callbacks
		if update.CallbackQuery != nil {
			callbackHandler(bot, update.CallbackQuery)
		}

		//process messages
		if update.Message != nil {
			messageHandler(bot, update.Message)
		}
	}

	http.DefaultTransport = &http.Transport{}
}

//collect photos
func checkPhotos(bot *tgbotapi.BotAPI) {
	for {
		select {
		case <-ticker.C:
			//if photos are present
			if len(photos) != 0 {
				lockPhotos.RLock()
				for k, v := range photos {
					var files []interface{}
					for _, p := range v {
						inpMedia := tgbotapi.NewInputMediaPhoto(p)
						inpMedia.Caption = caption
						files = append(files, inpMedia)
					}

					delete(photos, k)

					media := tgbotapi.NewMediaGroup(*c.TGAdmin, files)

					outMsg, err := bot.Send(media)
					if err != nil {
						log.LogRequestFile(fmt.Sprintf("[TELEGRAM]: %s", err))
						log.Errorf("[TELEGRAM]: %s", err)
						checkRetry(err)
					}

					user, err := db.GetEntry(
						&db.User{
							ChatID: strconv.FormatInt(k, 10),
						})
					if err != nil {
						log.LogRequestFile(fmt.Sprintf("[DB]: %s", err))
						log.Errorf("[DB]: %s", err)
					}
					queryID := postQuery(bot, user.Name)

					db.LockMessages.Lock()
					db.Messages[queryID] = &db.Post{
						User: *user, Time: outMsg.Time(), Media: &media,
					}
					db.LockMessages.Unlock()

					logData := fmt.Sprintf("User: %s suggest message, queryID: %d",
						user.Name, queryID)
					log.Infof(logData)
					log.LogRequestFile(logData)
				}
				lockPhotos.RUnlock()
			}
		default:
		}
	}
}
