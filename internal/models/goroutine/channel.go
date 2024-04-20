package goroutine

import (
	"botmanager/internal/models"
)

type GoroutinesPool struct {
	pool []*Goroutine
}

func NewPool() *GoroutinesPool {
	gp := new(GoroutinesPool)
	return gp
}

func (g GoroutinesPool) Exists(token string) bool {
	for goroutine := range g.pool {
		if g.pool[goroutine].botData.Token == token {
			return true
		}
	}
	return false
}

func (g GoroutinesPool) Get(token string) *Goroutine {
	for goroutine := range g.pool {
		if g.pool[goroutine].botData.Token == token {
			return g.pool[goroutine]
		}
	}
	return nil 
}

func (g *GoroutinesPool) Add(goroutine *Goroutine) {
	g.pool = append(g.pool, goroutine)
}

func (g *GoroutinesPool) Delete(goroutine *Goroutine) {
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
func SendMail(mail models.Mail) ChannelMessage {
	msg := ChannelMessage{
		MsgType: "mail",
		Value:   mail,
	}
	return msg
}
