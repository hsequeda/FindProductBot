package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"net/http"
)

var bot *tgbotapi.BotAPI

func init() {
	var err error
	initConfig()

	bot, err = tgbotapi.NewBotAPI(conf.botToken)
	if err != nil {
		logrus.Fatal(err)
	}

	u, err := bot.GetMe()
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Print(u)

	bot.Debug = true

	if conf.webHook.domain != "" {
		if conf.webHook.withCert {
			_, err = bot.SetWebhook(tgbotapi.NewWebhookWithCert(conf.webHook.domain, conf.webHook.certPath))
			if err != nil {
				logrus.Fatal(err)
			}

			go raiseServer(true)
		} else {
			_, err = bot.SetWebhook(tgbotapi.NewWebhook(conf.webHook.domain))
			if err != nil {
				logrus.Fatal(err)
			}

			go raiseServer(true)
		}

		info, err := bot.GetWebhookInfo()
		if err != nil {
			logrus.Fatal(err)
		}

		if info.LastErrorDate != 0 {
			logrus.Printf("[Telegram callback failed]%s", info.LastErrorMessage)
		}
	}
}

func main() {
	logrus.Println("Starting bot..")
	updateCh, err := getUpdateCh()
	if err != nil {
		logrus.Fatal(err)
	}

	for update := range updateCh {
		logrus.Println(update)
		switch {
		case update.CallbackQuery != nil:
			handleCallBackQuery(update.CallbackQuery)

		case update.Message != nil:
			if update.Message.Chat.Type == "private" {
				handlePrivateMessage(update.Message)
			} else {
				handlePublicMessage(update.Message)
			}

			break
		case update.InlineQuery != nil:
			handleInlineQuery(update.InlineQuery)
		}
	}
}

func raiseServer(withCert bool) {
	logrus.Println("Server Started")
	if withCert {
		err := http.ListenAndServeTLS("0.0.0.0:"+conf.webHook.port, conf.webHook.certPath, conf.webHook.keyPath, nil)
		if err != nil {
			logrus.Println(err)
		}
	} else {
		err := http.ListenAndServe("0.0.0.0:"+conf.webHook.port, nil)
		if err != nil {
			logrus.Println(err)
		}
	}
}

func getUpdateCh() (tgbotapi.UpdatesChannel, error) {
	if conf.webHook.domain != "" {
		return bot.ListenForWebhook(conf.webHook.path), nil
	} else {
		return bot.GetUpdatesChan(tgbotapi.UpdateConfig{
			Offset:  0,
			Timeout: 0,
		})
	}
}
