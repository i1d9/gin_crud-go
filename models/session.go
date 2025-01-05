package models

import (
	"time"
	"errors"
)

type Session struct {
	ID         int        `json:"id"`
	User_ID    string     `json:"user_id"`
	Token      string     `json:"token"`
	Expires_At *time.Time `json:"expires_at"`
	Type       string     `json:"type"`
	Created_At *time.Time `json:"created_at"`
	Updated_At *time.Time `json:"updated_at"`
}


func FindSessionbyToken(token String) (session Session, error) {

	return session, errors.New("Not Found")
}

func CreateSession()  {
	
}


func UpdateSession()  {
	
}


func DeleteSession()  {
	
}


func SearchSession()  {
	
}