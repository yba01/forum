package internal

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func Authenticated(w http.ResponseWriter, r *http.Request, page string, data any) (*http.Cookie, *Session, error) {
	c, err := r.Cookie("session_token")

	if err != nil {
		e := ExeTemp(w, page, data)
		if e != nil {
			Error500(w, r)
			return nil, nil, err
		}
		return nil, nil, err
	}

	sessionID := c.Value
	UserInfo, er := CheckSession(sessionID)
	if er != sql.ErrNoRows && er != nil {
		Error500(w, r)
		return nil, nil, er
	}
	if er == sql.ErrNoRows {
		DelCookie(w)
		http.Redirect(w, r, "/", code303)
		return nil, nil, er
	}
	//// FROM HERE THE SESSION EXISTS, LET'S CHECK IF IT IS EXPIRED OR NOT BY RECOVERING
	// EXPIRED_AT AND COMPARE IT WITH TIME.NOW()
	exp, err := GetExpired_AT(sessionID)
	if err != nil {
		Error500(w, r)
		return nil, nil, err
	}
	expiry, err := ParseTime(exp)
	if err != nil {
		Error500(w, r)
		return nil, nil, err
	}
	if expiry.Before(time.Now().UTC()) {
		e := DeleteSession(sessionID)
		if e != nil {
			Error500(w, r)
			return nil, nil, e
		}
		DelCookie(w)
		http.Redirect(w, r, "/login", code303)
		return nil, nil, nil
	}
	return c, UserInfo, nil
}

func ParseTime(timeStr string) (*time.Time, error) {
	layout2 := "2006-01-02 15:04:05"

	parsedTime, err := time.Parse(layout2, timeStr)
	parsedTime = parsedTime.UTC()
	if err != nil {
		return nil, err
	}

	return &parsedTime, nil
}

func SessionStart(user *User, w http.ResponseWriter) error {
	SessionToken := uuid.NewString()
	ExpiresAT := time.Now().UTC().Add(time.Hour)
	expired_at := ExpiresAT.Format("2006-01-02 15:04:05")
	maxAge := 4400
	if saveErr := SaveSession(SessionToken, user.Username, expired_at, user.ID); saveErr != nil {
		return saveErr
	}
	cookie := http.Cookie{
		Name:    "session_token",
		Value:   SessionToken,
		Expires: ExpiresAT,
		MaxAge:  maxAge,
	}
	http.SetCookie(w, &cookie)
	return nil
}
func CheckSession(sessionID string) (*Session, error) {
	row := DB.QueryRow("SELECT * FROM sessions WHERE session_id = ?", sessionID)
	var sess Session
	err := row.Scan(&sess.Token, &sess.User_ID, &sess.Username, &sess.Create_at, &sess.Expired_at)
	if err != nil {
		return nil, err
	}

	return &sess, nil

}

func SaveSession(token, username, expired_AT string, user_id int) error {
	_, err := DB.Exec("INSERT INTO sessions (session_id,user_id,username,expired_at) VALUES(?,?,?,?)", token, user_id, username, expired_AT)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}

func DeleteSession(sessionID string) error {
	_, err := DB.Exec("DELETE  FROM sessions WHERE session_id = ?", sessionID)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}
func DeleteSession2(userID int) error {
	_, err := DB.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	if err != nil {
		return nil
	}
	return nil
}

func GetExpired_AT(SessionID string) (string, error) {
	var expired string
	row := DB.QueryRow("SELECT expired_at FROM sessions WHERE session_id = ?", SessionID)
	err := row.Scan(&expired)
	if err != nil {
		return "", err
	}
	return expired, nil
}

// update the session time each time the user navigates to the page
func RefreshSession(sessionID string) error {
	ExpiresAT := time.Now().UTC().Add(time.Hour)
	expired_at := ExpiresAT.Format("2006-01-02 15:04:05")
	_, err := DB.Exec("UPDATE sessions SET expired_at = ? WHERE session_id = ? ", expired_at, sessionID)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}

func IsSessionExit(user *User, sessID string) bool {
	var id int
	row := DB.QueryRow("SELECT user_id FROM sessions where session_id = ? ", sessID)
	err := row.Scan(&id)
	if err != nil {
		return false
	}
	return user.ID == id
}

// Del from navigates
func DelCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)
}
