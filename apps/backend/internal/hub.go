// 'internal' 패키지에 속함을 선언합니다.
package internal

import (
	"context"
	"encoding/json" // json 패키지 추가
	"fmt"           // fmt 패키지 추가
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5" // jwt 패키지 추가
	"github.com/gorilla/websocket"
	"golang.org/x/time/rate" // rate 패키지 추가
	"log"
	"net/http"
	"os" // os 패키지 추가
)

// Message는 클라이언트와 Hub 간에 전달되는 메시지의 구조를 정의합니다.
type Message struct {
	Content        string `json:"content"`
	SenderID       string `json:"senderId"`
	SenderNickname string `json:"senderNickname"`
	// Avatar         string `json:"avatar"` // 아바타는 다음 단계에서 추가
}

// Client 구조체 확장
type Client struct {
	hub         *Hub
	conn        *websocket.Conn
	send        chan []byte
	limiter     *rate.Limiter
	AnonymousID string // 사용자의 고유 ID
	Nickname    string // 사용자의 닉네임
	// Avatar      string // 아바타
}

// Hub는 모든 활성 클라이언트를 관리하고 메시지를 모든 클라이언트에게 브로드캐스트합니다.
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan *Message // chan []byte에서 변경
	register   chan *Client
	unregister chan *Client
}

// NewHub는 새로운 Hub 인스턴스를 생성하고 초기화합니다.
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Message), // 타입 변경
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// WebSocket 업그레이더 설정 (main.go에서 이동)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ServeWs는 WebSocket 요청을 처리합니다.
// Hub 인스턴스를 받아와서 새로운 클라이언트를 생성하고 등록합니다.
func ServeWs(hub *Hub, c *gin.Context) {
	// 1. URL 쿼리에서 토큰 추출
	tokenString := c.Query("token")
	if tokenString == "" {
		log.Println("Token not found in query")
		return
	}

	// 2. 토큰 파싱 및 검증
	secretKey := os.Getenv("SECRET_KEY")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 서명 알고리즘이 HMAC인지 확인
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		log.Printf("Token parsing error: %v", err)
		return
	}

	var claims jwt.MapClaims
	var ok bool
	if claims, ok = token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		log.Println("Invalid token or claims")
		return
	}

	// 3. 클레임에서 사용자 정보 추출
	anonymousId, okId := claims["anonymousId"].(string)
	nickname, okNickname := claims["nickname"].(string)
	if !okId || !okNickname {
		log.Println("Invalid claims data type")
		return
	}

	// 4. WebSocket 연결 업그레이드 (인증 성공 후)
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// 5. Client 인스턴스 생성 및 정보 채우기
	limiter := rate.NewLimiter(rate.Limit(0.5), 5)
	client := &Client{
		hub:         hub,
		conn:        conn,
		send:        make(chan []byte, 256),
		limiter:     limiter,
		AnonymousID: anonymousId,
		Nickname:    nickname,
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
		_, content, err := c.conn.ReadMessage() // 받은 메시지는 순수 텍스트 내용(content)
		if err != nil {
			// ... 기존 에러 처리 ...
			break
		}

		if err := c.limiter.Wait(context.Background()); err != nil {
			// ... 기존 에러 처리 ...
			break
		}

		// Message 구조체 생성
		msg := &Message{
			Content:        string(content),
			SenderID:       c.AnonymousID,
			SenderNickname: c.Nickname,
		}
		// 구조체를 broadcast 채널로 전송
		c.hub.broadcast <- msg
	}
}

func (h *Hub) Run() {
	for {
		select {
		// ... register, unregister case는 변경 없음 ...
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
			// Message 구조체를 JSON으로 변환
			messageBytes, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshalling message: %v", err)
				continue
			}

			for client := range h.clients {
				select {
				case client.send <- messageBytes: // JSON 바이트를 전송
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// writePump는 Hub로부터 메시지를 받아 WebSocket 연결로 전송합니다.
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	// c.send 채널에 대한 for range 루프입니다.
	// 이 루프는 c.send 채널이 닫힐 때까지 계속 실행됩니다.
	for message := range c.send {
		// 채널에서 메시지를 성공적으로 꺼내오면,
		// 해당 메시지를 WebSocket 연결을 통해 클라이언트에게 전송합니다.
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			// 메시지 전송 중 에러가 발생하면 로그를 남기고 루프를 종료합니다.
			// defer에 의해 conn.Close()가 호출됩니다.
			log.Printf("error writing message: %v", err)
			break
		}
	}
	// 루프가 종료되었다는 것은 채널이 닫혔거나 전송 에러가 발생했다는 의미입니다.
	// defer가 연결을 정리해 줄 것입니다.
	log.Println("writePump finished for a client")
}
