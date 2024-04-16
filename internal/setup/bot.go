package setup

import (
	"botmanager/internal/models"
	"botmanager/internal/repos"
	"context"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	tg_models "github.com/go-telegram/bot/models"
	"github.com/jmoiron/sqlx"
)

var Goroutines = make(map[string]chan bool)

func GoroutineExists(token string) bool {
	_, ok := Goroutines[token]
	return ok
}

func BotWorker(token string, repo repos.ShopRepo) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	botData, _ := repo.Get(token)

	opts := []bot.Option{
		bot.WithDefaultHandler(func (ctx context.Context, b *bot.Bot, update *tg_models.Update) {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   update.Message.Text,
			})	
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

	// get channel that created
	channel := Goroutines[token]

	for {
		select {
		case msg := <-channel:
			switch msg {
			case false:
				return
			case true:
				go b.Start(ctx)
			}
		}
	}
}

// runs all bots from DB
func InitBots(db *sqlx.DB) error {
	repo := repos.NewShopRepo(db)
	bots, err := repo.Select()
	if err != nil {
		return err
	}	
	
	for i := range bots {
		token := bots[i].Token

		ch := make(chan bool)
		Goroutines[token] = ch
		go BotWorker(token, repo)

		ch <- true
	}	

	return nil 
}


