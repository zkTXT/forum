package webAPI

import (
	"FORUM-GO/forumGO"
	"database/sql"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
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

// Index handles the home page display
func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	categories := forumGO.GetCategories(database)
	icons := forumGO.GetCategoriesIcons(database)
	postsByCategories := forumGO.GetPostsByCategories(database)
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
	if !isLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	t, _ := template.ParseGlob("public/HTML/*.html")
	t.ExecuteTemplate(w, "createThread.html", nil)
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

func Profil(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !isLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	t, _ := template.ParseGlob("public/HTML/*.html")
	t.ExecuteTemplate(w, "profil.html", nil)
}

func ProfilOther(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !isLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	t, _ := template.ParseGlob("public/HTML/*.html")
	t.ExecuteTemplate(w, "profilother.html", nil)
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
