package internal

import "net/http"

func Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		Error500(w, r)
		return
	}

	sessionID := c.Value
	//set a new cookie for same name and value in order to delete it with maxAge = -1
	DelCookie(w)

	if e := DeleteSession(sessionID); e != nil{
		Error500(w,r)
		return
	}

	http.Redirect(w,r,"/",code303)
}
