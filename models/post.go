package models

import "time"

type Post struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	User_ID string `json:"user_id"`
	Inserted_At *time.Time `json:"inserted_at"`
	Updated_At *time.Time `json:"updated_at"`
}



func CreatePost()  {
	
}


func UpdatePost()  {
	
}


func DeletePost()  {
	
}


func SearchPost()  {
	
}