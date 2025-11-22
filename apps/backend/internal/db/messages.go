package db

import (
	"context"
	"time"
)

// DB에 저장/조회할 메시지 구조체
type DBMessage struct {
	ID             int64     `json:"id"`
	RoomID         string    `json:"roomId"`
	SenderID       string    `json:"senderId"`
	SenderNickname string    `json:"senderNickname"`
	Avatar         string    `json:"avatar"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"createdAt"`
}

// SaveMessage: 메시지를 DB에 저장합니다.
func SaveMessage(roomID, senderID, nickname, avatar, content string) error {
	query := `
        INSERT INTO messages (room_id, sender_id, sender_nickname, avatar, content, created_at)
        VALUES ($1, $2, $3, $4, $5, NOW())
    `
	// Pool 변수는 database.go에서 정의된 것을 사용합니다 (같은 패키지이므로 가능)
	_, err := Pool.Exec(context.Background(), query, roomID, senderID, nickname, avatar, content)
	return err
}

// GetRecentMessages: 특정 방의 최근 메시지 N개를 불러옵니다.
func GetRecentMessages(roomID string, limit int) ([]DBMessage, error) {
	query := `
        SELECT id, room_id, sender_id, sender_nickname, avatar, content, created_at
        FROM messages
        WHERE room_id = $1
        ORDER BY created_at DESC
        LIMIT $2
    `
	rows, err := Pool.Query(context.Background(), query, roomID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []DBMessage
	for rows.Next() {
		var msg DBMessage
		err := rows.Scan(&msg.ID, &msg.RoomID, &msg.SenderID, &msg.SenderNickname, &msg.Avatar, &msg.Content, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	// 역순 정렬 (과거 -> 최신)
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}