package models

import (
	"time"

	"context"
	"fmt"
	"math/rand"

	"github.com/jackc/pgx/v5/pgxpool"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

type Session struct {
	ID          int        `json:"id"`
	User_ID     int        `json:"user_id"`
	Token       string     `json:"token"`
	Status      string     `json:"status"`
	Expires_At  *time.Time `json:"expires_at"`
	Type        string     `json:"type"`
	Inserted_At *time.Time `json:"inserted_at"`
	Updated_At  *time.Time `json:"updated_at"`
}

func FindSessionbyToken(pool *pgxpool.Pool, token string) (Session, error) {

	ctx := context.Background()

	var session Session

	query := "SELECT id, user_id, token, expires_at, type, status, inserted_at, updated_at FROM sessions WHERE status = 'active' AND token = $1"
	if err := pool.QueryRow(ctx, query, token).Scan(
		&session.ID, &session.User_ID, &session.Token, &session.Status, &session.Expires_At, &session.Type, &session.Inserted_At, &session.Updated_At,
	); err != nil {
		return session, fmt.Errorf("Session by Token:  %v", err)
	}

	return session, nil

}

func CreateSession(pool *pgxpool.Pool, user_id int) (int, error) {

	id := 0
	ctx := context.Background()

	token := GenerateToken(32)

	query := `INSERT INTO sessions (user_id, token, status, type, expires_at) VALUES ($1, $2, $3, $4, $5)`

	err := pool.QueryRow(ctx, query, user_id, token, "active", "access_token", time.Now().Add(time.Hour*1))
	if err != nil {
		return id, fmt.Errorf("create session: %v", err)
	}

	return id, nil

}

func UpdateSession() {

}

func DeleteSession(pool *pgxpool.Pool, id int) (int, error) {

	ctx := context.Background()
	rows := 0
	query := "DELETE FROM sessions WHERE ID = $1"

	res, err := pool.Exec(ctx, query, id)

	if err != nil {
		return rows, fmt.Errorf("delete session: %v", err)
	}
	rows = int(res.RowsAffected())

	return rows, nil

}

func SearchSession() {

}

func GetSessionByID(pool *pgxpool.Pool, id int) (Session, error) {

	ctx := context.Background()
	var session Session

	query := "SELECT id, user_id, token, status, expires_at, type, inserted_at, updated_at FROM sessions WHERE ID = $1"
	if err := pool.QueryRow(ctx, query, id).Scan(
		&session.ID, &session.User_ID, &session.Token, &session.Status, &session.Expires_At, &session.Type, &session.Inserted_At, &session.Updated_At,
	); err != nil {
		return session, fmt.Errorf("session by id:  %v", err)
	}

	return session, nil

}

func GetSessionByTokenAndTokenType(pool *pgxpool.Pool, token string, token_type string) (Session, error) {

	ctx := context.Background()
	var session Session

	query := "SELECT id, user_id, token, status, expires_at, type, inserted_at, updated_at FROM sessions WHERE token = $1 and type = $2"
	if err := pool.QueryRow(ctx, query, token, token_type).Scan(
		&session.ID, &session.User_ID, &session.Token, &session.Status, &session.Expires_At, &session.Type, &session.Inserted_At, &session.Updated_At,
	); err != nil {
		return session, fmt.Errorf("session by id:  %v", err)
	}

	return session, nil

}

func GenerateRandomString(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateToken(length int) string {
	return GenerateRandomString(length, charset)
}
