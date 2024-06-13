package webAPI

import (
	"FORUM-GO/forumGO"
	"fmt"
	"html/template"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type Error struct {
	Message string
}

func DeleteCommentAPI(w http.ResponseWriter, r *http.Request) {
	userRole := getUserRoleFromRequest(r)
	if userRole != "moderator" && userRole != "admin" {
		http.Error(w, "Vous n'êtes pas autorisé à effectuer cette action", http.StatusUnauthorized)
		return
	}

	commentID := r.FormValue("id") // Utiliser "id" au lieu de "commentID"
	if commentID == "" {
		http.Error(w, "ID du commentaire manquant", http.StatusBadRequest)
		return
	}

	if _, err := database.Exec("DELETE FROM comments WHERE id = ?", commentID); err != nil {
		http.Error(w, "Erreur lors de la suppression du commentaire", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Le commentaire a été supprimé avec succès"))
}

// DeletePostAPI supprime une publication spécifiée par son ID (accessible uniquement aux modérateurs et aux administrateurs)
func DeletePostAPI(w http.ResponseWriter, r *http.Request) {
	userRole := getUserRoleFromRequest(r)
	if userRole != "moderator" && userRole != "admin" {
		http.Error(w, "Vous n'êtes pas autorisé à effectuer cette action", http.StatusUnauthorized)
		return
	}

	postID := r.FormValue("id") // Utiliser "id" au lieu de "postID"
	if postID == "" {
		http.Error(w, "ID de la publication manquant", http.StatusBadRequest)
		return
	}

	if _, err := database.Exec("DELETE FROM posts WHERE id = ?", postID); err != nil {
		http.Error(w, "Erreur lors de la suppression de la publication", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("La publication a été supprimée avec succès"))
}

// getUserRoleFromRequest récupère le rôle de l'utilisateur à partir de la requête HTTP
func getUserRoleFromRequest(r *http.Request) string {
	cookie, err := r.Cookie("user_role")
	if err != nil {
		return "guest"
	}
	return cookie.Value
}

// RegisterApi handles the registration process
func RegisterApi(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("ParseForm() error: %v", err), http.StatusInternalServerError)
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	sessionValue := uuid.NewV4().String()
	sessionExpiration := time.Now().Add(31 * 24 * time.Hour)

	if username == "" || email == "" || password == "" {
		http.Redirect(w, r, "/register?err=invalid_informations", http.StatusSeeOther)
		return
	}
	if !forumGO.UsernameNotTaken(database, username) {
		http.Redirect(w, r, "/register?err=username_taken", http.StatusSeeOther)
		return
	}
	if !forumGO.EmailNotTaken(database, email) {
		http.Redirect(w, r, "/register?err=email_taken", http.StatusSeeOther)
		return
	}
	forumGO.AddUser(database, username, email, password, sessionValue, sessionExpiration.Format("2006-01-02 15:04:05"))
	sessionCookie := http.Cookie{Name: "SESSION", Value: sessionValue, Expires: sessionExpiration, Path: "/"}
	http.SetCookie(w, &sessionCookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

// LoginApi handles user login
func LoginApi(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("ParseForm() error: %v", err), http.StatusInternalServerError)
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")
	user := r.FormValue("user")

	username, storedEmail, storedPassword := forumGO.GetUserInfo(database, email, user)
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	if username == "" && storedEmail == "" && storedPassword == "" {
		fmt.Printf("Login failed (email not found) for %s at %s\n", email, currentTime)
		http.Redirect(w, r, "/login?err=invalid_email", http.StatusSeeOther)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)); err != nil {
		fmt.Printf("Login failed (wrong password) for %s at %s\n", email, currentTime)
		http.Redirect(w, r, "/login?err=invalid_password", http.StatusSeeOther)
		return
	}
	sessionExpiration := time.Now().Add(31 * 24 * time.Hour)
	sessionValue := uuid.NewV4().String()
	sessionCookie := http.Cookie{Name: "SESSION", Value: sessionValue, Expires: sessionExpiration, Path: "/"}
	http.SetCookie(w, &sessionCookie)
	forumGO.UpdateCookie(database, sessionValue, sessionExpiration, storedEmail)
	fmt.Printf("Logged in user: %s with email: %s at %s\n", username, storedEmail, currentTime)
	http.Redirect(w, r, "/", http.StatusFound)
}

// LogoutAPI handles user logout
func LogoutAPI(w http.ResponseWriter, r *http.Request) {
	sessionCookie, _ := r.Cookie("SESSION")
	username := forumGO.GetUser(database, sessionCookie.Value)
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	if sessionCookie != nil {
		username := forumGO.GetUser(database, sessionCookie.Value)
		forumGO.Logout(database, username)
	}
	fmt.Printf("User %s logged out at %s\n", username, currentTime)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// isLoggedIn verifies user login status
func isLoggedIn(r *http.Request) bool {
	sessionCookie, err := r.Cookie("SESSION")
	if err != nil {
		return false
	}
	if !forumGO.CheckCookie(database, sessionCookie.Value) {
		return false
	}
	expiration := forumGO.GetExpires(database, sessionCookie.Value)
	return !isExpired(expiration)
}

// isExpired checks if a session is expired
func isExpired(expiration string) bool {
	expirationTime, _ := time.Parse("2006-01-02 15:04:05", expiration)
	return time.Now().After(expirationTime)
}

// Register displays the registration page
func Register(w http.ResponseWriter, r *http.Request) {
	error := r.URL.Query().Get("err")
	payload := Error{Message: ""}
	switch error {
	case "invalid_informations":
		payload.Message = "Invalid informations"
	case "email_taken":
		payload.Message = "Email already taken"
	case "username_taken":
		payload.Message = "Username already taken"
	}
	t, _ := template.ParseGlob("public/HTML/*.html")
	t.ExecuteTemplate(w, "registerForm.html", payload)
}

// Login displays the login page
func Login(w http.ResponseWriter, r *http.Request) {
	error := r.URL.Query().Get("err")
	payload := Error{Message: ""}
	switch error {
	case "invalid_email":
		payload.Message = "Invalid email"
	case "invalid_password":
		payload.Message = "Invalid password"
	}
	t, _ := template.ParseGlob("public/HTML/*.html")
	t.ExecuteTemplate(w, "signinForm.html", payload)
}
