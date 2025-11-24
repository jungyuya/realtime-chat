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
	"github.com/jungyuya/realtime-chat/backend/internal/db" // [중요] db 패키지 임포트
	"golang.org/x/time/rate"
)

// [수정] Heartbeat 설정 최적화
const (
	// 클라이언트에게 Ping을 보내는 주기 (10초로 단축하여 연결 유지 강화)
	pingPeriod = 10 * time.Second

	// Pong 응답을 기다리는 최대 시간
	pongWait = 60 * time.Second

	// 메시지 쓰기 제한 시간
	writeWait = 10 * time.Second
)

type Message struct {
	Content        string    `json:"content"`
	SenderID       string    `json:"senderId"`
	SenderNickname string    `json:"senderNickname"`
	Avatar         string    `json:"avatar"`
	Timestamp      time.Time `json:"timestamp"`
}

type Client struct {
	hub         *Hub
	conn        *websocket.Conn
	send        chan []byte
	limiter     *rate.Limiter
	AnonymousID string
	Nickname    string
	Avatar      string
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("Client registered: %s (%s). Total: %d", client.Nickname, client.AnonymousID, len(h.clients))

			// [디버깅 로그 추가] DB 조회 시작 알림
			log.Println("DEBUG: Start fetching recent messages from DB...")

			recentMessages, err := db.GetRecentMessages("global", 50)

			// [디버깅 로그 추가] DB 조회 결과 알림
			log.Printf("DEBUG: Finished fetching. Error: %v, Count: %d", err, len(recentMessages))

			if err != nil {
				log.Printf("Failed to fetch recent messages for %s: %v", client.Nickname, err)
			} else {
				// [수정] 루프에 'SendingHistory'라는 이름을 붙입니다.
			SendingHistory:
				for _, dbMsg := range recentMessages {
					msg := &Message{
						Content:        dbMsg.Content,
						SenderID:       dbMsg.SenderID,
						SenderNickname: dbMsg.SenderNickname,
						Avatar:         dbMsg.Avatar,
						Timestamp:      dbMsg.CreatedAt,
					}

					messageBytes, err := json.Marshal(msg)
					if err != nil {
						log.Printf("Failed to marshal recent message: %v", err)
						continue
					}

					select {
					case client.send <- messageBytes:
					default:
						close(client.send)
						delete(h.clients, client)
						// [수정] 단순히 break가 아니라, 'SendingHistory' 루프를 종료하라고 명시합니다.
						break SendingHistory
					}
				}
			}

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client unregistered: %s. Total: %d", client.Nickname, len(h.clients))
			}

		case message := <-h.broadcast:
			// [기존] 메시지를 DB에 저장 (비동기)
			go func(msg *Message) {
				err := db.SaveMessage("global", msg.SenderID, msg.SenderNickname, msg.Avatar, msg.Content)
				if err != nil {
					log.Printf("Failed to save message to DB: %v", err)
				}
			}(message)

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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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
	avatar, okAvatar := claims["avatar"].(string)
	if !okAvatar {
		avatar = "avatar-1"
	}

	if !okId || !okNickname {
		log.Println("Invalid claims data type")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
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

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	// [수정] Pong 핸들러 및 ReadDeadline 설정
	c.conn.SetReadLimit(4096) // 메시지 크기 제한을 좀 더 넉넉하게 (4KB)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, content, err := c.conn.ReadMessage()
		if err != nil {
			// [추가] 에러 로그 상세화
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Client %s disconnected unexpectedly: %v", c.Nickname, err)
			} else {
				log.Printf("Client %s disconnected: %v", c.Nickname, err)
			}
			break
		}

		if err := c.limiter.Wait(context.Background()); err != nil {
			log.Printf("Rate limiter error for %s: %v", c.Nickname, err)
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

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
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

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Write error for %s: %v", c.Nickname, err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Ping error for %s: %v", c.Nickname, err)
				return
			}
		}
	}
}
