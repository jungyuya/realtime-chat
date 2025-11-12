package internal

import (
	"crypto/sha1"
	//"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 미리 정의된 형용사와 명사 배열
var adjectives = []string{"익명의", "용감한", "재빠른", "총명한", "친절한", "고요한", "빛나는", "현명한"}
var nouns = []string{"사자", "호랑이", "코끼리", "기린", "돌고래", "독수리", "거북이", "고래"}

// SessionRequest는 세션 생성 요청의 본문(body) 구조를 정의합니다.
// 이제 클라이언트로부터 닉네임을 받지 않습니다.
type SessionRequest struct {
	AnonymousID string `json:"anonymousId" binding:"required"`
}

// SessionResponse는 세션 생성 응답의 구조를 정의합니다.
type SessionResponse struct {
	Token string `json:"token"`
}

// CreateSessionHandler는 POST /api/session 요청을 처리하는 핸들러입니다.
func CreateSessionHandler(c *gin.Context) {
	var req SessionRequest
	// c.ShouldBindJSON은 요청 본문의 JSON을 req 구조체에 바인딩(매핑)합니다.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// anonymousId를 기반으로 "결정론적 랜덤" 닉네임을 생성합니다.
	// 즉, 동일한 ID는 항상 동일한 닉네임을 갖게 됩니다.
	hasher := sha1.New()
	hasher.Write([]byte(req.AnonymousID))
	hashBytes := hasher.Sum(nil)

	// 해시 값의 일부를 사용하여 배열의 인덱스를 결정합니다.
	adjIndex := int(hashBytes[0]) % len(adjectives)
	nounIndex := int(hashBytes[1]) % len(nouns)
	nickname := fmt.Sprintf("%s %s", adjectives[adjIndex], nouns[nounIndex])

	// 환경 변수에서 비밀 키를 가져옵니다.
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Secret key not found"})
		return
	}

	// JWT 클레임(Claim)을 설정합니다. 클레임은 토큰에 담을 정보 조각들입니다.
	claims := jwt.MapClaims{
		"anonymousId": req.AnonymousID,
		"nickname":    nickname, // 직접 입력받는 대신, 생성된 닉네임을 사용
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