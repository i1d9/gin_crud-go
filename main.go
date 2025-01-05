package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/i1d9/gin_crud-go/authenication"
	"github.com/i1d9/gin_crud-go/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"

	"log"
	"os"
)

func loadConfig() {
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}

func createDatabaseConnectionPool() (*pgxpool.Pool, error) {
	// Create a new connection dbpool
	database_url := viper.GetString("DATABASE_URL")
	dbpool, err := pgxpool.New(context.Background(), database_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection dbpool: %v\n", err)
		os.Exit(1)
	}

	return dbpool, err

}

func setupDatabaseTables(pool *pgxpool.Pool) {

	ctx := context.Background()
	pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS users (id bigserial primary key,first_name text not null, last_name text, surname text, email text not null unique,username text not null unique, mobile_number text not null unique, password text not null, inserted_at timestamp, updated_at timestamp)`)
	pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS sessions (id bigserial primary key, user_id int references(users) not null, token text not null unique, status text, expires_at timestamp not null, type text not null, inserted_at timestamp, updated_at timestamp)`)
	pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS posts (id bigserial primary key, user_id int references(users) not null, title text, body text, inserted_at timestamp, updated_at timestamp)`)

	pool.Exec(ctx, `CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email)`)
	pool.Exec(ctx, `CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users(username)`)

	pool.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions(token)`)
	pool.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id)`)
	pool.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_sessions_status ON sessions(status)`)

}

func main() {

	loadConfig()
	dbpool, _ := createDatabaseConnectionPool()
	setupDatabaseTables(dbpool)
	defer dbpool.Close()

	router := gin.Default()

	// Authenctication routes
	auth := router.Group("/auth")
	{
		auth.POST("/login", func(c *gin.Context) {
			authenication.Login(c, dbpool)
		})

		auth.POST("/register", func(c *gin.Context) {
			authenication.Register(c, dbpool)
		})

		auth.POST("/logout", func(c *gin.Context) {
			authenication.Logout(c, dbpool)
		})

	}

	// User Routes
	user := router.Group("/users")
	user.Use(middleware.VerifyAccessToken(dbpool))
	{
		user.GET("/search", func(c *gin.Context) {
			authenication.Login(c, dbpool)
		})

		user.GET("/profile", func(c *gin.Context) {
			authenication.Login(c, dbpool)
		})

	}

	// Post Routes
	post := router.Group("/posts")
	post.Use(middleware.VerifyAccessToken(dbpool))
	{
		post.POST("/create", func(c *gin.Context) {
			authenication.Login(c, dbpool)
		})
	}

	router.Run(":8080")
}
