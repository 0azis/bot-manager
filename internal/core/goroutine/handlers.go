package goroutine

import (
	"botmanager/internal/core/domain"
	"botmanager/internal/core/utils"
	"context"
	"fmt"
	"strconv"

	"github.com/0azis/bot"
	"github.com/0azis/bot/models"
)

func (g goroutine) InitShopHandlers() {
	g.bot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, g.startHandler)
	g.bot.SetDefaultHandler(g.listenMessages)
}

func (g goroutine) InitHomeHandlers() {
	g.bot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, g.setButton)
	g.bot.RegisterHandler(bot.HandlerTypeMessageText, "🔐 Получить код", bot.MatchTypeExact, g.sendCode)
}

func (g goroutine) startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	user := update.Message.From
	botData, _ := b.GetMe(ctx)
	shopData, _ := g.store.Shop.GetByBotID(strconv.FormatInt(botData.ID, 10))

	b.SetChatMenuButton(ctx, &bot.SetChatMenuButtonParams{
		ChatID: user.ID,
		MenuButton: models.MenuButtonWebApp{
			Type: "web_app",
			Text: shopData.TitleButton,
			WebApp: models.WebAppInfo{
				URL: domain.WebLink(shopData.ID),
			},
		},
	})

	url := utils.GetTelegramPhoto(b, ctx, user.ID)

	telegramID := strconv.FormatInt(user.ID, 10)
	if !g.store.Subscriber.IsSubscribed(telegramID, shopData.ID) {
		s := domain.Subscriber{
			UserName:   user.Username,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			AvatarUrl:  url,
			TelegramID: telegramID,
			ShopID:     shopData.ID,
		}
		g.store.Subscriber.Insert(s)
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: user.ID,
		Text:   shopData.FirstLaunch,
	})
}

func (g goroutine) listenMessages(ctx context.Context, b *bot.Bot, update *models.Update) {
	botData, _ := b.GetMe(ctx)

	shop, _ := g.store.Shop.GetByBotID(strconv.FormatInt(botData.ID, 10))

	subscriber, _ := g.store.Subscriber.GetByTelegramID(strconv.FormatInt(update.Message.From.ID, 10))

	subscriberMsg := domain.Message{
		Text:         update.Message.Text,
		SubscriberID: subscriber.ID,
		IsFromUser:   true,
		BotID:        shop.ID,
	}

	err := g.store.Message.Insert(subscriberMsg)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.From.ID,
			Text:   "Произошла ошибка при отправке сообщения менеджеру. Попробуйте позже.",
		})
	}

	owner, _ := g.store.User.Get(shop.UserID)

	newMsgNotification := domain.Notification{
		UserID:   owner.TelegramID,
		Text:     "Вам пришло новое сообщение",
		Link:     "https://tgrocket.ru/app/messages",
		LinkText: "Читать сообщения",
	}
	homeBot := g.pool.GetHomeBot()
	homeBot.SendNotification(newMsgNotification)
}

func (g goroutine) sendCode(ctx context.Context, b *bot.Bot, update *models.Update) {
	user := update.Message.From

	code := utils.GenerateCode()
	err := g.redisDB.SetCode(strconv.FormatInt(update.Message.From.ID, 10), code)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: user.ID,
			Text:   "Неизвестная ошибка при отправке кода. Пожалуйста попробуйте позже",
		})
		return
	}

	url := utils.GetTelegramPhoto(b, ctx, user.ID)

	u := domain.User{
		TelegramID:   strconv.FormatInt(user.ID, 10),
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		UserName:     user.Username,
		IsPremium:    user.IsPremium,
		LanguageCode: user.LanguageCode,
		AvatarUrl:    url,
	}
	err = g.store.User.Insert(u)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.From.ID,
			Text:   "Неизвестная ошибка при отправке кода. Пожалуйста попробуйте позже",
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		Text:      fmt.Sprintf("*Код активен 30 секунд:* `%s`", code),
		ParseMode: "Markdown",
	})
}

func (g goroutine) setButton(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := &models.ReplyKeyboardMarkup{
		Keyboard: [][]models.KeyboardButton{
			{
				{Text: "🔐 Получить код"},
			},
		},
		ResizeKeyboard: true,
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.From.ID,
		Text:        "Запросите код",
		ReplyMarkup: kb,
	})

}
