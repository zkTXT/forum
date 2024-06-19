package webAPI

import (
	"FORUM-GO/forumGO"
	"database/sql"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	IsLoggedIn bool
	Username   string
}

type HomePage struct {
	User              User
	Categories        []string
	Icons             []string
	PostsByCategories [][]forumGO.Post
}

type PostsPage struct {
	User  User
	Title string
	Posts []forumGO.Post
	Icon  string
}

type PostPage struct {
	User User
	Post forumGO.Post
}

var database *sql.DB

func SetDatabase(db *sql.DB) {
	database = db
}

func init() {
	var err error
	database, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		panic(err)
	}
}

// Index handles the home page display
func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Fetch categories and posts from the database
	categories := forumGO.GetCategories(database)
	icons := forumGO.GetCategoriesIcons(database)
	postsByCategories := forumGO.GetPostsByCategories(database)

	// Check if alphabetical sorting is requested
	alphabetical := r.URL.Query().Get("alphabetical") == "true"
	if alphabetical {
		sort.Strings(categories)
	}

	// Check if user is logged in
	if isLoggedIn(r) {
		cookie, _ := r.Cookie("SESSION")
		username := forumGO.GetUser(database, cookie.Value)
		payload := HomePage{
			User:              User{IsLoggedIn: true, Username: username},
			Categories:        categories,
			Icons:             icons,
			PostsByCategories: postsByCategories,
		}
		t, _ := template.ParseGlob("public/HTML/*.html")
		t.ExecuteTemplate(w, "forum.html", payload)
		return
	}

	payload := HomePage{
		User:              User{IsLoggedIn: false},
		Categories:        categories,
		Icons:             icons,
		PostsByCategories: postsByCategories,
	}
	t, _ := template.ParseGlob("public/HTML/*.html")
	t.ExecuteTemplate(w, "forum.html", payload)
}

// DisplayPost renders a specific post
func DisplayPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	postID := r.URL.Query().Get("id")
	post := forumGO.GetPost(database, postID)
	comments := forumGO.GetComments(database, postID)
	post.Comments = comments
	payload := PostPage{
		Post: post,
	}
	if isLoggedIn(r) {
		cookie, _ := r.Cookie("SESSION")
		username := forumGO.GetUser(database, cookie.Value)
		payload.User = User{IsLoggedIn: true, Username: username}
	} else {
		payload.User = User{IsLoggedIn: false}
	}
	t, _ := template.ParseGlob("public/HTML/*.html")
	t.ExecuteTemplate(w, "detail.html", payload)
}

