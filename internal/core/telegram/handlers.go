package telegram

import (
	"botmanager/internal/core/domain"
	"context"
	"strconv"

	"github.com/go-telegram/bot"
	tg_domain "github.com/go-telegram/bot/models"
)

func (g goroutine) initHandlers() {
	g.bot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, g.startHandler)
}

func (g goroutine) startHandler(ctx context.Context, b *bot.Bot, update *tg_domain.Update) {
	b.SetChatMenuButton(ctx, &bot.SetChatMenuButtonParams{
		ChatID: update.Message.From.ID,
		MenuButton: tg_domain.MenuButtonWebApp{
			Type: "web_app",
			Text: g.botData.TitleButton,
			WebApp: tg_domain.WebAppInfo{
				URL: domain.WebLink(g.botData.ID),
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
	if !g.store.Subscriber.IsSubscribed(telegramID, g.botData.ID) {
		s := domain.Subscriber{
			UserName:   update.Message.From.Username,
			FirstName:  update.Message.From.FirstName,
			LastName:   update.Message.From.LastName,
			AvatarUrl:  url,
			TelegramID: telegramID,
			ShopID:     g.botData.ID,
		}
		g.store.Subscriber.Insert(s)
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   g.botData.FirstLaunch,
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
