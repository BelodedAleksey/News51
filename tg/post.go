package tg

import (
	"fmt"
	"net/http"
	"news/db"
	"news/log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//SendPost f
func SendPost(post *db.Post, bot *tgbotapi.BotAPI) {
	log.Infof("SendPost!")
	if post.Message != nil && post.Message.Text != "" {
		text := ""
		if post.User.Anon {
			text = fmt.Sprintf("%s\n_Прислал наш подписчик в_ @MurmanNewsBot", post.Message.Text)
		} else {
			text = fmt.Sprintf(
				"%s\n_Прислал наш подписчик @%s в_ @MurmanNewsBot", post.Message.Text, post.User.Name)
		}

		msg := tgbotapi.NewMessageToChannel(*c.TGChannel, text)
		msg.ParseMode = "markdown"

		_, err := bot.Send(msg)
		if err != nil {
			log.LogRequestFile(fmt.Sprintf("[TELEGRAM]: %s", err))
			log.Errorf("[TELEGRAM]: %s", err)
			checkRetry(err)
		}
	}

	if post.Media != nil {
		media := post.Media
		media.ChatID = *c.TGChatID
		_, err := bot.Send(post.Media)
		if err != nil {
			log.LogRequestFile(fmt.Sprintf("[TELEGRAM]: %s", err))
			log.Errorf("[TELEGRAM]: %s", err)
			checkRetry(err)
		}
	}

}

//SendMessage f
func SendMessage(message, pic string) {
	http.DefaultTransport = transport

	bot, err := tgbotapi.NewBotAPI(*c.TGToken)
	if err != nil {
		log.LogRequestFile(fmt.Sprintf("[TELEGRAM]: %s", err))
		log.Fatalf("[TELEGRAM]: %s", err)
	}
	bot.Debug = false
	log.Infof("Authorized on account %s", bot.Self.UserName)

	//post picture
	if pic != "" {
		photo := tgbotapi.NewPhotoUpload(*c.TGChatID, nil)
		photo.FileID = pic
		photo.UseExisting = true
		//photo.Caption = message

		_, err = bot.Send(photo)
		if err != nil {
			log.LogRequestFile(fmt.Sprintf("[TELEGRAM]: %s", err))
			log.Errorf("[TELEGRAM]: %s", err)
			checkRetry(err)
		}
	}

	msg := tgbotapi.NewMessageToChannel(*c.TGChannel, message)
	msg.ParseMode = "markdown"

	_, err = bot.Send(msg)
	if err != nil {
		log.LogRequestFile(fmt.Sprintf("[TELEGRAM]: %s", err))
		log.Errorf("[TELEGRAM]: %s", err)
		checkRetry(err)
	}
	http.DefaultTransport = &http.Transport{}
}

//check if pause for retry is needed
func checkRetry(err error) {
	if strings.Contains(err.Error(), "retry after") {
		spl := strings.Split(err.Error(), "retry after ")
		pause, err := strconv.Atoi(spl[len(spl)-1])
		if err != nil {
			log.Errorf("[STRCONV]: %s", err)
		}
		time.Sleep(time.Duration(pause) * time.Second)
	}
}
