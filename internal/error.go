package internal

import "net/http"


const (
	code200 = 200
	code303 = 303
	code500 = 500
	code400 = 400
	code404 = 404
	code405 = 405
)

func Error500(w http.ResponseWriter, r *http.Request) {
	Errorpage = "Internal server Error"
	code = 500
	http.Redirect(w, r, "/error", code303)
}
func Error405(w http.ResponseWriter, r *http.Request) {
	Errorpage= "Method Not allowed"
	code = 405
	http.Redirect(w, r, "/error", code303)

}
func Error400(w http.ResponseWriter, r *http.Request) {
	Errorpage= "Bad Request"
	code = 400
	http.Redirect(w, r, "/error", code303)
}
func Error404(w http.ResponseWriter, r *http.Request) {
	Errorpage= "Page Not Found"
	code = 404
	http.Redirect(w, r, "/error", code303)
}
func ErrorUrl(r *http.Request) bool {
	path := r.URL.Path
	if path != "/" &&
		path != "/register" &&
		path != "/login" &&
		path != "/post" &&
		path != "/error" &&
		path != "/logout" &&
		path != "/forumInfo" &&
		path != "/forumfilter" &&
		path != "reaction" &&
		path != "/reactioncom" &&
		path != "/comment" &&
		path != "/postSubmit" {
		return true
	}
	return false
}

func Error(w http.ResponseWriter, r *http.Request) {

	if Errorfile {
		data := NewError{
			Mess: Errorpage,
			Code: code,
		}
		w.WriteHeader(code500)
		if er := ExeTemp(w, "error.html", data); er != nil {
			http.Redirect(w, r, "",code500)
			
		}
		return
	}

	switch code {
	case 200:
		w.WriteHeader(code200)
	case 303:
		w.WriteHeader(code303)
	case 400:
		w.WriteHeader(code400)
	case 404:
		w.WriteHeader(code404)
	case 405:
		w.WriteHeader(code405)
	case 500:
		w.WriteHeader(code500)
	}
	data := NewError{
		Mess: Errorpage,
		Code: code,
	}

	if er := ExeTemp(w, "error.html", data); er != nil {
		http.Redirect(w, r, "", code500)
	}
}

