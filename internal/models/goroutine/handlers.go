package goroutine

import (
	"botmanager/internal/models"
	"context"
	"strconv"

	"github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
)

type BotHandlers interface {
	initHandlers()
	startHandler(ctx context.Context, b *bot.Bot, update *tg_models.Update)
}

func (g Goroutine) initHandlers() {
	g.bot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, g.startHandler)	
}

func (g Goroutine) startHandler(ctx context.Context, b *bot.Bot, update *tg_models.Update) {
	b.SetChatMenuButton(ctx, &bot.SetChatMenuButtonParams{
		ChatID: update.Message.From.ID,
		MenuButton: tg_models.MenuButtonWebApp{
			Type: "web_app",
			Text: g.botData.TitleButton,
			WebApp: tg_models.WebAppInfo{
				URL: models.WebLink(g.botData.ID),
			},
		},
	})

	photos, _ := b.GetUserProfilePhotos(ctx, &bot.GetUserProfilePhotosParams{
		UserID: update.Message.From.ID,
	})
	file, _ := b.GetFile(ctx, &bot.GetFileParams{
		FileID: photos.Photos[0][0].FileID,
	})

	url := b.FileDownloadLink(file)

	telegramID := strconv.FormatInt(update.Message.From.ID, 10)
	if !g.store.Subscriber().IsSubscribed(telegramID, g.botData.ID) {
		s := models.Subscriber{
			UserName:   update.Message.From.Username,
			FirstName:  update.Message.From.FirstName,
			LastName:   update.Message.From.LastName,
			AvatarUrl:  url,
			TelegramID: telegramID,
			ShopID:     g.botData.ID,
		}
		g.store.Subscriber().Insert(s)
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   g.botData.FirstLaunch,
	})
}
