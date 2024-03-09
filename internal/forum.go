package internal

import (
	"net/http"
	"strings"
)

func Forum(w http.ResponseWriter, r *http.Request) {
	if Errorfile {
		Error500(w, r)
		return
	}

	if ErrorUrl(r) {
		Error404(w, r)
		return

	}
	PostInf, err := DB.Query("SELECT p.id, p.title, u.username, p.category_id FROM post p JOIN users u ON p.user_id = u.id")

	if err != nil {
		Error500(w, r)
		return
	}

	Categories, err2 := DB.Query("SELECT category FROM categories")

	if err2 != nil {
		Error500(w, r)
		return
	}

	var category string

	var Post post
	var forumpage forumhome

	for Categories.Next() {
		Categories.Scan(&category)
		forumpage.Categories = append(forumpage.Categories, category)
	}
	
	for PostInf.Next() {
		var categories string
		PostInf.Scan(&Post.Id, &Post.Title, &Post.Username, &categories)
		Post.Category = strings.Fields(categories)
		forumpage.Allposts = append(forumpage.Allposts, Post)
	}
	if len(forumpage.Allposts) < 1 {
		forumpage.ErrorString = "There is nothing here..."
	}
	c, Sessioning, err := Authenticated(w, r, "forum.html", forumpage)
	if err != nil {
		return
	}
	if c != nil {
		sessionID := c.Value
		forumpage.IsConnect = true
		forumpage.Username = Sessioning.Username
		e := RefreshSession(sessionID)
		if e != nil {
			Error500(w, r)
			return
		}
	}

	if er := ExeTemp(w, "forum.html", forumpage); er != nil {
		Error500(w, r)
		return
	}

}

func Filterforum(w http.ResponseWriter, r *http.Request) {
	if ErrorUrl(r) {
		Error404(w, r)
		return

	}

	category_in := r.URL.Query().Get("id")


	if category_in == "createdpost" {
		CreatedPostFilter(w, r)
		return
	}

	if category_in == "LikedPost" {
		LikedPostFilter(w, r)
		return
	}

	PostInf, err := DB.Query("SELECT p.id, p.title, u.username, p.category_id FROM post p JOIN users u ON p.user_id = u.id WHERE p.category_id LIKE ?", "% "+category_in+" %")

	if err != nil {
		Error500(w, r)
		return
	}

	Categories, err2 := DB.Query("SELECT category FROM categories")

	if err2 != nil {
		Error500(w, r)
		return
	}

	var category string

	var Post post
	var forumpage forumhome

	for Categories.Next() {
		Categories.Scan(&category)
		forumpage.Categories = append(forumpage.Categories, category)
	}
	for PostInf.Next() {
		var categories string
		PostInf.Scan(&Post.Id, &Post.Title, &Post.Username, &categories)
		Post.Category = strings.Fields(categories)
		forumpage.Allposts = append(forumpage.Allposts, Post)
	}

	if len(forumpage.Allposts) < 1 {
		forumpage.ErrorString = "There is nothing here..."
	}

	c, sess, err := Authenticated(w, r, "forum.html", forumpage)
	if err != nil {
		return
	}
	if c != nil {
		sessionID := c.Value
		forumpage.IsConnect = true
		forumpage.Username = sess.Username
		e := RefreshSession(sessionID)
		if e != nil {
			Error500(w, r)
			return
		}
	}

	if er := ExeTemp(w, "forum.html", forumpage); er != nil {
		Error500(w, r)
		return
	}
}

func CreatedPostFilter(w http.ResponseWriter, r *http.Request) {
	if ErrorUrl(r) {
		Error404(w, r)
		return

	}

	c, session, err := Authenticated(w, r, "login.html", nil)
	if err != nil || session == nil {
		return
	}
	if c != nil {
		sessionID := c.Value
		e := RefreshSession(sessionID)
		if e != nil {
			Error500(w, r)
			return
		}
	}

	userId := session.User_ID

	PostInf, err := DB.Query("SELECT p.id, p.title, u.username, p.category_id FROM post p JOIN users u ON p.user_id = u.id WHERE u.id = ?", userId)

	if err != nil {
		Error500(w, r)
		return
	}

	Categories, err2 := DB.Query("SELECT category FROM categories")

	if err2 != nil {
		Error500(w, r)
		return
	}

	var category string

	var Post post
	var forumpage forumhome

	for Categories.Next() {
		Categories.Scan(&category)
		forumpage.Categories = append(forumpage.Categories, category)
	}
	for PostInf.Next() {
		var categories string
		PostInf.Scan(&Post.Id, &Post.Title, &Post.Username, &categories)
		Post.Category = strings.Fields(categories)
		forumpage.Allposts = append(forumpage.Allposts, Post)
	}
	if len(forumpage.Allposts) < 1 {
		forumpage.ErrorString = "There is nothing here..."
	}
	forumpage.IsConnect = true
	forumpage.Username = session.Username
	if er := ExeTemp(w, "forum.html", forumpage); er != nil {
		Error500(w, r)
		return
	}
}

func LikedPostFilter(w http.ResponseWriter, r *http.Request) {
	if ErrorUrl(r) {
		Error404(w, r)
		return

	}

	c, session, err := Authenticated(w, r, "login.html", nil)
	if err != nil || session == nil {
		return
	}
	if c != nil {
		sessionID := c.Value
		e := RefreshSession(sessionID)
		if e != nil {
			Error500(w, r)
			return
		}
	}

	userId := session.User_ID

	Postid, err := DB.Query("SELECT post_id FROM reaction WHERE liked = 1 AND user_id = ?", userId)

	if err != nil {
		Error500(w, r)
		return
	}
	var Post post
	var post_id int
	var forumpage forumhome
	for Postid.Next() {
		Postid.Scan(&post_id)
		var categories string
		err := DB.QueryRow("SELECT p.id, p.title, u.username, p.category_id FROM post p JOIN users u ON p.user_id = u.id WHERE p.id = ?", post_id).Scan(&Post.Id, &Post.Title, &Post.Username, &categories)
		if err != nil {
			Error500(w, r)
			return
		}
		Post.Category = strings.Fields(categories)
		forumpage.Allposts = append(forumpage.Allposts, Post)
	}
	if len(forumpage.Allposts) < 1 {
		forumpage.ErrorString = "There is nothing here..."
	}

	Categories, err2 := DB.Query("SELECT category FROM categories")

	if err2 != nil {
		Error500(w, r)
		return
	}

	var category string

	for Categories.Next() {
		Categories.Scan(&category)
		forumpage.Categories = append(forumpage.Categories, category)
	}

	forumpage.IsConnect = true
	forumpage.Username = session.Username

	if er := ExeTemp(w, "forum.html", forumpage); er != nil {
		Error500(w, r)
		return
	}
}
