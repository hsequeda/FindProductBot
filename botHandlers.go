package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	storeClient "github.com/stdevHsequeda/CubanProductFinder"
	"strconv"
	"strings"
)

func (m MyBot) handleCallBackQuery(query *tgbotapi.CallbackQuery) {
	_, err := m.bot.AnswerCallbackQuery(tgbotapi.NewCallback(query.ID, query.Data))
	if err != nil {
		logrus.Warn(err)
	}
	var storeList []string
	for _, store := range provinces[query.Data] {
		storeList = append(storeList, store.name)
	}
	msg := fmt.Sprintf("<b>Ha seleccionado %s:</b>\n <b>Tiendas:</b>\n %s", query.Data, strings.Join(storeList, "\n"))
	m.sendUserPanel(query.Message.Chat.ID, msg)

	err = InsertUser(strconv.FormatInt(query.Message.Chat.ID, 10), query.Data)
	if err != nil {
		logrus.Warn(err)
	}

	_, err = m.bot.DeleteMessage(tgbotapi.DeleteMessageConfig{
		ChatID:    query.Message.Chat.ID,
		MessageID: query.Message.MessageID,
	})
	if err != nil {
		logrus.Warn(err)
	}
}

func (m MyBot) handlePublicMessage(message *tgbotapi.Message) {
	switch {
	case message.Text == "/help":
		// Send instructions
		m.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "For implement"))

	case strings.Split(message.Text, " ")[0] == "/buscar":
		// Search product
		user, err := GetUser(strconv.Itoa(message.From.ID))
		switch {
		case err == errValEmpty || err == errBucketEmpty:
			m.sendProvinceNotSelectError(message.Chat.ID)
			break
		case err != nil:
			logrus.Warn(err)
			return
		default:
			var pattern string
			splitText := strings.Split(message.Text, " ")
			if len(splitText) >= 1 {
				pattern = strings.Join(splitText[1:], " ")
			}

			prods, err := m.sc.SearchProduct(pattern)
			if err != nil {
				logrus.Print(err)
			}

			var results = make([]storeClient.Product, 0)
			for _, store := range provinces[user.Province] {
				for i := range prods {
					if prods[i].GetSection().GetStore().Name == store.name {
						results = append(results, prods[i])
					}
				}
			}

			m.sendResultMessage(message.Chat.ID, prods)
		}

	case strings.Split(message.Text, " ")[0] == "/subscribirme":
		// Buscar
	default:
		// inserte un comando valido
	}
}

