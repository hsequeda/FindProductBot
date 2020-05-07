package main

import (
	httpClient "findTuEnvioBot/client"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"net/http"
)

type MyBot struct {
	botName string
	bot     *tgbotapi.BotAPI
	client  *httpClient.Client
}

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	conf, err := initConfig()
	if err != nil {
		logrus.Fatalln(err)
		return
	}
	logrus.Println("Starting bot..")
	mBot := initBot(conf)

	updateCh, errCh := mBot.getUpdateCh(conf)
	if errCh != nil {
		logrus.Fatal(errCh)
	}

	for update := range updateCh {
		switch {
		case update.CallbackQuery != nil:
			mBot.handleCallBackQuery(update.CallbackQuery)

		case update.Message != nil:
			if update.Message.Chat.Type == "private" {
				mBot.handlePrivateMessage(update.Message)
			} else {
				mBot.handlePublicMessage(update.Message)
			}

			break
		case update.InlineQuery != nil:
			mBot.handleInlineQuery(update.InlineQuery)
		}
	}
}

func initBot(conf config) MyBot {
	bot, err := tgbotapi.NewBotAPI(conf.BotToken)
	if err != nil {
		logrus.Fatal(err)
	}

	u, err := bot.GetMe()
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Print(u)

	bot.Debug = true
	mBot := MyBot{bot: bot, client: httpClient.NewClient(), botName: conf.BotName}
	if conf.WebHook.Domain == "" {
		return mBot
	}

	if conf.WebHook.WithCert {
		_, err = bot.SetWebhook(tgbotapi.NewWebhookWithCert(conf.WebHook.Domain, conf.WebHook.CertPath))
		if err != nil {
			logrus.Fatal(err)
		}

		go raiseServer(conf)
	} else {
		_, err = bot.SetWebhook(tgbotapi.NewWebhook(conf.WebHook.Domain))
		if err != nil {
			logrus.Fatal(err)
		}

		go raiseServer(conf)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		logrus.Fatal(err)
	}

	if info.LastErrorDate != 0 {
		logrus.Printf("[Telegram callback failed]%s", info.LastErrorMessage)
	}

	return mBot
}

func raiseServer(conf config) {
	logrus.Println("Server Started")
	if conf.WebHook.WithCert {
		err := http.ListenAndServeTLS("0.0.0.0:"+conf.WebHook.Port, conf.WebHook.CertPath, conf.WebHook.KeyPath, nil)
		if err != nil {
			logrus.Println(err)
		}
	} else {
		err := http.ListenAndServe("0.0.0.0:"+conf.WebHook.Port, nil)
		if err != nil {
			logrus.Println(err)
		}
	}
}

func (m *MyBot) getUpdateCh(conf config) (tgbotapi.UpdatesChannel, error) {
	if conf.WebHook.Domain != "" {
		return m.bot.ListenForWebhook(conf.WebHook.Path), nil
	}
	return m.bot.GetUpdatesChan(tgbotapi.UpdateConfig{
		Offset:  0,
		Timeout: 0,
	})
}
