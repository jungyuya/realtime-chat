package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"golang.org/x/time/rate"
)

const (
	// 클라이언트에게 Ping을 보내는 주기 (반드시 pongWait보다 짧아야 함)
	pingPeriod = 50 * time.Second

	// Pong 응답을 기다리는 최대 시간
	pongWait = 60 * time.Second

	// 메시지 쓰기 제한 시간
	writeWait = 10 * time.Second
)

// Message는 클라이언트와 Hub 간에 전달되는 메시지의 구조를 정의합니다.
type Message struct {
	Content        string    `json:"content"`
	SenderID       string    `json:"senderId"`
	SenderNickname string    `json:"senderNickname"`
	Avatar         string    `json:"avatar"`
	Timestamp      time.Time `json:"timestamp"`
}

// Client는 Hub와 WebSocket 연결 사이의 중개자 역할을 합니다.
type Client struct {
	hub         *Hub
	conn        *websocket.Conn
	send        chan []byte
	limiter     *rate.Limiter
	AnonymousID string
	Nickname    string
	Avatar      string
}

// Hub는 모든 활성 클라이언트를 관리하고 메시지를 모든 클라이언트에게 브로드캐스트합니다.
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
}

// NewHub는 새로운 Hub 인스턴스를 생성하고 초기화합니다.
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Run은 Hub를 별도의 고루틴으로 실행합니다.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Println("New client registered. Total clients:", len(h.clients))

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Println("Client unregistered. Total clients:", len(h.clients))
			}

		case message := <-h.broadcast:
			messageBytes, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshalling message: %v", err)
				continue
			}
			for client := range h.clients {
				select {
				case client.send <- messageBytes:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// WebSocket 업그레이더 설정
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ServeWs는 WebSocket 요청을 처리하고 인증을 수행합니다.
func ServeWs(hub *Hub, c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		log.Println("Token not found in query")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	secretKey := os.Getenv("SECRET_KEY")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		log.Printf("Token parsing error: %v", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var claims jwt.MapClaims
	var ok bool
	if claims, ok = token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		log.Println("Invalid token or claims")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	anonymousId, okId := claims["anonymousId"].(string)
	nickname, okNickname := claims["nickname"].(string)
	if !okId || !okNickname {
		log.Println("Invalid claims data type for id or nickname")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// avatar 클레임은 선택적으로 처리
	avatar, okAvatar := claims["avatar"].(string)
	if !okAvatar {
		avatar = "avatar-1"
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	limiter := rate.NewLimiter(rate.Limit(0.5), 5)
	client := &Client{
		hub:         hub,
		conn:        conn,
		send:        make(chan []byte, 256),
		limiter:     limiter,
		AnonymousID: anonymousId,
		Nickname:    nickname,
		Avatar:      avatar,
	}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

// readPump는 WebSocket 연결에서 메시지를 읽어 Hub로 전달합니다.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	// [수정] Heartbeat 설정: 읽기 제한 시간 및 Pong 핸들러 설정
	c.conn.SetReadLimit(512) // 메시지 최대 크기 제한 (선택 사항)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		// Pong을 받으면 제한 시간을 다시 늘려줌 (연결 유지)
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, content, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		if err := c.limiter.Wait(context.Background()); err != nil {
			log.Printf("rate limiter wait error: %v", err)
			break
		}

		msg := &Message{
			Content:        string(content),
			SenderID:       c.AnonymousID,
			SenderNickname: c.Nickname,
			Avatar:         c.Avatar,
			Timestamp:      time.Now(),
		}
		c.hub.broadcast <- msg
	}
}

// writePump는 Hub로부터 메시지를 받아 WebSocket 연결로 전송하고, 주기적으로 Ping을 보냅니다.
func (c *Client) writePump() {
	// [수정] 주기적으로 신호를 보내는 Ticker 생성
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			// [수정] 쓰기 제한 시간 설정
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub가 채널을 닫으면 연결 종료 메시지 전송
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// 메시지 전송
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("error writing message: %v", err)
				return
			}

		// [추가] Ticker가 울릴 때마다 Ping 메시지 전송
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}