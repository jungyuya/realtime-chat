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
	"github.com/jungyuya/realtime-chat/backend/internal/db"
	"golang.org/x/time/rate"
)

const (
	pingPeriod = 10 * time.Second
	pongWait   = 60 * time.Second
	writeWait  = 10 * time.Second
)

type Message struct {
	Content        string    `json:"content"`
	SenderID       string    `json:"senderId"`
	SenderNickname string    `json:"senderNickname"`
	Avatar         string    `json:"avatar"`
	Timestamp      time.Time `json:"timestamp"`
	RoomID         string    `json:"roomId"` // [추가] 메시지가 속한 방 ID
}

type Client struct {
	hub         *Hub
	conn        *websocket.Conn
	send        chan []byte
	limiter     *rate.Limiter
	AnonymousID string
	Nickname    string
	Avatar      string
	RoomID      string // [추가] 클라이언트가 현재 접속한 방 ID
}

type Hub struct {
	// [수정] 방 ID별로 클라이언트 맵을 관리 (RoomID -> Client -> bool)
	rooms      map[string]map[*Client]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		rooms:      make(map[string]map[*Client]bool), // 초기화
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// [수정] 해당 방이 없으면 생성
			if h.rooms[client.RoomID] == nil {
				h.rooms[client.RoomID] = make(map[*Client]bool)
			}
			// 해당 방에 클라이언트 등록
			h.rooms[client.RoomID][client] = true
			
			log.Printf("Client registered to room [%s]: %s. Total in room: %d", client.RoomID, client.Nickname, len(h.rooms[client.RoomID]))

			// [수정] 해당 방(client.RoomID)의 최근 메시지 불러오기
			recentMessages, err := db.GetRecentMessages(client.RoomID, 50)
			if err != nil {
				log.Printf("Failed to fetch recent messages for %s: %v", client.Nickname, err)
			} else {
			SendingHistory:
				for _, dbMsg := range recentMessages {
					msg := &Message{
						Content:        dbMsg.Content,
						SenderID:       dbMsg.SenderID,
						SenderNickname: dbMsg.SenderNickname,
						Avatar:         dbMsg.Avatar,
						Timestamp:      dbMsg.CreatedAt,
						RoomID:         dbMsg.RoomID,
					}

					messageBytes, err := json.Marshal(msg)
					if err != nil {
						continue
					}

					select {
					case client.send <- messageBytes:
					default:
						close(client.send)
						delete(h.rooms[client.RoomID], client)
						break SendingHistory
					}
				}
			}

		case client := <-h.unregister:
			// [수정] 해당 방에서 클라이언트 제거
			if clients, ok := h.rooms[client.RoomID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)
					
					// 방이 비었으면 방 자체를 삭제 (메모리 관리)
					if len(clients) == 0 {
						delete(h.rooms, client.RoomID)
					}
					log.Printf("Client unregistered from room [%s]: %s", client.RoomID, client.Nickname)
				}
			}

		case message := <-h.broadcast:
			// [수정] 메시지를 해당 방(message.RoomID)의 DB에 저장
			go func(msg *Message) {
				err := db.SaveMessage(msg.RoomID, msg.SenderID, msg.SenderNickname, msg.Avatar, msg.Content)
				if err != nil {
					log.Printf("Failed to save message to DB: %v", err)
				}
			}(message)

			messageBytes, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshalling message: %v", err)
				continue
			}

			// [수정] 해당 방에 있는 클라이언트들에게만 전송
			if clients, ok := h.rooms[message.RoomID]; ok {
				for client := range clients {
					select {
					case client.send <- messageBytes:
					default:
						close(client.send)
						delete(clients, client)
					}
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
	// [추가] URL 쿼리에서 room 파라미터 추출 (기본값: global)
	roomID := c.DefaultQuery("room", "global")

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
		RoomID:      roomID, // [추가] Client에 RoomID 설정
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
	
	c.conn.SetReadLimit(4096)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { 
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil 
	})

	for {
		_, content, err := c.conn.ReadMessage()
		if err != nil {
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
			RoomID:         c.RoomID, // [추가] 메시지에 RoomID 포함
		}
		c.hub.broadcast <- msg
	}
}

// ... writePump는 변경 없음 ...
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