package tools

import (
	"botmanager/internal/models"
	"botmanager/internal/repos"
	"context"
	"os"
	"os/signal"
	"strconv"

	"github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
	"github.com/jmoiron/sqlx"
)

func BotWorker(token string, shopRepo repos.ShopRepo, subRepo repos.SubscriberRepo, pool *models.GoroutinesPool) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	botData, _ := shopRepo.Get(token)

	opts := []bot.Option{
		bot.WithDefaultHandler(func(ctx context.Context, b *bot.Bot, update *tg_models.Update) {
			b.SetChatMenuButton(ctx, &bot.SetChatMenuButtonParams{
				ChatID: update.Message.Chat.ID,
				MenuButton: tg_models.MenuButtonWebApp{
					Type: "web_app",
					Text: botData.TitleButton,
					WebApp: tg_models.WebAppInfo{
						URL: models.WebLink(botData.ID),
					},
				},
			})
			b.SetMyDescription(ctx, &bot.SetMyDescriptionParams{
				Description: botData.Description,
			})
		}),
	}

	b, _ := bot.New(token, opts...)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, func(ctx context.Context, b *bot.Bot, update *tg_models.Update) {
		tgID := strconv.FormatInt(update.Message.From.ID, 10)
		botInfo, _ := b.GetMe(ctx)
		shop, _ := shopRepo.GetBy("bot_id", strconv.FormatInt(botInfo.ID, 10))

		photos, _ := b.GetUserProfilePhotos(ctx, &bot.GetUserProfilePhotosParams{
			UserID: update.Message.From.ID,
		})
		file, _ := b.GetFile(ctx, &bot.GetFileParams{
			FileID: photos.Photos[0][0].FileID,
		})

		url := b.FileDownloadLink(file)

		if !subRepo.IsSubscribed(tgID, shop.ID) {
			s := models.Subscriber{
				UserName:   update.Message.From.Username,
				FirstName:  update.Message.From.FirstName,
				LastName:   update.Message.From.LastName,
				AvatarUrl:  url,
				TelegramID: tgID,
				ShopID:     shop.ID,
			}
			subRepo.Insert(s)
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   botData.FirstLaunch,
		})
	})

	// get channel that created
	channel := pool.Get(token)

	for {
		select {
		case msg := <-channel:

			switch msg.MsgType {
			case "status":

				switch msg.Value {
				case false:
					close(channel)
					return
				case true:
					go b.Start(ctx)
				}

			case "mail":
				mail := msg.Value.(models.Mail)
				subs, _ := subRepo.Select(mail.ShopID)

				for i := range subs {
					kb := &tg_models.InlineKeyboardMarkup{
						InlineKeyboard: [][]tg_models.InlineKeyboardButton{},
					}

					if mail.AddButton {
						kb.InlineKeyboard = append(kb.InlineKeyboard, []tg_models.InlineKeyboardButton{{Text: "Открыть меню", WebApp: &tg_models.WebAppInfo{URL: models.WebLink(botData.ID)}}})
					}

					if mail.PhotoLink != "" {
						b.SendPhoto(ctx, &bot.SendPhotoParams{
							ChatID:  subs[i].TelegramID,
							Photo:   &tg_models.InputFileString{Data: "https://tgrocket.ru/api/uploads/share/" + mail.PhotoLink}, 
							Caption: mail.Text,
							ReplyMarkup: kb,
						})
					} else {
						b.SendMessage(ctx, &bot.SendMessageParams{
							ChatID: subs[i].TelegramID,
							Text: mail.Text,
							ReplyMarkup: kb,
						})
					}
				}
			}
		}
	}
}

// runs all bots from DB
func InitBots(db *sqlx.DB, pool *models.GoroutinesPool) error {
	shopRepo := repos.NewShopRepo(db)
	subRepo := repos.NewSubscriberRepo(db)

	bots, err := shopRepo.Select()
	if err != nil {
		return err
	}

	for i := range bots {
		token := bots[i].Token

		ch := make(chan models.ChannelMessage)
		pool.Add(token, ch)

		go BotWorker(token, shopRepo, subRepo, pool)

		msg := models.WorkType(true)

		ch <- msg
	}

	return nil
}
