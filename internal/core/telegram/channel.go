package telegram

import (
	"botmanager/internal/core/domain"
)

type GoroutinesPool struct {
	pool []*goroutine
}

func NewPool() *GoroutinesPool {
	gp := new(GoroutinesPool)
	return gp
}

func (g GoroutinesPool) Exists(ID string) bool {
	for goroutine := range g.pool {
		if g.pool[goroutine].botData.ID == ID{
			return true
		}
	}
	return false
}

func (g GoroutinesPool) Get(ID string) goroutineInterface {
	for goroutine := range g.pool {
		if g.pool[goroutine].botData.ID == ID {
			return g.pool[goroutine]
		}
	}
	return nil
}

func (g *GoroutinesPool) Add(goroutine *goroutine) {
	g.pool = append(g.pool, goroutine)
}

func (g *GoroutinesPool) Delete(goroutine *goroutine) {
	for i := range g.pool {
		if g.pool[i] == goroutine {
			g.pool = append(g.pool[:i], g.pool[i+1:]...)
			break
		}
	}
}

type ChannelMessage struct {
	MsgType string
	Value   any
}

// Return message for channel, that start the bot
func StatusWork(status bool) ChannelMessage {
	msg := ChannelMessage{
		MsgType: "status",
		Value:   status,
	}
	return msg
}

// Return message for channel, that send mail
func SendMail(mail domain.Mail) ChannelMessage {
	msg := ChannelMessage{
		MsgType: "mail",
		Value:   mail,
	}
	return msg
}

// Return message for channel, that send message
func SendMessage(message domain.Message) ChannelMessage {
	msg := ChannelMessage{
		MsgType: "message",
		Value:   message,
	}
	return msg
}