func (m MyBot) handlePrivateMessage(privateMsg *tgbotapi.Message) {
	if privateMsg.IsCommand() {
		switch {
		case privateMsg.Text == "/help":
			m.sendInstructions(privateMsg.Chat.ID)
			// Send instructions
		case strings.Split(privateMsg.Text, " ")[0] == "/start":
			if len(strings.Split(privateMsg.Text, " ")) >= 2 {
				if strings.Split(privateMsg.Text, " ")[1] == "start" {
					m.sendInlineKeyboardSelectProvince(privateMsg.Chat.ID)
				}
			}

			// Send instructions
			m.sendUserPanel(privateMsg.Chat.ID, "Seleccione la opcion que desee realizar:")
		case strings.Split(privateMsg.Text, " ")[0] == "/buscar":
			// Search Product
			user, err := GetUser(strconv.Itoa(privateMsg.From.ID))
			switch {
			case err == errValEmpty || err == errBucketEmpty:
				m.sendProvinceNotSelectError(privateMsg.Chat.ID)
				break
			case err != nil:
				logrus.Warn(err)
				return
			default:
				var pattern string
				splitText := strings.Split(privateMsg.Text, " ")
				if len(splitText) >= 1 {
					pattern = strings.Join(splitText[1:], " ")
				}

				prods, err := m.sc.SearchProduct(pattern)
				if err != nil {
					logrus.Print(err)
				}

				var results = make([]storeClient.Product, 0)
				for _, store := range provinces[user.Province] {
					for i := range prods {
						if prods[i].GetSection().GetStore().Name == store.name {
							results = append(results, prods[i])
						}
					}
				}

				m.sendResultMessage(privateMsg.Chat.ID, prods)
			}

		case strings.Split(privateMsg.Text, " ")[0] == "/subscribirme":
			// Subscribe
		default:
			m.sendInsertCommandValidError(privateMsg.Chat.ID)
			m.sendInstructions(privateMsg.Chat.ID)
		}
		return
	}
	switch privateMsg.Text {
	case "🆘 Help":
		// Send instuctions
		m.sendInstructions(privateMsg.Chat.ID)
	case "🗺 Seleccionar Provincia":
		// Send province list
		m.sendInlineKeyboardSelectProvince(privateMsg.Chat.ID)
		break
	case "📋 Adicionar subscripcion":
		m.bot.Send(tgbotapi.NewMessage(privateMsg.Chat.ID, "For implement"))
	// Add subscription
	case "👤 Mi Perfil":
		user, err := GetUser(strconv.FormatInt(privateMsg.Chat.ID, 10))
		switch {
		case err == errValEmpty || err == errBucketEmpty:
			m.sendInlineKeyboardSelectProvince(privateMsg.Chat.ID)
			break
		case err != nil:
			logrus.Warn(err)
			break
		default:
			msg := tgbotapi.NewMessage(privateMsg.Chat.ID, fmt.Sprintf(
				"👤 <b>Usuario:</b> %s,\n 🗺 <b>Provincia:</b> %s", privateMsg.From.FirstName, user.Province))
			msg.ParseMode = "html"
			_, err := m.bot.Send(msg)
			if err != nil {
				logrus.Warn(err)
				break
			}
		}
	default:
		// Search Product
		user, err := GetUser(strconv.Itoa(privateMsg.From.ID))
		switch {
		case err == errValEmpty || err == errBucketEmpty:
			m.sendProvinceNotSelectError(privateMsg.Chat.ID)
			break
		case err != nil:
			logrus.Warn(err)
			return
		default:
			prods, err := m.sc.SearchProduct(privateMsg.Text)
			if err != nil {
				logrus.Print(err)
			}

			var results = make([]storeClient.Product, 0)
			for _, store := range provinces[user.Province] {
				for i := range prods {
					if prods[i].GetSection().GetStore().Name == store.name {
						results = append(results, prods[i])
					}
				}
			}

			m.sendResultMessage(privateMsg.Chat.ID, prods)
		}
	}
}

func (m MyBot) handleInlineQuery(query *tgbotapi.InlineQuery) {
	if len(query.Query) >= 2 {
		user, err := GetUser(strconv.Itoa(query.From.ID))
		switch {
		case err == errValEmpty || err == errBucketEmpty:
			_, err := m.bot.AnswerInlineQuery(tgbotapi.InlineConfig{
				InlineQueryID: query.ID,
				Results: []interface{}{
					tgbotapi.NewInlineQueryResultArticleHTML(uuid.New().String(),
						"Necesitas empezar una conversacion conmigo primero.",
						fmt.Sprintf("Necesitas empezar una conversacion "+
							"conmigo para que me digas tu provincia"+
							" <a href=\"https://t.me/%s?start=start\">Empezar conversacion.</a>", m.botName)),
				},
			})

			if err != nil {
				logrus.Warn(err)
			}
			break
		case err != nil:
			logrus.Warn(err)
			return
		default:
			prods, err := m.sc.SearchProduct(query.Query)
			if err != nil {
				logrus.Print(err)
			}

			var results = make([]storeClient.Product, 0)
			for _, store := range provinces[user.Province] {
				for i := range prods {
					if prods[i].GetSection().GetStore().Name == store.name {
						results = append(results, prods[i])
					}
				}
			}

			if len(prods) >= 1 {
				m.sendQueryResultList(prods, query.ID)
			}
		}
	}
}
