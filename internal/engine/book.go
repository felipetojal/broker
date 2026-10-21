package engine

import (
		"container/heap"
		"sync"
		"time"
		"fmt"
		"github.com/felipetojal/broker/internal/ws"
		"encoding/json"
)

type Clock interface {
		Now() time.Time
}

type OrderBook struct {
		SellPQ SellPriorityQueue
		BuyPQ BuyPriorityQueue
		Orders map[string]*Order
		Trades []Trade
		broadcast Broadcast
		mu sync.Mutex
}

func NewOrderBook(b Broadcast) *OrderBook {
		ob := &OrderBook{
				SellPQ : SellPriorityQueue{OrderPriorityQueue: make(OrderPriorityQueue,0)},
				BuyPQ  : BuyPriorityQueue{OrderPriorityQueue: make(OrderPriorityQueue,0)},
				Orders : make(map[string]*Order),
				broadcast: b,
		}

		heap.Init(&ob.BuyPQ)
		heap.Init(&ob.SellPQ)
		return ob
}


func (ob *OrderBook) AddOrder(ord *Order) {
		ob.mu.Lock()
		defer ob.mu.Unlock()

		ob.Orders[ord.Id] = ord

		if ord.OrderSide == Buy {
				heap.Push(&ob.BuyPQ, ord)
		}else if ord.OrderSide == Sell {
				heap.Push(&ob.SellPQ, ord)
		}

		ob.Match()

}

func (ob *OrderBook) CancellOrder(ord *Order){
		ob.mu.Lock()
		defer ob.mu.Unlock()

		o,exists := ob.Orders[ord.Id]

		if exists {
				o.Status=false
		}
}

func (ob *OrderBook) deleteOrder(ord *Order, id string) {
				delete(ob.Orders, id)
}

func (ob *OrderBook) broadcastTrade(t Trade){
		if ob.broadcast == nil {
				return
		}
		b, err := json.Marshal(t)
		if err == nil {
				ob.broadcast.Send(b)
		}
}

func (ob *OrderBook) Match(){
		for {
				if ob.BuyPQ.Len() == 0  || ob.SellPQ.Len() == 0 {
						break
				}
				topBuy := ob.BuyPQ.Top()
				topSell := ob.SellPQ.Top()

				if !topBuy.Status {
						heap.Pop(&ob.BuyPQ)
						continue
				}
				if !topSell.Status {
						heap.Pop(&ob.SellPQ)
						continue
				}

				if topBuy.Price < topSell.Price {
						break
				}

				if topBuy.Quantity > topSell.Quantity {
						topBuy.Quantity -= topSell.Quantity
						t := Trade{
								IdBuy: topBuy.Id,
								IdSell: topSell.Id,
								Price : topBuy.Price,
								Quantity: topSell.Quantity,
								TradeAt: time.Now(),
						}
						ob.Trades = append(ob.Trades, t)
						ob.broadcastTrade(t)
						ob.deleteOrder(topSell, topSell.Id)
						heap.Pop(&ob.SellPQ)
				} else if topBuy.Quantity < topSell.Quantity {
						topSell.Quantity -= topBuy.Quantity
						t := Trade{
								IdBuy: topBuy.Id,
								IdSell: topSell.Id,
								Price : topBuy.Price,
								Quantity: topBuy.Quantity,
								TradeAt: time.Now(),
						}
						ob.Trades = append(ob.Trades, t)
						ob.broadcastTrade(t)
						ob.deleteOrder(topBuy, topBuy.Id)
						heap.Pop(&ob.BuyPQ)
				} else{
						ob.deleteOrder(topSell, topSell.Id)
						ob.deleteOrder(topBuy, topBuy.Id)
						t := Trade{
								IdBuy: topBuy.Id,
								IdSell: topSell.Id,
								Price : topBuy.Price,
								Quantity: topSell.Quantity,
								TradeAt: time.Now(),
						}
						ob.Trades = append(ob.Trades, t)
						ob.broadcastTrade(t)
						heap.Pop(&ob.SellPQ)
						heap.Pop(&ob.BuyPQ)
				}
		}
}

func (ob *OrderBook) PrintState(ticker string) {
		ob.mu.Lock()
		defer ob.mu.Unlock()

		fmt.Printf("\n========== LIVRO DE OFERTAS: %s ==========\n", ticker)

		fmt.Println("\n[ TRADES EXECUTADOS ]")
		fmt.Println("-------------------------------------------------")
		if len(ob.Trades) == 0 {
				fmt.Println("Nenhum trade executado.")
		} else {
				for _, t := range ob.Trades {
						fmt.Printf("✅ Negócio Fechado: %d ações a R$ %d | Comprador: %s -> Vendedor: %s\n",
								t.Quantity, t.Price, t.IdBuy, t.IdSell)
				}
		}

		fmt.Println("\n[ FILA DE VENDAS (ASKS) ] - Menor preço no topo")
		fmt.Println("-------------------------------------------------")
		if ob.SellPQ.Len() == 0 {
				fmt.Println("Fila vazia.")
		} else {
				for _, o := range ob.SellPQ.OrderPriorityQueue {
						if o.Status {
								fmt.Printf("🔴 Venda | ID: %s \t| Preço: R$ %d \t| Qtd: %d\n", o.Id, o.Price, o.Quantity)
						}
				}
		}

		fmt.Println("\n[ FILA DE COMPRAS (BIDS) ] - Maior preço no topo")
		fmt.Println("-------------------------------------------------")
		if ob.BuyPQ.Len() == 0 {
				fmt.Println("Fila vazia.")
		} else {
				for _, o := range ob.BuyPQ.OrderPriorityQueue {
						if o.Status {
								fmt.Printf("🔵 Compra | ID: %s \t| Preço: R$ %d \t| Qtd: %d\n", o.Id, o.Price, o.Quantity)
						}
				}
		}
		fmt.Println("=================================================\n")
}
