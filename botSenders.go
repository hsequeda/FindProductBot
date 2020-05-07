package main

import (
	"findTuEnvioBot/products"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (m MyBot) sendInlineKeyboardSelectProvince(chatId int64) {
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
	_, err := m.bot.Send(msg)

	if err != nil {
		logrus.Warn(err)
	}
}

func (m MyBot) sendUserPanel(chatId int64, text string) {
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
	msg.ParseMode = "html"
	msg.ReplyMarkup = replyKeyboard
	_, err := m.bot.Send(msg)
	if err != nil {
		logrus.Warn(err)
	}
}

func (m MyBot) sendQueryResultList(list []products.Product, inlineQueryID string) {
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

	_, err := m.bot.AnswerInlineQuery(tgbotapi.InlineConfig{
		InlineQueryID: inlineQueryID,
		Results:       resultList,
		CacheTime:     10000,
	})
	if err != nil {
		logrus.Warn(err)
	}
}

func (m MyBot) sendResultMessage(chatId int64, productList []products.Product) {
	for _, prod := range productList {
		rawMsg := fmt.Sprintf(
			`
		<b>Producto: %s</b>,
		<b>Precio: %s</b>,
		<b>Tienda: %s</b>,
		<a href="%s">Ver Producto</a>,
`, prod.GetName(), prod.GetPrice(), prod.GetStore(), prod.GetLink())

		msg := tgbotapi.NewMessage(chatId, rawMsg)
		msg.ParseMode = "html"
		_, err := m.bot.Send(msg)
		if err != nil {
			logrus.Warn(err)
		}
	}
}

func (m MyBot) sendProvinceNotSelectError(chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "❌ Necesita seleccionar una provincia.")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Seleccionar provincia!", fmt.Sprintf("https://t.me/%s?start=start", m.botName)),
		),
	)
	_, err := m.bot.Send(msg)
	if err != nil {
		logrus.Warn(err)
	}
}

func (m MyBot) sendInsertCommandValidError(chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "❌ Inserte un comando valido.")
	_, err := m.bot.Send(msg)
	if err != nil {
		logrus.Warn(err)
	}
}

func (m MyBot) sendInstructions(chatId int64) {
	msg := tgbotapi.NewMessage(chatId,
		`Este bot te ayuda a encontrar 🔍 productos en las tiendas virtuales.
			 Hasta ahora abarca todas las tiendas de <a href="https://www.tuenvio.cu">Tu Envio</a> y la tienda de
			<a href="https://5tay42.xetid.cu">5ta y 42</a> con la meta de añadir las restantes.
			> <b>Primero que todo:</b> 
			- Para empezar a usar este bot tiene que iniciar un chat privado y seleccionar una Provincia. Usted
			puede cambiar la provincia cuando desee. La provincia que seleccione es en la cual se realizaran las busquedas.

			<b>Modos de uso:</b>
			1- Privado(El chat privado con el bot): En este modo usted va a tener acceso al comando /buscar 'Producto' 
			para buscar un producto, tambien puede hacerlo escribiendo en el chat privado lo que quiere buscar.
			Ademas va a tener a disposicion un listado de botones que le haran la vida mas facil☺️ brindandole las opciones de:
				1. Cambiar o añadir su provincia. 
				2. Mostrar esta ayuda para si se olvida de algo. 
				3. Ver su perfil donde vera su usuario con la provincia que tiene vinculada.
				4. Subscribirse a un patron de busqueda para notificarle cuando encontremos algo.(Pendiente)
			2- Publico( Añadiendo el bot a un grupo para uso publico): En este modo usted va a tener acceso al comando
			'/buscar', si no ha iniciado una conversacion con el bot y le hace un pedido se le mostrara un boton de enlace
			 para realizar esta tarea.
			3- Inline: A este modo se accede escriviendo '@buscarTuEnvioBot 'patron a buscar'. Puede acceder a este modo desde
             cualquier parte de telegram, lo mismo un grupo como un chat privado, cuando escriba un patron tiene que esperar
             unos segundos para que se realice el pedido asi que uselo con calma.
			`,
	)
	msg.ParseMode = "html"
	_, err := m.bot.Send(msg)
	if err != nil {
		logrus.Warn(err)
	}
}
