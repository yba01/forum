package main

import (
	"fmt"
	data "forum/database"
	in "forum/internal"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

const port = ":8081"


func main() {
	tmpl, err := template.ParseGlob("templates/*")
	if err != nil {
		in.Errorfile = true
	}
	// send a template
	if tmpl == nil{
		fmt.Println("Error Parsing file")
		return
	}
	in.Tmpl = tmpl
	// initialize the database
	InitDB, err := data.InitDB()
	if err != nil {
		in.Errorfile = true
	}
	
	if InitDB == nil{
		fmt.Println("Error Opening database")
		return
	}
	in.DB = InitDB
	// create all tables in our database initDB
	if err := data.CreateTables(InitDB); err != nil {
		in.Errorfile = true // if we have bad request sql in our query sql, we have an internal error
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// handle root

	http.HandleFunc("/register", in.Register)
	http.HandleFunc("/registerAuth", in.RegisterAuth)
	http.HandleFunc("/logout", in.Logout)
	http.HandleFunc("/login", in.Login)
	http.HandleFunc("/loginAuth", in.LoginAuth)

	http.HandleFunc("/error", in.Error)

	http.HandleFunc("/", in.Forum)
	http.HandleFunc("/forumInfo", in.ForumInfo)
	http.HandleFunc("/forumfilter", in.Filterforum)
	http.HandleFunc("/reaction", in.ReactionHandle)
	http.HandleFunc("/reactioncom", in.ReactioncomHandle)
	http.HandleFunc("/comment", in.ComentaryInsert)
	http.HandleFunc("/post", in.Post)
	http.HandleFunc("/postSubmit", in.PostSubmit)
	fmt.Println("Server is running on http://localhost"+port)
	http.ListenAndServe(port, nil)

}
