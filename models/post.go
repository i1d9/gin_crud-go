package models

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Post struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Body        string     `json:"body"`
	User_ID     string     `json:"user_id"`
	Inserted_At *time.Time `json:"inserted_at"`
	Updated_At  *time.Time `json:"updated_at"`
}

func Get(c *gin.Context, pool *pgxpool.Pool) {

	search_term, exists := c.GetQuery("term")

	if exists {

		users, err := SearchPost(pool, search_term)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, users)
	} else {

		users, err := GetPosts(pool)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, users)

	}

}

func GetPosts(pool *pgxpool.Pool) ([]Post, error) {

	ctx := context.Background()

	query := `SELECT id, title, body, user_id, inserted_at, updated_at FROM posts`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post

		err := rows.Scan(&post.ID, &post.Title, &post.Body, &post.User_ID, &post.Inserted_At, &post.Updated_At)

		if err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}
		posts = append(posts, post)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration error: %w", rows.Err())
	}

	return posts, nil
}

func CreatePost(pool *pgxpool.Pool, title string, body string, user_id int) error {

	ctx := context.Background()

	query := `INSERT INTO posts (title, body, user_id,  inserted_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := pool.Exec(ctx, query, title, body, user_id, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("create post: %v", err)
	}

	return nil
}

func Create(c *gin.Context, pool *pgxpool.Pool) {

	title := c.PostForm("title")
	body := c.PostForm("body")

	session, exists := c.Get("session")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	castedSession, _ := session.(Session)
	err := CreatePost(pool, title, body, castedSession.User_ID)

	if err != nil {
		c.JSON(400, gin.H{
			"error":   err.Error(),
			"message": "Invalid parameter were provided",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Post created successfully",
	})

}

func UpdatePost(pool *pgxpool.Pool, post_id int, title string, body string) (int, error) {

	ctx := context.Background()
	rows := 0
	query := "UPDATE posts SET title = $1, body = $2, updated_at = $3 WHERE ID = $4"

	res, err := pool.Exec(ctx, query, title, body, time.Now(), post_id)

	if err != nil {
		return rows, fmt.Errorf("update user: %v", err)
	}
	rows = int(res.RowsAffected())

	return rows, nil
}

func Update(c *gin.Context, pool *pgxpool.Pool) {

	post_id, exists := c.GetQuery("id")

	if exists {
		title := c.PostForm("title")
		body := c.PostForm("body")

		post_id, err := strconv.Atoi(post_id)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   "Invalid post ID",
				"message": "Post ID must be an integer",
			})
			return
		}

		_, err = UpdatePost(pool, post_id, title, body)

		if err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"message": "Invalid parameter were provided",
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "Post update successfully",
		})
	} else {

		c.JSON(400, gin.H{
			"error":   "",
			"message": "Invalid parameter were provided",
		})
		return
	}

}

func DeletePost(pool *pgxpool.Pool, post_id int) (int, error) {

	ctx := context.Background()
	rows := 0
	query := "DELETE FROM posts WHERE ID = $1"

	res, err := pool.Exec(ctx, query, post_id)

	if err != nil {
		return rows, fmt.Errorf("delete post: %v", err)
	}
	rows = int(res.RowsAffected())

	return rows, nil

}

func Delete(c *gin.Context, pool *pgxpool.Pool) {

	post_id, exists := c.GetQuery("id")

	if exists {

		post_id, err := strconv.Atoi(post_id)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   "Invalid post ID",
				"message": "Post ID must be an integer",
			})
			return
		}

		_, err = DeletePost(pool, post_id)

		if err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"message": "Invalid parameter were provided",
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "Post deleted successfully",
		})
	} else {

		c.JSON(400, gin.H{
			"error":   "",
			"message": "Invalid parameter were provided",
		})
		return
	}

}

func SearchPost(pool *pgxpool.Pool, term string) ([]Post, error) {

	ctx := context.Background()

	query := `
		SELECT 
			id, 
			title, 
			body, 
			user_id,
			inserted_at, 
			updated_at 
		
		FROM posts
		WHERE
			title ILIKE '%' || $1 || '%'
	`

	rows, err := pool.Query(ctx, query, term)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post

		err := rows.Scan(&post.ID, &post.Title, &post.Body, &post.User_ID, &post.Inserted_At, &post.Updated_At)

		if err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}
		posts = append(posts, post)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration error: %w", rows.Err())
	}

	return posts, nil

}
