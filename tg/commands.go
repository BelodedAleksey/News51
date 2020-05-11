package tg

import (
	"fmt"
	"news/db"
	"news/log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//post Accept/Decline message with query
func postQuery(bot *tgbotapi.BotAPI, username string) int {
	messageToSend := fmt.Sprintf("User: %s\n *Запостить сообщение?*", username)
	id := *c.TGAdmin
	msg := tgbotapi.NewMessage(id, messageToSend)
	msg.ParseMode = "markdown"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("Да", "Accept"),
			tgbotapi.NewInlineKeyboardButtonData("Нет", "Decline"),
		},
	)

	outMsg, err := bot.Send(msg)
	if err != nil {
		log.LogRequestFile(fmt.Sprintf("[TELEGRAM]: %s", err))
		log.Errorf("[TELEGRAM]: %s", err)
		checkRetry(err)
	}
	return outMsg.MessageID
}

// cmd /start
func startCmd(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	messageToSend := fmt.Sprintf(
		`Привет, %s! Чтобы отправить сообщение,
		 сначала введи команду _/anonymous_ и выбери как бот будет получать посты
		 анонимно или нет, затем введи сообщение и отправь`,
		message.From.FirstName,
	)
	id := message.Chat.ID
	msg := tgbotapi.NewMessage(id, messageToSend)
	msg.ParseMode = "markdown"

	_, err := bot.Send(msg)
	if err != nil {
		log.LogRequestFile(fmt.Sprintf("[TELEGRAM]: %s", err))
		log.Errorf("[TELEGRAM]: %s", err)
		checkRetry(err)
	}

	user := &db.User{
		Anon:   true,
		Name:   message.From.FirstName,
		ChatID: strconv.FormatInt(id, 10),
	}
	err = db.UpdateEntry(user)
	if err != nil {
		log.LogRequestFile(fmt.Sprintf("[DB]: %s", err))
		log.Errorf("[DB]: %s", err)
	}

	log.Infof("User: %s, id: %d", message.From.FirstName, id)
}

// cmd /anonymous
func anonymousCmd(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	messageToSend := "*Как бы вы хотели отправить пост?*"
	markup := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("Анонимно", "Anon"),
			tgbotapi.NewInlineKeyboardButtonData("Неанонимно", "Nonanon"),
		},
	)
	id := message.Chat.ID
	msg := tgbotapi.NewMessage(id, messageToSend)
	msg.ParseMode = "markdown"
	msg.ReplyMarkup = markup

	_, err := bot.Send(msg)
	if err != nil {
		log.LogRequestFile(fmt.Sprintf("[TELEGRAM]: %s", err))
		log.Errorf("[TELEGRAM]: %s", err)
		checkRetry(err)
	}
}

//default
func defaultCmd(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {

	if message.Text != "" {
		id := *c.TGAdmin
		msg := tgbotapi.NewMessage(id, message.Text)
		msg.ParseMode = "markdown"

		_, err := bot.Send(msg)
		if err != nil {
			log.LogRequestFile(fmt.Sprintf("[TELEGRAM]: %s", err))
			log.Errorf("[TELEGRAM]: %s", err)
			checkRetry(err)
		}

		queryID := postQuery(bot, message.From.FirstName)

		user, err := db.GetOrCreateEntry(
			&db.User{
				ChatID: strconv.FormatInt(message.Chat.ID, 10),
				Name:   message.From.FirstName,
			})
		if err != nil {
			log.LogRequestFile(fmt.Sprintf("[DB]: %s", err))
			log.Errorf("[DB]: %s", err)
		}

		db.LockMessages.Lock()
		db.Messages[queryID] = &db.Post{
			User: *user, Message: message, Time: message.Time(),
		}
		db.LockMessages.Unlock()

		logData := fmt.Sprintf("User: %s suggest message, queryID: %d",
			user.Name, queryID)
		log.Infof(logData)
		log.LogRequestFile(logData)
	}

	//collect photos
	if message.Photo != nil {
		if message.Caption != "" {
			caption = message.Caption
		}

		lockPhotos.Lock()
		defer lockPhotos.Unlock()
		photos[message.Chat.ID] = append(photos[message.Chat.ID], (*message.Photo)[0].FileID)
	}

	//NOT SUPPORTED
	if message.Video != nil {
		log.Infof("Video message")
		msg := tgbotapi.NewMessage(message.Chat.ID, "Извините, видео пока не поддерживается")
		msg.ParseMode = "markdown"

		_, err := bot.Send(msg)
		if err != nil {
			log.LogRequestFile(fmt.Sprintf("[TELEGRAM]: %s", err))
			log.Errorf("[TELEGRAM]: %s", err)
			checkRetry(err)
		}
	}

	if message.VideoNote != nil {
		log.Infof("Videonote message")
		msg := tgbotapi.NewMessage(message.Chat.ID, "Извините, видео пока не поддерживается")
		msg.ParseMode = "markdown"

		_, err := bot.Send(msg)
		if err != nil {
			log.LogRequestFile(fmt.Sprintf("[TELEGRAM]: %s", err))
			log.Errorf("[TELEGRAM]: %s", err)
			checkRetry(err)
		}
	}

	if message.Animation != nil {
		log.Infof("Animation message")
		msg := tgbotapi.NewMessage(message.Chat.ID, "Извините, анимация пока не поддерживается")
		msg.ParseMode = "markdown"

		_, err := bot.Send(msg)
		if err != nil {
			log.LogRequestFile(fmt.Sprintf("[TELEGRAM]: %s", err))
			log.Errorf("[TELEGRAM]: %s", err)
			checkRetry(err)
		}
	}

	if message.Audio != nil {
		log.Infof("Audio message")
		msg := tgbotapi.NewMessage(message.Chat.ID, "Извините, аудио пока не поддерживается")
		msg.ParseMode = "markdown"

		_, err := bot.Send(msg)
		if err != nil {
			log.LogRequestFile(fmt.Sprintf("[TELEGRAM]: %s", err))
			log.Errorf("[TELEGRAM]: %s", err)
			checkRetry(err)
		}
	}

	if message.Document != nil {
		log.Infof("Document message")
		msg := tgbotapi.NewMessage(message.Chat.ID, "Извините, документы пока не поддерживается")
		msg.ParseMode = "markdown"

		_, err := bot.Send(msg)
		if err != nil {
			log.LogRequestFile(fmt.Sprintf("[TELEGRAM]: %s", err))
			log.Errorf("[TELEGRAM]: %s", err)
			checkRetry(err)
		}
	}

}
