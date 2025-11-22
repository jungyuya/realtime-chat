package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors" // CORS 미들웨어 임포트 추가
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/jungyuya/realtime-chat/backend/internal"
	"github.com/jungyuya/realtime-chat/backend/internal/db"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, reading from environment variables")
	}

	// [추가] 데이터베이스 초기화 (연결 및 테이블 생성)
	// 주의: 로컬 개발 시 .env에 DB 정보가 없으면 여기서 에러가 나고 서버가 꺼질 수 있습니다.
	// DB 연결 정보가 있을 때만 실행하도록 조건문을 걸거나,
	// 로컬 테스트를 위해 .env에 DB 정보를 채워야 합니다. (Step 4.2에서 진행 예정)
	// 일단은 코드를 추가해둡니다.
	if os.Getenv("DB_HOST") != "" {
		db.InitDB()
	} else {
		log.Println("DB_HOST not set, skipping database initialization")
	}

	hub := internal.NewHub()
	go hub.Run()

	router := gin.Default()

	// CORS 미들웨어 설정을 더 유연하게 변경합니다.
	config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"http://localhost:5173"} // 특정 Origin만 허용하는 대신,
	config.AllowAllOrigins = true // 모든 Origin을 허용합니다. (개발 환경에서 유용)
	// 추가적으로 허용할 헤더를 명시할 수 있습니다.
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	// OPTIONS 메소드를 명시적으로 허용합니다.
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}

	router.Use(cors.New(config))

	// ... 나머지 코드 (api group, ws route, router.Run) ...
	api := router.Group("/api")
	{
		api.POST("/session", internal.CreateSessionHandler)
	}

	router.GET("/ws", func(c *gin.Context) {
		internal.ServeWs(hub, c)
	})

	router.Run(":8080")
}
