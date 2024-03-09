package internal

import (
	"fmt"
	"net/http"
	"strconv"
)

// The Post function handles POST requests, checks for errors, authenticates users, refreshes sessions,
// queries categories from the database, and executes a template to render the post.html page.
func Post(w http.ResponseWriter, r *http.Request) {

	if Errorfile {
		Error500(w, r)
		return
	}

	c , _,err := Authenticated(w,r,"login.html",nil)
	if err != nil{
		return
	}
	if c != nil{
		sessionID := c.Value
		e := RefreshSession(sessionID)
		if e != nil{
			Error500(w,r)
			return
		}
	}

	Categories, err := DB.Query("SELECT id, category FROM categories")
	if err != nil {
		Error500(w, r)
		return
	}

	categories := make(map[int]string)

	var id int
	var category string

	for Categories.Next() {
		Categories.Scan(&id, &category)
		categories[id] = category
	}
	
	if er := ExeTemp(w, "post.html", categories); er != nil {
		Error500(w, r)
		return
	}
}

// The function `PostSubmit` handles the submission of a post with validation checks and database
// operations in a Go web application.
func PostSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error405(w, r)
		return
	}

	_, Session, e := Authenticated(w, r,"login.html",nil)

	if e != nil || Session == nil {
		return
	}
	
	sessionID := Session.User_ID

	title := r.FormValue("subject")

	if len(title) > 30 {
		Error400(w, r)
		return
	}

	message := r.FormValue("Thepost")
	category := r.Form["category"]

	if title == "" || message == "" || len(category) == 0 {
		Error400(w, r)
		return
	}

	var category_id string
	for _, id := range category {
		categoryid, er := strconv.Atoi(id)
		if er != nil {
			Error400(w, r)
			return
		}
		if categoryid < 1 || categoryid > 6 {
			Error400(w, r)
			return
		}
		var onecategory string
		queryerr := DB.QueryRow("SELECT category FROM categories WHERE id = ?", categoryid).Scan(&onecategory)

		if queryerr != nil {
			Error400(w, r)
			return
		}
		category_id += " " + onecategory + " "
	}

	post := &PostValue{
		Title:   title,
		User_id: sessionID,
		Message: message,
		Category_id: category_id,
	}
	err := CreatePost(post)

	if err != nil {
		Error500(w, r)
		return
	}

	http.Redirect(w, r, "/", code303)
}

// The CreatePost function inserts a new post into a database table.
func CreatePost(post *PostValue) error {
	_, err := DB.Exec("INSERT INTO post (user_id, category_id, title, message) VALUES(?, ?, ?, ?)",  post.User_id, post.Category_id, post.Title, post.Message)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}
 