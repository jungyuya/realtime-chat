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

// Message는 클라이언트와 Hub 간에 전달되는 메시지의 구조를 정의합니다.
type Message struct {
	Content        string    `json:"content"`
	SenderID       string    `json:"senderId"`
	SenderNickname string    `json:"senderNickname"`
	Avatar         string    `json:"avatar"`
	Timestamp      time.Time `json:"timestamp"` // 타입을 time.Time으로 변경
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
	if !okId || !okNickname { // avatar에 대한 검사는 일단 제거
		log.Println("Invalid claims data type for id or nickname")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// avatar 클레임은 선택적으로 처리
	avatar, okAvatar := claims["avatar"].(string)
	if !okAvatar {
		avatar = "avatar-1" // 토큰에 avatar 정보가 없으면, 기본값으로 avatar-1을 사용
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

// writePump는 Hub로부터 메시지를 받아 WebSocket 연결로 전송합니다.
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("error writing message: %v", err)
			break
		}
	}
	log.Println("writePump finished for a client")
}
