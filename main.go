package main

import (
	"FORUM-GO/forumGO"
	"FORUM-GO/webAPI"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Post struct {
	Id         int
	Username   string
	Title      string
	Categories []string
	Content    string
	CreatedAt  string
	UpVotes    int
	DownVotes  int
	Comments   []Comment
}

type Comment struct {
	Id        int
	PostId    int
	Username  string
	Content   string
	CreatedAt string
}

// Database
var database *sql.DB

func main() {
	// check if DB exists
	var _, err = os.Stat("database.db")

	// create DB if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create("database.db")
		if err != nil {
			return
		}
		defer file.Close()
	}

	database, _ = sql.Open("sqlite3", "./database.db")

	forumGO.CreateUsersTable(database)
	forumGO.CreatePostTable(database)
	forumGO.CreateCommentTable(database)
	forumGO.CreateVoteTable(database)
	forumGO.CreateCategoriesTable(database)
	forumGO.CreateCategories(database)
	forumGO.CreateCategoriesIcons(database)

	webAPI.SetDatabase(database)

	fs := http.FileServer(http.Dir("public"))
	router := http.NewServeMux()
	fmt.Println("Starting server on port 8080")
	fmt.Println("http://localhost:8000:")

	router.HandleFunc("/", webAPI.Index)
	router.HandleFunc("/register", webAPI.Register)
	router.HandleFunc("/login", webAPI.Login)
	router.HandleFunc("/admin", webAPI.Admin)
	router.HandleFunc("/profil", webAPI.Profil)
	router.HandleFunc("/post", webAPI.DisplayPost)
	router.HandleFunc("/filter", webAPI.GetPostsByApi)
	router.HandleFunc("/newpost", webAPI.NewPost)
	router.HandleFunc("/newcategory", webAPI.NewCategory)
	router.HandleFunc("/add-category", webAPI.AddCategory)
	router.HandleFunc("/api/register", webAPI.RegisterApi)
	router.HandleFunc("/api/login", webAPI.LoginApi)
	router.HandleFunc("/api/logout", webAPI.LogoutAPI)
	router.HandleFunc("/api/createpost", webAPI.CreatePostApi)
	router.HandleFunc("/api/comments", webAPI.CommentsApi)
	router.HandleFunc("/api/vote", webAPI.VoteApi)
	router.HandleFunc("/userprofile", webAPI.UserProfile)
	router.HandleFunc("/api/updateUser", webAPI.UpdateUsername)
	router.HandleFunc("/api/updatePassword", webAPI.UpdatePassword)

	router.Handle("/public/", http.StripPrefix("/public/", fs))
	http.ListenAndServe(":8000", router)
}
