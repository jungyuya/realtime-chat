package internal

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Profile은 명사(동물 이름)와 해당 아바타 파일 키를 묶는 구조체입니다.
type Profile struct {
	Noun      string
	AvatarKey string
}

// 미리 정의된 형용사와 프로필 배열
var adjectives = []string{ "익명의", "용감한", "재빠른", "총명한", "친절한", "고요한", "유쾌한", "배고픈", "성실한"}
var profiles = []Profile{
	{Noun: "사자", AvatarKey: "avatar-1"},
	{Noun: "호랑이", AvatarKey: "avatar-2"},
	{Noun: "코알라", AvatarKey: "avatar-3"},
	{Noun: "곰탱이", AvatarKey: "avatar-4"},
	{Noun: "돌고래", AvatarKey: "avatar-5"},
	{Noun: "토끼", AvatarKey: "avatar-6"},
	{Noun: "고양이", AvatarKey: "avatar-7"},
	{Noun: "여우", AvatarKey: "avatar-8"},
	{Noun: "판다", AvatarKey: "avatar-9"},
	{Noun: "강아지", AvatarKey: "avatar-10"},
	{Noun: "올빼미", AvatarKey: "avatar-11"},
	{Noun: "고슴도치", AvatarKey: "avatar-12"},
	{Noun: "사슴", AvatarKey: "avatar-13"},
	{Noun: "기린", AvatarKey: "avatar-14"},
	{Noun: "수달", AvatarKey: "avatar-15"},
	{Noun: "다람쥐", AvatarKey: "avatar-16"},

	// 추가 아바타와 이미지에 맞춰 계속 추가할 예정.
}

// SessionRequest는 닉네임을 받지 않습니다.
// SessionRequest 구조체 수정
type SessionRequest struct {
	AnonymousID string `json:"anonymousId" binding:"required"`
}

type SessionResponse struct {
	Token string `json:"token"`
}

// CreateSessionHandler 함수 수정
func CreateSessionHandler(c *gin.Context) {
	var req SessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hasher := sha1.New()
	hasher.Write([]byte(req.AnonymousID))
	hashBytes := hasher.Sum(nil)

	adjIndex := int(hashBytes[0]) % len(adjectives)
	profileIndex := int(hashBytes[1]) % len(profiles)

	selectedProfile := profiles[profileIndex]
	nickname := fmt.Sprintf("%s %s", adjectives[adjIndex], selectedProfile.Noun)
	avatar := selectedProfile.AvatarKey

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Secret key not found"})
		return
	}

	claims := jwt.MapClaims{
		"anonymousId": req.AnonymousID,
		"nickname":    nickname,
		"avatar":      avatar,
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
		"iat":         time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, SessionResponse{Token: tokenString})
}
