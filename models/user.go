package models

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matthewhartstonge/argon2"
	"time"
)

type User struct {
	ID            int        `json:"id"`
	Gender        string     `json:"gender"`
	First_Name    string     `json:"first_name"`
	Last_Name     string     `json:"last_name"`
	Surname       string     `json:"surname"`
	Username      string     `json:"username"`
	Email         string     `json:"email"`
	Password      string     `json:"password"`
	Mobile_Number string     `json:"mobile_number"`
	Inserted_At   *time.Time `json:"inserted_at"`
	Updated_At    *time.Time `json:"updated_at"`
}

func GetAuthUser(pool *pgxpool.Pool, identifier string) (User, error) {

	ctx := context.Background()

	var user User
	query := `SELECT id, first_name, last_name, surname, email, password, username, mobile_number, gender,inserted_at, updated_at FROM users WHERE email = $1 OR username = $1`

	if err := pool.QueryRow(ctx, query, identifier).Scan(
		&user.ID, &user.First_Name, &user.Last_Name, &user.Surname, &user.Email, &user.Password, &user.Username, &user.Mobile_Number, &user.Gender, &user.Inserted_At, &user.Updated_At,
	); err != nil {
		return user, fmt.Errorf("get user by identifier  %v", err)
	}

	return user, nil

}

func GetUsers(pool *pgxpool.Pool) ([]User, error) {

	ctx := context.Background()

	query := `SELECT id, first_name, last_name, surname, email, password, username, mobile_number, gender,inserted_at, updated_at FROM users`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User

		err := rows.Scan(&user.ID, &user.First_Name, &user.Last_Name, &user.Surname, &user.Email, &user.Password, &user.Username, &user.Mobile_Number, &user.Gender, &user.Inserted_At, &user.Updated_At)

		if err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}
		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration error: %w", rows.Err())
	}

	return users, nil
}

func GetUserbyID(pool *pgxpool.Pool, id int) (User, error) {

	ctx := context.Background()

	var user User

	query := "SELECT id, first_name, last_name, surname, email, username, mobile_number, password, gender,inserted_at, updated_at FROM users WHERE ID = $1"
	if err := pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.First_Name, &user.Last_Name, &user.Surname, &user.Email, &user.Username, &user.Mobile_Number, &user.Password, &user.Gender, &user.Inserted_At, &user.Updated_At,
	); err != nil {
		return user, fmt.Errorf("user by id:  %v", err)
	}

	return user, nil
}

func CreateUser(pool *pgxpool.Pool, first_name string, last_name string, surname string, email string, username string, mobile_number string, gender string, password string) error {

	argon := argon2.DefaultConfig()

	ctx := context.Background()

	encoded_password, encoded_password_err := argon.HashEncoded([]byte(password))
	if encoded_password_err != nil {
		return fmt.Errorf("encoded password error: %v", encoded_password_err)
	}

	query := `INSERT INTO users (first_name, last_name, surname, email, username, mobile_number,gender, password, inserted_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := pool.Exec(ctx, query, first_name, last_name, surname, email, username, mobile_number, gender, string(encoded_password), time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("create user: %v", err)
	}

	return nil

}

func UpdateUser(pool *pgxpool.Pool, user_id int, first_name string, last_name string, email string, mobile_number string) (int, error) {

	ctx := context.Background()
	rows := 0
	query := "UPDATE users SET first_name = $1, last_name = $2, email = $3, mobile_number = $4 WHERE ID = $5"

	res, err := pool.Exec(ctx, query, first_name, last_name, email, mobile_number, user_id)

	if err != nil {
		return rows, fmt.Errorf("update user: %v", err)
	}
	rows = int(res.RowsAffected())

	return rows, nil
}

func DeleteUser(pool *pgxpool.Pool, id int) (int, error) {

	ctx := context.Background()
	rows := 0
	query := "DELETE FROM users WHERE ID = $1"

	res, err := pool.Exec(ctx, query, id)

	if err != nil {
		return rows, fmt.Errorf("delete user: %v", err)
	}
	rows = int(res.RowsAffected())

	return rows, nil

}
func SearchUsers(pool *pgxpool.Pool, term string) ([]User, error) {

	ctx := context.Background()

	query := `
		SELECT 
			id, 
			first_name, 
			last_name, 
			surname, 
			email, 
			username, 
			mobile_number,
			password, 
			gender,
			inserted_at, 
			updated_at 
		
		FROM users
		WHERE
			username ILIKE '%' || $1 || '%' OR
			first_name ILIKE '%' || $1 || '%' OR
			last_name ILIKE '%' || $1 || '%' OR
			surname ILIKE '%' || $1 || '%' OR
			email ILIKE '%' || $1 || '%' OR
			mobile_number ILIKE '%' || $1 || '%'
	`

	rows, err := pool.Query(ctx, query, term)

	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User

		err := rows.Scan(&user.ID, &user.First_Name, &user.Last_Name, &user.Surname, &user.Email, &user.Password, &user.Username, &user.Mobile_Number, &user.Gender, &user.Inserted_At, &user.Updated_At)

		if err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}
		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration error: %w", rows.Err())
	}

	return users, nil
}
