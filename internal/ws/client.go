package ws

import (
		"bytes"
		"encoding/json"
		"log"
		"net/http"
		"strconv"
		"time"
		"github.com/felipetojal/broker/internal/engine"
		"github.com/gorilla/websocket"
)

const (
		writeWait = 10 * time.Second
		pongWait = 60 * time.Second
		pingPeriord = (pongWait * 9)/10
		maxMessageSize=512
)

var (
		newline= []byte{'\n'}
		space=[]byte{' '}
)

var upgrader = websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
				return true
		},
}

type Client struct {
		hub *Hub
		conn *websocket.Conn
		ex *engine.Exchange
		send chan []byte
}

func (c *Client) readPump(){
		defer func() {
				c.hub.unregister <- c
				c.conn.Close()
		}()

		c.conn.SetReadLimit(maxMessageSize)
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		c.conn.SetPongHandler(func(string) error {
				c.conn.SetReadDeadline(time.Now().Add(pongWait));
				return nil
		})

		for{
				_, message, err := c.conn.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
								log.Printf("error: %v", err)
						}
						break
				}
				message=bytes.TrimSpace(bytes.Replace(message,newline,space,-1))
				var orders []engine.Order
				er := json.Unmarshal(message, &orders)
				if er != nil {
						log.Println("Error: %v", er)
						break;
				}

				for _, order := range orders {
						c.ex.ProcessOrder(order.Ticker, &engine.Order{
								Id : strconv.FormatInt(order.IdUser, 10) + strconv.FormatInt(time.
Now().UnixNano(), 10),
								Price: order.Price,
								Quantity: order.Quantity,
								In : time.Hour.Milliseconds(),
								CreateAt: time.Now(),
								OrderSide: order.OrderSide,
								Status:  true,
								IdUser: order.IdUser,
								Ticker: order.Ticker,
						})
				}
		}
}


func (c *Client) writePump() {
		ticker := time.NewTicker(pingPeriord)
		defer func() {
				ticker.Stop()
				c.conn.Close()
		}()

		for {
				select {
				case message, ok := <-c.send:
						c.conn.SetWriteDeadline(time.Now().Add(writeWait))
						if !ok {
								c.conn.WriteMessage(websocket.CloseMessage, []byte{})
								return
						}
						w,err := c.conn.NextWriter(websocket.TextMessage)
						if err != nil {
								return
						}
						w.Write(message)

						n := len(c.send)
						for i:=0; i < n; i++ {
								w.Write(newline)
								w.Write(<-c.send)
						}

						if err:= w.Close(); err != nil {
								return
						}

				case <-ticker.C:
						c.conn.SetWriteDeadline(time.Now().Add(writeWait))
						if err := c.conn.WriteMessage(websocket.PingMessage,nil); err != nil {
								return
						}
				}
		}
}


func SeveWs(hub *Hub, ex *engine.Exchange, w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w,r,nil)

		if err != nil {
				log.Println(err)
				return
		}

		client := &Client{
				hub: hub,
				conn: conn,
				ex: ex,
				send : make(chan []byte,256),
		}
		client.hub.register <- client

		go client.writePump()
		go client.readPump()
}
