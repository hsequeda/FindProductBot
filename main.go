package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

type MyBot struct {
	*tgbotapi.BotAPI
}

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	conf, err := initConfig()
	if err != nil {
		logrus.Fatalln(err)
		return
	}
	logrus.Println("Starting bot..")
	bot := initBot(conf)

	updateCh, errCh := bot.getUpdateCh(conf)
	if errCh != nil {
		logrus.Fatal(errCh)
	}

	for update := range updateCh {
		logrus.Print("a", update)
		switch {
		case update.Message != nil:
			if update.Message.Chat.Type == "private" {
				if prov, ok := isProvince(update.Message.Text); ok {
					err := InsertUser(strconv.Itoa(update.Message.From.ID), prov)
					if err != nil {
						logrus.Warn(err)
						continue
					}
					logrus.Println(bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(
						"Bien hecho!!\n Haz "+"seleccionado como provincia '%s'. Ahora las busquedas que me"+
							" pidas las realizare en las tiendas que estan en esa provincia ☺", update.Message.Text))))
					logrus.Println(bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Puedes cambiar de provincia"+
						" cada vez que quieras como se muestra abajo ⬇")))
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Seleccione su provincia:")
				msg.ReplyMarkup = getProvKeyboard()
				logrus.Println(bot.Send(msg))
			} else {
				logrus.Println(bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "No me gustan los grupos! 😠")))
			}

			break
		case update.InlineQuery != nil:
			if len(update.InlineQuery.Query) >= 2 {
				user, err := GetUser(strconv.Itoa(update.InlineQuery.From.ID))
				switch {
				case err == errValEmpty || err == errBucketEmpty:
					_, err := bot.AnswerInlineQuery(tgbotapi.InlineConfig{
						InlineQueryID: update.InlineQuery.ID,
						Results: []interface{}{
							tgbotapi.NewInlineQueryResultArticleHTML(uuid.New().String(), "Necesitas empezar una"+
								" conversacion conmigo primero.", "Necesitas empezar una conversacion "+
								"conmigo para poder usarme <a href=\"https://t.me/findTuEnvioBot\">Empezar conversacion.</a>"),
						},
					})
					if err != nil {
						logrus.Print(err)
						continue
					}
					break
				case err != nil:
					logrus.Print(err)
					continue
				default:
					var productList = make([]Product, 0)
					prov := provinces[user.Province]
					for _, store := range prov.stores {
						products, err := GetProductsByPattern(store.rawName, update.InlineQuery.Query)
						if err != nil {
							logrus.Print(err)
							continue
						}

						for i := range products.Content {
							products.Content[i].Store = store.name
							productList = append(productList, products.Content[i])
						}
					}

					result, err := getQueryResultList(productList)
					if err != nil {
						logrus.Print(err)
						continue
					}

					_, err = bot.AnswerInlineQuery(tgbotapi.InlineConfig{
						InlineQueryID: update.InlineQuery.ID,
						Results:       result,
					})
					if err != nil {
						logrus.Print(err)
						continue
					}

				}
			}

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

	if conf.WebHook.Domain == "" {
		return MyBot{bot}
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
	return MyBot{bot}
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

func (b *MyBot) getUpdateCh(conf config) (tgbotapi.UpdatesChannel, error) {
	if conf.WebHook.Domain != "" {
		return b.ListenForWebhook(conf.WebHook.Path), nil
	}
	return b.GetUpdatesChan(tgbotapi.UpdateConfig{
		Offset:  0,
		Timeout: 0,
	})
}

func getQueryResultList(list []Product) ([]interface{}, error) {
	var resultList = make([]interface{}, 0)

	for _, prod := range list {
		msg := fmt.Sprintf(
			`
		<b>Producto: %s</b>,
		<b>Precio: %s</b>,
		<b>Tienda: %s</b>,
		<a href="%s">Enlace</a>,
		`, strings.TrimSpace(prod.Name), strings.TrimSpace(prod.Price), strings.TrimSpace(prod.Store), strings.TrimSpace(prod.Link))

		inlineQueryResult := tgbotapi.NewInlineQueryResultArticleHTML(uuid.New().String(), strings.TrimSpace(prod.Name), msg)
		resultList = append(resultList, inlineQueryResult)
	}
	return resultList, nil
}

func getProvKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboardList := make([][]tgbotapi.KeyboardButton, 0)
	for _, prov := range provinces {
		keyboardList = append(keyboardList, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(prov.name)))
	}

	return tgbotapi.NewReplyKeyboard(keyboardList...)

}

func isProvince(text string) (string, bool) {
	for key := range provinces {
		if provinces[key].name == text {
			return key, true
		}
	}
	return "", false
}
