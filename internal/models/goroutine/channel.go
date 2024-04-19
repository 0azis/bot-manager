package goroutine

import "botmanager/internal/models"

// {"token": *channel}
// basic message form
type GoroutinesPool []Goroutine 

func NewPool() GoroutinesPool {
	var gp GoroutinesPool
	return gp
}  

func (g GoroutinesPool) Exists(token string) bool {
	for goroutine := range g {
		if g[goroutine].botData.Token == token {
			return true
		}	
	}
	return false
}

func (g GoroutinesPool) Get(token string) Goroutine {
	for goroutine := range g {
		if g[goroutine].botData.Token == token {
			return g[goroutine]
		}
	}
	return Goroutine{} 
}

func (g GoroutinesPool) Add(goroutine *Goroutine) {
	g = append(g, *goroutine)
}

func (g GoroutinesPool) Delete(token string) {
	var delEl int
	for goroutine := range g {
		if g[goroutine].botData.Token == token {
			delEl = goroutine
			g = append(g[:delEl], g[delEl+1:]...)
		}
	}
}

type ChannelMessage struct {
	MsgType string
	Value any 
} 

// Return message for channel, that start the bot
func StatusWork(status bool) ChannelMessage {
	msg := ChannelMessage{
		MsgType: "status", 	
		Value: status,
	}
	return msg
}

// Return message for channel, that send mail 
func SendMail(mail models.Mail) ChannelMessage {
	msg := ChannelMessage{
		MsgType: "mail",
		Value: mail,
	}
	return msg
}
