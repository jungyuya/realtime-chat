package internal

import (
	"net/http"
	"os"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// SessionRequest는 세션 생성 요청의 본문(body) 구조를 정의합니다.
// `json:"..."` 태그는 JSON 데이터를 Go 구조체 필드에 매핑하는 방법을 알려줍니다.
type SessionRequest struct {
	AnonymousID string `json:"anonymousId" binding:"required"`
	Nickname    string `json:"nickname" binding:"required,min=2,max=15"`
}

// SessionResponse는 세션 생성 응답의 구조를 정의합니다.
type SessionResponse struct {
	Token string `json:"token"`
}

// CreateSessionHandler는 POST /api/session 요청을 처리하는 핸들러입니다.
func CreateSessionHandler(c *gin.Context) {
	var req SessionRequest
	// c.ShouldBindJSON은 요청 본문의 JSON을 req 구조체에 바인딩(매핑)합니다.
	// 만약 nickname의 길이가 2~15자가 아니거나, 필드가 누락되면 자동으로 에러를 반환합니다.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 환경 변수에서 비밀 키를 가져옵니다.
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Secret key not found"})
		return
	}

	// JWT 클레임(Claim)을 설정합니다. 클레임은 토큰에 담을 정보 조각들입니다.
	claims := jwt.MapClaims{
		"anonymousId": req.AnonymousID,
		"nickname":    req.Nickname,
		// 'exp' (Expiration Time)는 토큰의 만료 시간을 나타냅니다.
		// 지금으로부터 24시간 뒤에 만료되도록 설정합니다.
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		// 'iat' (Issued At)는 토큰이 발급된 시간을 나타냅니다.
		"iat": time.Now().Unix(),
	}

	// HMAC-SHA256 알고리즘으로 서명할 새 토큰을 생성합니다.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 비밀 키를 사용하여 토큰에 서명하고, 완전한 토큰 문자열을 얻습니다.
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// 성공적으로 생성된 토큰을 클라이언트에게 응답으로 보냅니다.
	c.JSON(http.StatusOK, SessionResponse{Token: tokenString})
}