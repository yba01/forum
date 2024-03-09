package internal

import (
	"database/sql"
	"html/template"
	"io"
)

var (
	Errorfile bool   // internal error file or database
	Errorpage string // error page
	code      int
	Tmpl      *template.Template
	DB        *sql.DB
)

type data struct {
	IsReg      bool
	ErrorLog   string
	ErrorAlert string
}

type NewError struct {
	Mess string
	Code int
}

type User struct {
	ID           int
	Username     string
	Email        string
	Password     string
	HashPassword string
}

type Session struct {
	Token      string
	User_ID    int
	Username   string
	Create_at  string
	Expired_at string
}

type PostValue struct {
	ID          int
	User_id     int
	Category_id string
	Title       string
	Message     string
}

type comment struct {
	Com_id   int
	Like     int
	Dislike  int
	Username string
	Comment  string
}

type ppost struct {
	User        string
	Post_id     string
	User_id     int
	Title       string
	Message     string
	Comments    []comment
	Likes       int
	Dislikes    int
	IsConnect   bool
	Online_user string
}

type post struct {
	Id       int
	Username string
	Title    string
	Category []string
}

type forumhome struct {
	Allposts    []post
	Categories  []string
	IsConnect   bool
	Username    string
	ErrorString string
}

func ExeTemp(w io.Writer, name string, data any) error {
	return Tmpl.ExecuteTemplate(w, name, data)
}
