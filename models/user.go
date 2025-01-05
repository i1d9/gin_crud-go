package models

import (
	"time"
	"errors"
)

type User struct {
	ID        int    `json:"id"`
	First_Name string `json:"first_name"`
	Last_Name  string `json:"last_name"`
	Surname   string `json:"surname"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Mobile_Number string `json:"mobile_number"`
	Created_At *time.Time `json:"created_at"`
	Updated_At *time.Time `json:"updated_at"`
}



func GetUserbyID(id int) (user User, error) {

	return user, errors.New("Not Found")
}

func CreateUser()  {
	
}


func UpdateUser()  {
	
}


func DeleteUser()  {
	
}


func SearchUser()  {
	
}