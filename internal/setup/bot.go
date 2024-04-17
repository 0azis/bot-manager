package setup

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

var Goroutines = make(map[string]chan bool)

func GoroutineExists(token string) bool {
	_, ok := Goroutines[token]
	return ok
}

func BotWorker(token string, shopRepo repos.ShopRepo, subRepo repos.SubscriberRepo) {
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
		shopID, _ := shopRepo.ShopByBot(strconv.FormatInt(botInfo.ID, 10))

		s := models.Subscriber{
			UserName:   update.Message.From.Username,
			FirstName:  update.Message.From.FirstName,
			LastName:   update.Message.From.LastName,
			AvatarUrl:  "",
			TelegramID: tgID,
			ShopID:     shopID,
		}
		subRepo.Insert(s)

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   botData.FirstLaunch,
		})
	})

	// get channel that created
	channel := Goroutines[token]

	for {
		select {
		case msg := <-channel:
			switch msg {
			case false:
				close(channel)
				return
			case true:
				go b.Start(ctx)
			}
		}
	}
}

// runs all bots from DB
func InitBots(db *sqlx.DB) error {
	shopRepo := repos.NewShopRepo(db)
	subRepo := repos.NewSubscriberRepo(db)

	bots, err := shopRepo.Select()
	if err != nil {
		return err
	}

	for i := range bots {
		token := bots[i].Token

		ch := make(chan bool)
		Goroutines[token] = ch
		go BotWorker(token, shopRepo, subRepo)

		ch <- true
	}

	return nil
}
