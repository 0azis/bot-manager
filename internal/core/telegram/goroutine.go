package telegram

import (
	"botmanager/internal/adapter/repo"
	"botmanager/internal/core/domain"
	"context"

	tg_domain "github.com/go-telegram/bot/models"

	"github.com/go-telegram/bot"
)

type goroutineInterface interface {
	Start()
	Stop()
	SendMail(mail domain.Mail) error
	SendMessage(message domain.Message) error
}

type goroutine struct {
	channel chan ChannelMessage
	botData domain.Shop
	bot     *bot.Bot
	ctx     context.Context
	pool    *GoroutinesPool
	store   repo.Store
}

func New(botData domain.Shop, store repo.Store, pool *GoroutinesPool) (goroutineInterface, error) {
	var goroutine goroutine

	goroutine.store = store
	goroutine.pool = pool
	goroutine.botData = botData

	opts := []bot.Option{
		bot.WithDefaultHandler(goroutine.listenMessages),
	}

	b, err := bot.New(botData.Token, opts...)
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
func (g *goroutine) Start() {
	g.pool.Add(g)

	go g.run()

	msg := StatusWork(true)
	g.channel <- msg
}

// Stop the goroutine
func (g *goroutine) Stop() {
	msg := StatusWork(false)
	g.channel <- msg

	g.pool.Delete(g)
}

// Send message for all subscribers
func (g goroutine) SendMail(mail domain.Mail) error {
	subs, err := g.store.Subscriber.Select(mail.ShopID)
	if err != nil {
		return err
	}

	for i := range subs {
		kb := &tg_domain.InlineKeyboardMarkup{
			InlineKeyboard: [][]tg_domain.InlineKeyboardButton{},
		}

		if mail.AddButton {
			kb.InlineKeyboard = append(kb.InlineKeyboard, []tg_domain.InlineKeyboardButton{{Text: "Открыть меню", WebApp: &tg_domain.WebAppInfo{URL: domain.WebLink(g.botData.ID)}}})
		}

		if mail.PhotoLink != "" {
			g.bot.SendPhoto(g.ctx, &bot.SendPhotoParams{
				ChatID:      subs[i].TelegramID,
				Photo:       &tg_domain.InputFileString{Data: "https://tgrocket.ru/api/uploads/share/" + mail.PhotoLink},
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

func (g goroutine) SendMessage(message domain.Message) error {
	sub, err := g.store.Subscriber.Get(message.SubscriberID)
	if err != nil {
		return err
	}

	_, err = g.bot.SendMessage(g.ctx, &bot.SendMessageParams{
		ChatID: sub.TelegramID,
		Text:   message.Text,
	})

	return err
}

func (g *goroutine) run() {
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
				mail := msg.Value.(domain.Mail)
				g.SendMail(mail)
			case "message":
				message := msg.Value.(domain.Message)
				g.SendMessage(message)
			}
		}
	}
}
