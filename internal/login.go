package internal

import (
	"database/sql"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if Errorfile {
		Error500(w, r)
		return
	}
	
	if er := ExeTemp(w, "login.html", nil); er != nil {
		Error500(w, r)
	}
}

func LoginAuth(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		Error405(w, r)
		return
	}

	

	nameInput := r.FormValue("username")
	paswdInput := r.FormValue("password")

	if nameInput == "" || paswdInput == "" {
		Error400(w, r)
		return
	}

	paswdhash, err1 := getPaswd(nameInput)

	if err1 != sql.ErrNoRows && err1 != nil { // assure us that there are no errors in the database
		Error500(w, r)
		return
	}

	if err1 == sql.ErrNoRows { // Bad nameInput
		ExeTemp(w, "login.html", data{ErrorAlert: "BadLog"})
		return
	}
	// check password after that nameInput is valid
	err2 := bcrypt.CompareHashAndPassword([]byte(paswdhash), []byte(paswdInput))

	if err2 != nil {
		ExeTemp(w, "login.html", data{ErrorAlert: "BadLog"})
		return
	}
	user, err := getUserFromDB(nameInput)
	if err != nil {
		Error500(w, r)
		return
	}
	//--------------------------user is connected------------------------------------------------------//
	c, _ := r.Cookie("session_token")
	
	if c != nil {
		SessionID := c.Value
		if IsSessionExit(user, SessionID) {
			http.Redirect(w, r, "/", code303)
			return
		}

		DeleteSession(SessionID)
	}
	
	// ExeTemp(w, "login.html", data{ErrorAlert: "SessExist2"})
	
	DeleteSession2(user.ID)
	
	// create a session token
	err3 := SessionStart(user, w) // start the session
	if err3 != nil {
		Error500(w, r)
		return
	}

	http.Redirect(w, r, "/", code303)

}

func getPaswd(nameInput string) (string, error) {
	row := DB.QueryRow("SELECT password FROM users WHERE username = ?", nameInput)
	var paswdhash string
	err := row.Scan(&paswdhash)

	if err != nil {
		return "", err
	}

	return paswdhash, err //sqlError && nil
}
