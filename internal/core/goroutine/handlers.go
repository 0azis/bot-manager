package goroutine

import (
	"botmanager/internal/core/domain"
	"botmanager/internal/core/utils"
	"context"
	"strconv"

	"github.com/go-telegram/bot"
	tg_domain "github.com/go-telegram/bot/models"
)

func (g goroutine) InitShopHandlers() {
	g.bot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, g.startHandler)
}

func (g goroutine) InitHomeHandlers() {
	g.bot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, g.setButton)
	g.bot.RegisterHandler(bot.HandlerTypeMessageText, "🔐 Получить код", bot.MatchTypeExact, g.sendCode)
}

func (g goroutine) startHandler(ctx context.Context, b *bot.Bot, update *tg_domain.Update) {
	botData, _ := b.GetMe(ctx)
	shopData, _ := g.store.Shop.GetByBotID(strconv.FormatInt(botData.ID, 10))

	b.SetChatMenuButton(ctx, &bot.SetChatMenuButtonParams{
		ChatID: update.Message.From.ID,
		MenuButton: tg_domain.MenuButtonWebApp{
			Type: "web_app",
			Text: shopData.TitleButton,
			WebApp: tg_domain.WebAppInfo{
				URL: domain.WebLink(shopData.ID),
			},
		},
	})

	var url string
	photos, _ := b.GetUserProfilePhotos(ctx, &bot.GetUserProfilePhotosParams{
		UserID: update.Message.From.ID,
	})

	if photos.TotalCount == 0 {
		url = ""
	} else {
		file, _ := b.GetFile(ctx, &bot.GetFileParams{
			FileID: photos.Photos[0][0].FileID,
		})
		url = b.FileDownloadLink(file)
	}

	telegramID := strconv.FormatInt(update.Message.From.ID, 10)
	if !g.store.Subscriber.IsSubscribed(telegramID, shopData.ID) {
		s := domain.Subscriber{
			UserName:   update.Message.From.Username,
			FirstName:  update.Message.From.FirstName,
			LastName:   update.Message.From.LastName,
			AvatarUrl:  url,
			TelegramID: telegramID,
			ShopID:     shopData.ID,
		}
		g.store.Subscriber.Insert(s)
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   shopData.FirstLaunch,
	})
}

func (g goroutine) listenMessages(ctx context.Context, b *bot.Bot, update *tg_domain.Update) {
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
}

func (g goroutine) sendCode(ctx context.Context, b *bot.Bot, update *tg_domain.Update) {
	code := utils.GenerateCode()
	err := g.redisDB.SetCode(strconv.FormatInt(update.Message.From.ID, 10), code)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.From.ID,
			Text:   "Неизвестная ошибка при отправке кода. Пожалуйста попробуйте позже",
		})
		return
	}

	user := update.Message.From

	var url string
	photos, _ := b.GetUserProfilePhotos(ctx, &bot.GetUserProfilePhotosParams{
		UserID: user.ID,
	})

	if photos.TotalCount == 0 {
		url = ""
	} else {
		file, _ := b.GetFile(ctx, &bot.GetFileParams{
			FileID: photos.Photos[0][0].FileID,
		})
		url = b.FileDownloadLink(file)
	}

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
			Text:   err.Error(),
		})
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		Text:      "<b>Код активен 30 секунд:</b> " + code,
		ParseMode: "HTML",
	})
}

func (g goroutine) setButton(ctx context.Context, b *bot.Bot, update *tg_domain.Update) {
	kb := &tg_domain.ReplyKeyboardMarkup{
		Keyboard: [][]tg_domain.KeyboardButton{
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
