package models

// {"token": *channel} 
// basic message form
type GoroutinesPool struct {
	pool map[string]chan ChannelMessage
}

func NewGoroutinesPool() *GoroutinesPool {
	return &GoroutinesPool{
		pool: make(map[string]chan ChannelMessage),
	}
}

func (g GoroutinesPool) Exists(token string) bool {
	_, ok := g.pool[token]
	return ok
}

func (g GoroutinesPool) Get(token string) chan ChannelMessage {
	return g.pool[token] 
}

func (g *GoroutinesPool) Add(token string, value chan ChannelMessage) {
	g.pool[token] = value
}

func (g *GoroutinesPool) Delete(token string) {
	delete(g.pool, token)
}

type ChannelMessage struct {
	MsgType string
	Value any 
} 

func WorkType(status bool) ChannelMessage {
	msg := ChannelMessage{
		MsgType: "status", 	
		Value: status,
	}
	return msg
}

func MailType(mail Mail) ChannelMessage {
	msg := ChannelMessage{
		MsgType: "mail",
		Value: mail,
	}
	return msg
}
