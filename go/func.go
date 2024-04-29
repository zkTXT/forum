package forum

import (
	"net/http"
	"text/template"
)

// type Texte struct {
// 	Francais string
// 	Anglais  string
// }

// type AdminCredentials struct {
// 	UsernameL string
// 	PasswordL string
// }

// type AdminPageData struct {
// 	Success bool
// }

// var adminCredentials = AdminCredentials{

// 	UsernameL: "user",
// 	PasswordL: "pass2",
// }

var Error = template.Must(template.ParseFiles("./src/template/error.html"))
var home = template.Must(template.ParseFiles("./src/template/home.html"))
var principal = template.Must(template.ParseFiles("./src/template/principal.html"))
var acc = template.Must(template.ParseFiles("./src/template/acc.html"))

// var login = template.Must(template.ParseFiles("./src/template/login.html"))
// var admin = template.Must(template.ParseFiles("./src/template/admin.html"))

func ErrHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	p := "Error"
	err := Error.ExecuteTemplate(w, "error.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	p := Error
	err := home.ExecuteTemplate(w, "home.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func PrincipalPage(w http.ResponseWriter, r *http.Request) {
	p := Error
	err := principal.ExecuteTemplate(w, "principal.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AccPage(w http.ResponseWriter, r *http.Request) {
	p := Error
	err := acc.ExecuteTemplate(w, "acc.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// func authenticate(username, password string) bool {
// 	return username == adminCredentials.UsernameL && password == adminCredentials.PasswordL
// }

// func LoginHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "POST" {
// 		r.ParseForm()
// 		username := r.Form.Get("UserLogin")
// 		password := r.Form.Get("Password")

// 		if authenticate(username, password) {
// 			http.Redirect(w, r, "/admin", http.StatusSeeOther)
// 			return
// 		} else {
// 			p := "Invalid credentials"
// 			err := login.ExecuteTemplate(w, "login.html", p)
// 			if err != nil {
// 				http.Error(w, err.Error(), http.StatusInternalServerError)
// 				return
// 			}
// 		}
// 	}

// 	p := Error
// 	err := login.ExecuteTemplate(w, "login.html", p)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

// func AdminHandler(w http.ResponseWriter, r *http.Request) {

// 	data := AdminPageData{
// 		Success: false,
// 	}

// 	if r.URL.Query().Get("success") == "true" {
// 		data.Success = true
// 	}

// 	err := admin.ExecuteTemplate(w, "admin.html", data)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

// func ChangeCredentialsHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "POST" {
// 		r.ParseForm()
// 		newUsername := r.Form.Get("NewUsername")
// 		newPassword := r.Form.Get("NewPassword")

// 		adminCredentials.UsernameL = newUsername
// 		adminCredentials.PasswordL = newPassword

// 		http.Redirect(w, r, "/admin?success=true", http.StatusSeeOther)
// 		return
// 	}

// 	http.Redirect(w, r, "/admin", http.StatusSeeOther)
// }

// func NameHandler(w http.ResponseWriter, r *http.Request) {
// 	p := Error
// 	err := name.ExecuteTemplate(w, "name.html", p)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }
