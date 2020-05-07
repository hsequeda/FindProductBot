package main

import (
	"findTuEnvioBot/products"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func sendInlineKeyboardSelectProvince(chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "Seleccione una provincia:")

	var markup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("La Habana", "La Habana"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Pinar del Rio", "Pinar del Rio"),
			tgbotapi.NewInlineKeyboardButtonData("Artemisa", "Artemisa"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Mayabeque", "Mayabeque"),
			tgbotapi.NewInlineKeyboardButtonData("Matanzas", "Matanzas"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Cienfuegos", "Cienfuegos"),
			tgbotapi.NewInlineKeyboardButtonData("Villa Clara", "Villa Clara"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Sancti Spiritus", "Sancti Spiritus"),
			tgbotapi.NewInlineKeyboardButtonData("Ciego de Avila", "Ciego de Avila"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Camaguey", "Camaguey"),
			tgbotapi.NewInlineKeyboardButtonData("Las Tunas", "Las Tunas"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Granma", "Granma"),
			tgbotapi.NewInlineKeyboardButtonData("Santiago de Cuba", "Santiago de Cuba"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Guantanamo", "Guantanamo"),
			tgbotapi.NewInlineKeyboardButtonData("La Isla", "La Isla"),
		),
	)
	msg.ReplyMarkup = markup
	_, err := bot.Send(msg)

	if err != nil {
		logrus.Warn(err)
	}
}

func sendUserPanel(chatId int64, text string) {
	msg := tgbotapi.NewMessage(chatId, text)
	replyKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("👤 Mi Perfil"),
			tgbotapi.NewKeyboardButton("🆘 Help"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🗺 Seleccionar Provincia"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📋 Adicionar subscripcion"),
		),
	)
	msg.ReplyMarkup = replyKeyboard
	_, err := bot.Send(msg)
	if err != nil {
		logrus.Warn(err)
	}
}

func sendQueryResultList(list []products.Product, inlineQueryID string) {
	var resultList = make([]interface{}, 0)

	for _, prod := range list {
		msg := fmt.Sprintf(
			`
		<b>Producto: %s</b>,
		<b>Precio: %s</b>,
		<b>Tienda: %s</b>,
		<a href="%s">Ver Producto</a>,
`, prod.GetName(), prod.GetPrice(), prod.GetStore(), prod.GetLink())

		inlineQueryResult := tgbotapi.NewInlineQueryResultArticleHTML(uuid.New().String(),
			fmt.Sprintf("%s - %s", prod.GetName(), prod.GetStore()), msg)

		resultList = append(resultList, inlineQueryResult)
	}

	_, err := bot.AnswerInlineQuery(tgbotapi.InlineConfig{
		InlineQueryID: inlineQueryID,
		Results:       resultList,
		CacheTime:     10000,
	})
	if err != nil {
		logrus.Warn(err)
	}
}
