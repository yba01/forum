package internal

import (
	"database/sql"
	"net/http"
	"regexp"
	"strings"
	"unicode"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if Errorfile {
		Error500(w, r)
		return
	}


	if er := ExeTemp(w, "register.html", nil); er != nil {
		Error500(w, r)
	}

}

func RegisterAuth(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		Error405(w, r)
		return
	}

	newUser, err := getInputReg(r)
	if err != nil {
		Error500(w, r)
		return
	}

	newUser.Username = strings.ToLower(newUser.Username)
	newUser.Email = strings.ToLower(newUser.Email)

	//check if username, email and password are empty
	if newUser.Username == "" || newUser.Email == "" || newUser.Password == "" {
		Error400(w, r)
		return
	}

	username := newUser.Username
	email := newUser.Email
	password := newUser.Password

	if !CheckName(username) || !CheckPassword(password) {
		ExeTemp(w, "register.html", data{ErrorAlert: "BadCri"})
		return
	}
	if !CheckEmail(email) {
		ExeTemp(w, "register.html", data{ErrorAlert: "email"})
		return
	}

	er := IsNameOrEmailExist(username, email)

	if er != sql.ErrNoRows && er != nil { // assure us that there are no errors in the database
		Error500(w, r)
		return
	}
	if er == nil {
		ExeTemp(w, "register.html", data{ErrorAlert: "UnReg"})
		return
	}

	
	if err := CreateUser(newUser); err != nil {
		Error500(w, r)
		return
	}

	
	ExeTemp(w, "login.html", data{ErrorAlert: "ok"})


}

// check username for only alphaNumeric characters
func CheckName(name string) bool {
	var (
		IsAlphanumeric = true
		lengthName     = false
	)
	if 5 <= len(name) && len(name) <= 50 {
		lengthName = true
	}
	for _, char := range name {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			IsAlphanumeric = false
		}
	}

	return IsAlphanumeric && lengthName
}

// check password criteria
func CheckPassword(pwd string) bool {
	var (
		paswdLowercase, paswdUppercase, paswdNumber, paswdLength, paswdSpecial bool
		paswdNoSpace                                                           = true
	)
	for _, char := range pwd {
		switch {
		case unicode.IsLower(char):
			paswdLowercase = true
		case unicode.IsUpper(char):
			paswdUppercase = true
		case unicode.IsNumber(char):
			paswdNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			paswdSpecial = true
		case unicode.IsSpace(char):
			paswdNoSpace = false
		}
	}
	if 9 < len(pwd) && len(pwd) < 60 {
		paswdLength = true
	}
	if !paswdLowercase || !paswdUppercase || !paswdNumber || !paswdLength || !paswdSpecial || !paswdNoSpace {
		return false
	}
	return true
}

// Check if username or email is already used
func IsNameOrEmailExist(name, email string) error {
	row1 := DB.QueryRow("SELECT * FROM users WHERE username = ? ", name)
	row2 := DB.QueryRow("SELECT * FROM users WHERE email = ? ", email)
	user := User{}
	err1 := row1.Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	err2 := row2.Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err1 != nil && err2 != nil {
		return err1
	}

	return nil
}

// Check if email is correct
func CheckEmail(email string) bool {
	// Regular expression pattern for basic email
	pattern := "^[a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+.[a-zA-Z.]{2,15}$"

	// Compile the pattern
	regex := regexp.MustCompile(pattern)

	// Check if the email match pattern
	return regex.MatchString(email)
}
