package engine

import (
	"sync"
	"github.com/felipetojal/broker/internal/ws"
)

type Broadcast interface {
	Send(message []byte)
}

type Exchange struct {
		Books map[string]*OrderBook
		mu sync.RWMutex
		broadcast Broadcast
}

func NewExchange(b Broadcast) *Exchange {
		return &Exchange{
				Books : make(map[string]*OrderBook),
				broadcast:  b,
		}
}

func (ex *Exchange) ProcessOrder(ticker string, ord *Order) {
		ex.mu.RLock()
		book, exists := ex.Books[ticker]
		ex.mu.RUnlock()

		if !exists {
				ex.mu.Lock()
				book, exists = ex.Books[ticker]
				if !exists {
						book = NewOrderBook(ex.broadcast)
						ex.Books[ticker]=book
				}
				ex.mu.Unlock()
		}
		book.AddOrder(ord)
}
