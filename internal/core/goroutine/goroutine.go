package goroutine

import (
	"botmanager/internal/adapter/redis"
	"botmanager/internal/adapter/repo"
	"botmanager/internal/core/domain"
	"context"

	tg_domain "github.com/go-telegram/bot/models"

	"github.com/go-telegram/bot"
)

// 1. ShopBot - bot from DB, that implement user's shop
// 2. HomeBot - registration bot and home bot for users (Telegram CRM)

type goroutineInterface interface {
	Start()                                   // start bot
	Stop()                                    // stop bot
	SendMail(mail domain.Mail) error          // send mail to all subscribers
	SendMessage(message domain.Message) error // send message to the current user
	InitShopHandlers()                        // init handlers for shop bot
	InitHomeHandlers()                        // init handlers for home bot
	// SendNotification(notification domain.Notification)
}

type goroutine struct {
	token   string
	channel chan channelMessage
	bot     *bot.Bot
	ctx     context.Context
	pool    *GoroutinesPool
	store   repo.Store
	redisDB redis.RedisInterface
}

func NewShopBot(token string, pool *GoroutinesPool, store repo.Store) (goroutineInterface, error) {
	var goroutine goroutine

	goroutine.pool = pool
	goroutine.token = token
	goroutine.store = store

	opts := []bot.Option{
		bot.WithDefaultHandler(goroutine.listenMessages),
	}
	b, err := bot.New(token, opts...)
	if err != nil {
		return &goroutine, err
	}
	goroutine.bot = b

	ch := make(chan channelMessage)
	goroutine.channel = ch

	return &goroutine, err
}

func NewHomeBot(token string, pool *GoroutinesPool, store repo.Store, redisDB redis.RedisInterface) (goroutineInterface, error) {
	var goroutine goroutine

	goroutine.pool = pool
	goroutine.store = store
	goroutine.token = token
	goroutine.redisDB = redisDB

	b, err := bot.New(token)
	if err != nil {
		return &goroutine, err
	}
	goroutine.bot = b

	ch := make(chan channelMessage)
	goroutine.channel = ch

	return &goroutine, err
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

// Send mail for all subscribers
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
			kb.InlineKeyboard = append(kb.InlineKeyboard, []tg_domain.InlineKeyboardButton{{Text: "Открыть меню", WebApp: &tg_domain.WebAppInfo{URL: domain.WebLink(mail.ShopID)}}})
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

// Send message to current user
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

func (g goroutine) SendCode(userID string, code string) error {
	g.redisDB.SetCode(userID, code)
	_, err := g.bot.SendMessage(g.ctx, &bot.SendMessageParams{
		ChatID: userID,
		Text:   code,
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
