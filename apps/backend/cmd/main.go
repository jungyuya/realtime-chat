package main

import (
	"log"

	"github.com/gin-contrib/cors" // CORS 미들웨어 임포트 추가
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/jungyuya/realtime-chat/backend/internal"
)

func main() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file")
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