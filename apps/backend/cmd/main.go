package main

import (
	"log"

	"github.com/gin-contrib/cors" // CORS 미들웨어 임포트 추가
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/jungyuya/realtime-chat/backend/internal"
)

func main() {
	// [수정] .env 파일 로드를 시도하되, 파일이 없어도 에러를 발생시키지 않습니다.
	// godotenv.Load()는 파일이 없을 경우 조용히 무시하고 넘어갑니다.
	// 이를 통해 로컬 개발 환경에서는 .env 파일을 사용하고,
	// 컨테이너 환경(파일이 없는)에서는 환경 변수를 직접 사용하게 됩니다.
	err := godotenv.Load() // 경로를 명시하지 않으면 현재 디렉토리에서 .env를 찾습니다.
	if err != nil {
		// .env 파일이 없는 것은 에러가 아니므로, 경고 메시지만 출력합니다.
		log.Println("Warning: .env file not found, reading from environment variables")
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