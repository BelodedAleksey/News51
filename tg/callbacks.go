package tg

import (
	"fmt"
	"news/db"
	"news/log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//Handle callbacks
func callbackHandler(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	var message, logData string
	//get user
	user, err := db.GetOrCreateEntry(
		&db.User{
			ChatID: strconv.FormatInt(callback.Message.Chat.ID, 10),
			Name:   callback.Message.From.FirstName,
		})
	if err != nil {
		log.LogRequestFile(fmt.Sprintf("[DB]: %s", err))
		log.Errorf("[DB]: %s", err)
	}
	switch callback.Data {
	case "Anon":
		message = "*Спасибо*, в случае одобрения - ваше сообщение будет отправлено анонимно."
		user.Anon = true
		db.UpdateEntry(user)

		logData = fmt.Sprintf("User: %s became Anon", user.Name)
		break
	case "Nonanon":
		message = "*Спасибо*, в случае одобрения - мы укажем ваше авторство."
		user.Anon = false
		db.UpdateEntry(user)

		logData = fmt.Sprintf("User: %s became NonAnon", user.Name)
		break
	case "Accept":
		message = fmt.Sprintf("Ваш пост №%d был опубликован", callback.Message.MessageID)

		db.LockMessages.RLock()
		m := db.Messages[callback.Message.MessageID]
		db.LockMessages.RUnlock()

		user = &m.User
		delete(db.Messages, callback.Message.MessageID)
		SendPost(m, bot)

		logData = fmt.Sprintf("User: %s message, queryID: %d was published",
			m.User.Name, callback.Message.MessageID)

		//delete message with buttons
		delMsg := tgbotapi.NewDeleteMessage(*c.TGAdmin, callback.Message.MessageID)
		_, err := bot.DeleteMessage(delMsg)
		if err != nil {
			log.LogRequestFile(fmt.Sprintf("[TELEGRAM]: %s", err))
			log.Fatalf("[TELEGRAM]: %s", err)
		}
		break
	case "Decline":
		message = fmt.Sprintf("Ваш пост №%d был отклонен", callback.Message.MessageID)

		db.LockMessages.RLock()
		m := db.Messages[callback.Message.MessageID]
		db.LockMessages.RUnlock()

		user = &m.User
		delete(db.Messages, callback.Message.MessageID)

		logData = fmt.Sprintf("User: %s message id: %d was declined",
			user.Name, callback.Message.MessageID)

		//delete message with buttons
		delMsg := tgbotapi.NewDeleteMessage(*c.TGAdmin, callback.Message.MessageID)
		_, err := bot.DeleteMessage(delMsg)
		if err != nil {
			log.LogRequestFile(fmt.Sprintf("[TELEGRAM]: %s", err))
			log.Fatalf("[TELEGRAM]: %s", err)
		}
		break
	default:
		break
	}

	log.Infof(logData)
	log.LogRequestFile(logData)

	chatID, _ := strconv.ParseInt(user.ChatID, 10, 64)
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "markdown"

	_, err = bot.Send(msg)
	if err != nil {
		log.LogRequestFile(fmt.Sprintf("[TELEGRAM]: %s", err))
		log.Fatalf("[TELEGRAM]: %s", err)
	}

	_, err = bot.AnswerCallbackQuery(
		tgbotapi.NewCallback(
			callback.ID,
			message,
		),
	)
	if err != nil {
		log.LogRequestFile(fmt.Sprintf("[TELEGRAM]: %s", err))
		log.Fatalf("[TELEGRAM]: %s", err)
	}
}

//process messages
func messageHandler(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	switch message.Text {
	case "/start":
		startCmd(bot, message)
		break
	case "/anonymous":
		anonymousCmd(bot, message)
		break
	default:
		defaultCmd(bot, message)
		break
	}
}