// GetPostsByApi filters and retrieves posts based on criteria
func GetPostsByApi(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("by")
	switch filter {
	case "category":
		category := r.URL.Query().Get("category")
		posts := forumGO.GetPostsByCategory(database, category)
		icon := forumGO.GetCategoryIcon(database, category)
		payload := PostsPage{
			Title: "Posts in category " + category,
			Posts: posts,
			Icon:  icon,
		}
		if isLoggedIn(r) {
			cookie, _ := r.Cookie("SESSION")
			username := forumGO.GetUser(database, cookie.Value)
			payload.User = User{IsLoggedIn: true, Username: username}
		}
		t, _ := template.ParseGlob("public/HTML/*.html")
		t.ExecuteTemplate(w, "posts.html", payload)
	case "myposts":
		if isLoggedIn(r) {
			cookie, _ := r.Cookie("SESSION")
			username := forumGO.GetUser(database, cookie.Value)
			posts := forumGO.GetPostsByUser(database, username)
			payload := PostsPage{
				User:  User{IsLoggedIn: true, Username: username},
				Title: "My posts",
				Posts: posts,
				Icon:  "fa-user",
			}
			t, _ := template.ParseGlob("public/HTML/*.html")
			t.ExecuteTemplate(w, "posts.html", payload)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	case "liked":
		if isLoggedIn(r) {
			cookie, _ := r.Cookie("SESSION")
			username := forumGO.GetUser(database, cookie.Value)
			likedPosts := forumGO.GetLikedPosts(database, username)
			payload := PostsPage{
				User:  User{IsLoggedIn: true, Username: username},
				Title: "Posts liked by me",
				Posts: likedPosts,
				Icon:  "fa-heart",
			}
			t, _ := template.ParseGlob("public/HTML/*.html")
			t.ExecuteTemplate(w, "posts.html", payload)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	default:
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// NewPost renders the page for creating a new post
func NewPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	categories := forumGO.GetCategories(database)
	if isLoggedIn(r) {
		cookie, _ := r.Cookie("SESSION")
		username := forumGO.GetUser(database, cookie.Value)
		payload := HomePage{
			User:       User{IsLoggedIn: true, Username: username},
			Categories: categories,
		}
		t, _ := template.ParseGlob("public/HTML/*.html")
		t.ExecuteTemplate(w, "createThread.html", payload)
		return
	}
	payload := HomePage{
		User:       User{IsLoggedIn: false},
		Categories: categories,
	}
	t, _ := template.ParseGlob("public/HTML/*.html")
	t.ExecuteTemplate(w, "createThread.html", payload)
}

// inArray checks if a string exists within a slice
func inArray(input string, array []string) bool {
	for _, element := range array {
		if element == input {
			return true
		}
	}
	return false
}

func Admin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !isLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	t, _ := template.ParseGlob("public/HTML/*.html")
	t.ExecuteTemplate(w, "admin.html", nil)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	statement, err := database.Prepare("DELETE FROM posts WHERE id = ?")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	_, err = statement.Exec(postID)
	if err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	commentID, err := strconv.Atoi(r.FormValue("comment_id"))
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	statement, err := database.Prepare("DELETE FROM comments WHERE id = ?")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	_, err = statement.Exec(commentID)
	if err != nil {
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func UpdateRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	role := r.FormValue("role")
	if role != "user" && role != "moderator" && role != "admin" {
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	statement, err := database.Prepare("UPDATE users SET role = ? WHERE id = ?")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	_, err = statement.Exec(role, userID)
	if err != nil {
		http.Error(w, "Failed to update role", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func Profil(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !isLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	cookie, _ := r.Cookie("SESSION")
	username := forumGO.GetUser(database, cookie.Value)
	_, email, _ := forumGO.GetUserInfo(database, cookie.Value)

	payload := struct {
		IsLoggedIn bool
		Username   string
		Email      string
	}{
		IsLoggedIn: true,
		Username:   username,
		Email:      email,
	}

	t, _ := template.ParseGlob("public/HTML/*.html")
	t.ExecuteTemplate(w, "profil.html", payload)
}

func NewCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !isLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	t, _ := template.ParseGlob("public/HTML/*.html")
	t.ExecuteTemplate(w, "newcategory.html", nil)
}

func AddCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !isLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	icon := r.FormValue("icon")

	// Insert new category into the database
	statement, err := database.Prepare("INSERT INTO categories (name, icon) VALUES (?, ?)")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	_, err = statement.Exec(name, icon)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Redirect to a success page or back to the form
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// GetLatestPostByUser retrieves the most recent post by the given username
func GetLatestPostByUser(database *sql.DB, username string) (forumGO.Post, error) {
	var post forumGO.Post
	query := `SELECT id, username, title, categories, content, created_at, upvotes, downvotes FROM posts WHERE username = ? ORDER BY created_at DESC LIMIT 1`
	row := database.QueryRow(query, username)
	var categories string
	err := row.Scan(&post.Id, &post.Username, &post.Title, &categories, &post.Content, &post.CreatedAt, &post.UpVotes, &post.DownVotes)
	if err != nil {
		return post, err
	}
	post.Categories = strings.Split(categories, ",")
	return post, nil
}

func UserProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !isLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	username := r.URL.Query().Get("username")
	email, _ := forumGO.GetUserByUsername(database, username)

	latestPost, err := forumGO.GetLatestPostByUser(database, username)
	if err != nil {
		// Handle error if needed
	}

	payload := struct {
		IsLoggedIn bool
		Username   string
		Email      string
		LatestPost forumGO.Post
	}{
		IsLoggedIn: true,
		Username:   username,
		Email:      email,
		LatestPost: latestPost,
	}

	t, _ := template.ParseGlob("public/HTML/*.html")
	t.ExecuteTemplate(w, "userprofile.html", payload)
}

// UpdateUsername handles the username update request
func UpdateUsername(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !isLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	cookie, err := r.Cookie("SESSION")
	if err != nil {
		http.Error(w, "Could not get session cookie", http.StatusInternalServerError)
		return
	}

	oldUsername := forumGO.GetUser(database, cookie.Value)
	newUsername := r.FormValue("newUsername")

	// Check if the new username is taken
	if !forumGO.UsernameNotTaken(database, newUsername) {
		http.Error(w, "Username already taken", http.StatusBadRequest)
		return
	}

	// Update the username in the database
	err = forumGO.UpdateUsername(database, oldUsername, newUsername)
	if err != nil {
		http.Error(w, "Could not update username", http.StatusInternalServerError)
		return
	}

	// Redirect to profile page
	http.Redirect(w, r, "/profil", http.StatusSeeOther)
}

// UpdatePassword handles the password update request
func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !isLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	cookie, err := r.Cookie("SESSION")
	if err != nil {
		http.Error(w, "Could not get session cookie", http.StatusInternalServerError)
		return
	}

	username := forumGO.GetUser(database, cookie.Value)
	currentPassword := r.FormValue("currentPassword")
	newPassword := r.FormValue("newPassword")
	confirmPassword := r.FormValue("confirmPassword")

	if newPassword != confirmPassword {
		http.Error(w, "New passwords do not match", http.StatusBadRequest)
		return
	}

	// Verify current password
	storedPassword := forumGO.GetPasswordByUsername(database, username)
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(currentPassword))
	if err != nil {
		http.Error(w, "Current password is incorrect", http.StatusBadRequest)
		return
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Could not hash new password", http.StatusInternalServerError)
		return
	}

	// Update the password in the database
	err = forumGO.UpdatePassword(database, username, string(hashedPassword))
	if err != nil {
		http.Error(w, "Could not update password", http.StatusInternalServerError)
		return
	}

	// Redirect to profile page
	http.Redirect(w, r, "/profil", http.StatusSeeOther)
}

func GetUserRole(username string) (string, error) {
	var role string
	err := database.QueryRow("SELECT role FROM users WHERE username = ?", username).Scan(&role)
	if err != nil {
		return "", err
	}
	return role, nil
}
