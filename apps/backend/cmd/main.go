package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/jungyuya/realtime-chat/backend/internal"
	"github.com/jungyuya/realtime-chat/backend/internal/db"
)

func main() {
	// 1. 환경 변수 로드
	// 로컬 개발 환경에서는 .env 파일을 읽고,
	// 컨테이너/클라우드 환경에서는 시스템 환경 변수를 사용합니다.
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, reading from environment variables")
	}

	// 2. 데이터베이스 초기화
	// DB_HOST 환경 변수가 있을 때만 DB 연결을 시도합니다.
	if os.Getenv("DB_HOST") != "" {
		db.InitDB()
	} else {
		log.Println("DB_HOST not set, skipping database initialization")
	}

	// 3. WebSocket Hub 초기화 및 실행
	hub := internal.NewHub()
	go hub.Run()

	// 4. Gin 라우터 설정
	router := gin.Default()

	// 5. CORS 설정
	// 프론트엔드 도메인과 로컬 개발 주소를 허용합니다.
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"https://chat.jungyu.xyz", // 실제 서비스 도메인
		"http://localhost:5173",   // 로컬 개발 환경
		"http://127.0.0.1:5173",
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}

	router.Use(cors.New(config))

	// 6. [중요] GKE Ingress 헬스 체크용 루트 경로 핸들러
	// 로드밸런서가 서비스 상태를 확인할 때 이 경로로 요청을 보냅니다.
	// 200 OK를 반환해야 502 Bad Gateway 에러가 발생하지 않습니다.
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Health Check OK")
	})

	// 7. API 라우트 설정
	api := router.Group("/api")
	{
		api.POST("/session", internal.CreateSessionHandler)
	}

	// 8. WebSocket 라우트 설정
	router.GET("/ws", func(c *gin.Context) {
		internal.ServeWs(hub, c)
	})

	// 9. 서버 실행
	router.Run(":8080")
}
