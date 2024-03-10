package internal

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
)

func ForumInfo(w http.ResponseWriter, r *http.Request) {
	post_id := r.URL.Query().Get("id")
	id, erid := strconv.Atoi(post_id)
	if erid != nil {
		Error500(w, r)
		return
	}

	var Post ppost
	Post.Post_id = post_id
	err := DB.QueryRow("SELECT title, message, user_id FROM post WHERE id = ?", post_id).Scan(&Post.Title, &Post.Message, &Post.User_id)
	err1 := DB.QueryRow("SELECT username FROM users WHERE id = ?", Post.User_id).Scan(&Post.User)
	if err != nil {
		Error400(w, r)
		return
	}
	if err1 != nil {
		Error400(w, r)
		return
	}
	
	total, errsum := DB.Query("SELECT SUM(liked), SUM(disliked) FROM reaction WHERE post_id = ?", post_id)
	if errsum != nil {
		Error500(w, r)
		return
	}
	for total.Next() {
		total.Scan(&Post.Likes, &Post.Dislikes)
	}

	var Commentary comment

	Comments, err := DB.Query("SELECT u.username, c.commentary, c.id FROM comments c JOIN users u ON c.user_id = u.id WHERE c.post_id = ?", id)
	if err != nil {
		Error500(w, r)
		return
	}
	for Comments.Next() {
		Comments.Scan(&Commentary.Username, &Commentary.Comment, &Commentary.Com_id)
		Commentary.Dislike = 0
		Commentary.Like = 0
		total, errsum := DB.Query("SELECT SUM(liked), SUM(disliked) FROM reactioncom WHERE com_id = ?", Commentary.Com_id)
		if errsum != nil {
			Error500(w, r)
			return
		}
		for total.Next() {
			total.Scan(&Commentary.Like, &Commentary.Dislike)
		}
		Post.Comments = append(Post.Comments, Commentary)
	}
	c, Sessioned, err := Authenticated(w, r, "forumInfo.html", Post)
	if err != nil || Sessioned == nil{
		return
	}
	if c != nil {
		sessionID := c.Value
		Post.IsConnect = true
		Post.Online_user = Sessioned.Username 
		e := RefreshSession(sessionID)
		if e != nil {
			Error500(w, r)
			return
		}
	}

	if er := ExeTemp(w, "forumInfo.html", Post); er != nil {
		Error500(w, r)
		return
	}
}

func ComentaryInsert(w http.ResponseWriter, r *http.Request) {


	post_id := r.URL.Query().Get("id")

	_, Session, err := Authenticated(w, r, "login.html", nil)
	if err != nil || Session == nil{
		return
	}
	var userID int
	if Session != nil{
		userID = Session.User_ID
	}
	

	commentary := r.FormValue("commentary")

	if commentary == "" || onlySpace(commentary){
		Error400(w, r)
		return
	}

	commentary = strings.Trim(commentary," ")

	_, err = DB.Exec("INSERT INTO comments (user_id, post_id, commentary) VALUES(?, ?, ?)", userID, post_id, commentary)
	if err != nil {
		Error500(w, r)
		return
	}

	http.Redirect(w, r, "forumInfo?id="+post_id, code303)

}

func ReactionHandle(w http.ResponseWriter, r *http.Request) {


	post_id := r.URL.Query().Get("id")
	_, Session, err := Authenticated(w, r, "login.html", nil)
	if err != nil || Session == nil {
		return
	}
	sessionID := Session.User_ID
	action := r.FormValue("action")
	liked, disliked := 0, 0
	switch action {
	case "like":
		liked, disliked = 1, 0
	case "dislike":
		liked, disliked = 0, 1
	}

	rows := DB.QueryRow("SELECT user_id, post_id, liked, disliked FROM reaction WHERE user_id=? AND post_id=?", sessionID, post_id)

	var id, id2, like, dislike int

	err = rows.Scan(&id, &id2, &like, &dislike)
	if like == liked && dislike == disliked {
		liked = 0
		disliked = 0
	}
	if err != nil {
		if err == sql.ErrNoRows {
			_, err1 := DB.Exec("INSERT INTO reaction (user_id, post_id, liked, disliked) VALUES(?, ?, ?, ?)", sessionID, post_id, liked, disliked)
			if err1 != nil {
				Error500(w, r)
				return
			}
		}
	} else {
		_, err1 := DB.Exec("UPDATE reaction SET liked = ?, disliked = ? WHERE post_id=? AND user_id=?", liked, disliked, post_id, sessionID)
		if err1 != nil {
			Error400(w, r)
			return
		}
	}

	http.Redirect(w, r, "forumInfo?id="+post_id, code303)

}

func ReactioncomHandle(w http.ResponseWriter, r *http.Request) {

	
	com_id := r.URL.Query().Get("id")
	_, Session, err := Authenticated(w, r, "login.html", nil)
	if err != nil || Session == nil {
		return
	}
	sessionID := Session.User_ID
	action := r.FormValue("actioncom")
	liked, disliked := 0, 0
	switch action {
	case "like":
		liked, disliked = 1, 0
	case "dislike":
		liked, disliked = 0, 1
	}

	rows := DB.QueryRow("SELECT user_id, com_id, liked, disliked FROM reactioncom WHERE user_id=? AND com_id=?", sessionID, com_id)
	var id int
	var id2 int
	var like int
	var dislike int
	err = rows.Scan(&id, &id2, &like, &dislike)
	if like == liked && dislike == disliked {
		liked = 0
		disliked = 0
	}
	if err != nil {
		if err == sql.ErrNoRows {
			_, err1 := DB.Exec("INSERT INTO reactioncom (user_id, com_id, liked, disliked) VALUES(?, ?, ?, ?)", sessionID, com_id, liked, disliked)
			if err1 != nil {
				Error500(w, r)
				return
			}
		}
	} else {
		_, err1 := DB.Exec("UPDATE reactioncom SET liked = ?, disliked = ? WHERE com_id=? AND user_id=?", liked, disliked, com_id, sessionID)
		if err1 != nil {
			Error500(w, r)
			return
		}
	}
	var post int
	err = DB.QueryRow("SELECT c.post_id FROM reactioncom r JOIN comments c ON c.id = r.com_id WHERE r.com_id = ?", com_id).Scan(&post)
	if err != nil {
		Error500(w, r)
		return
	}
	post_id := strconv.Itoa(post)
	http.Redirect(w, r, "forumInfo?id="+post_id, code303)
}
