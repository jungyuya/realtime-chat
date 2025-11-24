package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB 연결을 관리할 전역 변수 (Connection Pool)
var Pool *pgxpool.Pool

// InitDB: 데이터베이스 연결을 초기화하고 스키마를 생성합니다.
func InitDB() {
	var err error
	// 환경 변수에서 DB 접속 정보를 가져옵니다.
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// [디버깅 로그 추가] 비밀번호는 제외하고 연결 정보 출력
	log.Printf("DEBUG: Connecting to DB at %s:%s, DB: %s, User: %s",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"))

	// Connection Pool 설정을 생성합니다.
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Unable to parse DB config: %v", err)
	}

	// 연결 풀을 생성합니다. (백그라운드에서 연결 관리)
	Pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	// 연결 테스트 (Ping)
	if err := Pool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	log.Println("Connected to PostgreSQL successfully!")

	// 테이블 생성 (Migration)
	createTables()
}

// createTables: 필요한 테이블이 없으면 자동으로 생성합니다.
func createTables() {
	query := `
    CREATE TABLE IF NOT EXISTS messages (
        id BIGSERIAL PRIMARY KEY,
        room_id VARCHAR(255) NOT NULL,
        sender_id VARCHAR(255) NOT NULL,
        sender_nickname VARCHAR(50) NOT NULL,
        avatar VARCHAR(50) NOT NULL,
        content TEXT NOT NULL,
        created_at TIMESTAMPTZ DEFAULT NOW()
    );
    
    -- room_id와 시간순 정렬을 위한 인덱스 생성 (조회 성능 최적화)
    CREATE INDEX IF NOT EXISTS idx_messages_room_created ON messages (room_id, created_at DESC);
    `

	_, err := Pool.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
	log.Println("Database schema initialized.")
}
