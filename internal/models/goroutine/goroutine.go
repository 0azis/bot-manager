package goroutine

import (
	"botmanager/internal/models"
	"botmanager/internal/repos"
	"context"

	tg_models "github.com/go-telegram/bot/models"

	"github.com/go-telegram/bot"
)

type GoroutineInterface interface {
	Start()
	Stop()
	SendMessages(mail models.Mail) error
}

type Goroutine struct {
	channel chan ChannelMessage
	botData models.Shop
	bot     *bot.Bot
	ctx     context.Context
	pool    *GoroutinesPool
	store   repos.Store
}

func New(botData models.Shop, store repos.Store, pool *GoroutinesPool) (*Goroutine, error) {
	var goroutine Goroutine	

	goroutine.store = store
	goroutine.pool = pool
	goroutine.botData = botData

	b, err := bot.New(botData.Token)
	if err != nil {
		return &goroutine, err
	}
	goroutine.bot = b

	ch := make(chan ChannelMessage)
	goroutine.channel = ch

	goroutine.initHandlers()

	return &goroutine, nil
}

// Start the goroutine
func (g *Goroutine) Start() {
	g.pool.Add(g)

	go g.run()

	msg := StatusWork(true)
	g.channel <- msg
}

// Stop the goroutine
func (g *Goroutine) Stop() {
	msg := StatusWork(false)
	g.channel <- msg

	g.pool.Delete(g)
}

// Send message for all subscribers
func (g Goroutine) SendMessages(mail models.Mail) error {
	subs, err := g.store.Subscriber().Select(mail.ShopID)
	if err != nil {
		return err
	}

	for i := range subs {
		kb := &tg_models.InlineKeyboardMarkup{
			InlineKeyboard: [][]tg_models.InlineKeyboardButton{},
		}

		if mail.AddButton {
			kb.InlineKeyboard = append(kb.InlineKeyboard, []tg_models.InlineKeyboardButton{{Text: "Открыть меню", WebApp: &tg_models.WebAppInfo{URL: models.WebLink(g.botData.ID)}}})
		}

		if mail.PhotoLink != "" {
			g.bot.SendPhoto(g.ctx, &bot.SendPhotoParams{
				ChatID:      subs[i].TelegramID,
				Photo:       &tg_models.InputFileString{Data: "https://tgrocket.ru/api/uploads/share/" + mail.PhotoLink},
				Caption:     mail.Text,
				ReplyMarkup: kb,
			})
		} else {
			g.bot.SendMessage(g.ctx, &bot.SendMessageParams{
				ChatID:      subs[i].TelegramID,
				Text:        mail.Text,
				ReplyMarkup: kb,
			})
		}
	}

	return nil
}

func (g *Goroutine) run() {
	ctx, done := context.WithCancel(context.Background()) 
	defer done()
	g.ctx = ctx

	for {
		select {
		case msg := <-g.channel:
			switch msg.MsgType {
			case "status":
				switch msg.Value {
				case false:
					close(g.channel)
					return	
				case true:
					go g.bot.Start(g.ctx)
				}
			case "mail":
				mail := msg.Value.(models.Mail)
				g.SendMessages(mail)	
			}
		}
	}
}
