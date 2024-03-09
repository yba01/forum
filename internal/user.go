package internal

import (
	"database/sql"
	"fmt"
	"net/http"
	"golang.org/x/crypto/bcrypt"
)

// For Register
func CreateUser(user *User) error {
	_, err := DB.Exec("INSERT INTO users (username,email,password) VALUES(?,?,?)", user.Username, user.Email, user.HashPassword)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	return nil
}

// The user already exists in DB and we need his information
func getUserFromDB(name string) (*User, error) {

	newUser := User{}

	row := DB.QueryRow("SELECT * FROM users WHERE username = ?", name)

	err := row.Scan(&newUser.ID, &newUser.Username, &newUser.Password, &newUser.Email)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &newUser, err
}

// get data input when user registers
func getInputReg(r *http.Request) (*User, error) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	hashPasswd, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	user := &User{
		Username:     username,
		Email:        email,
		Password:     password,
		HashPassword: string(hashPasswd),
	}

	return user, nil
}
