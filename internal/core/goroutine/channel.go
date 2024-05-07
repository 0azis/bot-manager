package goroutine

type GoroutinesPool struct {
	pool []*goroutine

	// homebot token
	HomeBotToken string
}

func NewPool(hbtoken string) *GoroutinesPool {
	return &GoroutinesPool{
		HomeBotToken: hbtoken,
	}
}

func (g GoroutinesPool) Exists(token string) bool {
	for goroutine := range g.pool {
		if g.pool[goroutine].token == token {
			return true
		}
	}
	return false
}

func (g GoroutinesPool) Get(token string) *goroutine {
	for goroutine := range g.pool {
		if g.pool[goroutine].token == token {
			return g.pool[goroutine]
		}
	}
	return nil
}

func (g GoroutinesPool) GetHomeBot() *goroutine {
	return g.Get(g.HomeBotToken)
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
