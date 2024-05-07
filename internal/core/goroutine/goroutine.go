package goroutine

import (
	"botmanager/internal/adapter/secondary/database"
	"botmanager/internal/adapter/secondary/redis"
	"botmanager/internal/core/domain"
	"context"
	"fmt"

	"github.com/0azis/bot/models"

	"github.com/0azis/bot"
)

// 1. ShopBot - bot from DB, that implement user's shop
// 2. HomeBot - registration bot and home bot for users (Telegram CRM)

type goroutineInterface interface {
	Start()                                                  // start bot
	Stop()                                                   // stop bot
	SendMail(mail domain.Mail) error                         // send mail to all subscribers
	SendMessage(message domain.Message) error                // send message to the current user
	SendNotification(notification domain.Notification) error // send notification to home bot
	InitShopHandlers()                                       // init handlers for shop bot
	InitHomeHandlers()                                       // init handlers for home bot
}

type goroutine struct {
	token   string
	channel chan bool
	bot     *bot.Bot
	ctx     context.Context
	pool    *GoroutinesPool
	store   database.Store
	redisDB redis.RedisInterface
}

func New(token string, pool *GoroutinesPool, store database.Store, redisDB redis.RedisInterface) (goroutineInterface, error) {
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

	ch := make(chan bool)
	goroutine.channel = ch

	return &goroutine, err
}

// Start the goroutine
func (g *goroutine) Start() {
	g.pool.Add(g)

	go g.run()

	g.channel <- true
}

// Stop the goroutine
func (g *goroutine) Stop() {
	g.channel <- false

	g.pool.Delete(g)
}

// Send mail for all subscribers
func (g goroutine) SendMail(mail domain.Mail) error {
	subs, err := g.store.Subscriber.Select(mail.ShopID)
	if err != nil {
		return err
	}

	for i := range subs {
		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{},
		}

		if mail.AddButton {
			kb.InlineKeyboard = append(kb.InlineKeyboard, []models.InlineKeyboardButton{{Text: "Открыть меню", WebApp: &models.WebAppInfo{URL: domain.WebLink(mail.ShopID)}}})
		}

		if mail.PhotoLink != "" {
			g.bot.SendPhoto(g.ctx, &bot.SendPhotoParams{
				ChatID:      subs[i].TelegramID,
				Photo:       &models.InputFileString{Data: "https://tgrocket.ru/api/uploads/share/" + mail.PhotoLink},
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

func (g goroutine) SendNotification(notification domain.Notification) error {
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{},
	}

	if notification.Button {
		kb.InlineKeyboard = append(kb.InlineKeyboard, []models.InlineKeyboardButton{{Text: "Открыть меню", WebApp: &models.WebAppInfo{URL: domain.WebLink(notification.ShopID)}}})
	}

	if notification.Photo != "" {
		g.bot.SendPhoto(g.ctx, &bot.SendPhotoParams{
			ChatID:      notification.UserID,
			Photo:       &models.InputFileString{Data: notification.Photo},
			Caption:     notification.Text + "\n\n" + fmt.Sprintf("<a href='%s'>%s</a>", notification.Link, notification.LinkText),
			ReplyMarkup: kb,
			ParseMode:   "HTML",
		})
	} else {
		g.bot.SendMessage(g.ctx, &bot.SendMessageParams{
			ChatID:      notification.UserID,
			Text:        notification.Text + "\n\n" + fmt.Sprintf("<a href='%s'>%s</a>", notification.Link, notification.LinkText),
			ReplyMarkup: kb,
			ParseMode:   "HTML",
		})
	}
	return nil
}

func (g *goroutine) run() {
	ctx, done := context.WithCancel(context.Background())
	defer done()
	g.ctx = ctx

	for {
		select {
		case msg := <-g.channel:
			switch msg {
			case false:
				close(g.channel)
				return
			case true:
				go g.bot.Start(g.ctx)
			}

		}
	}
}
