package forumGO

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// GetPost by id returns a Post struct with the post data
func GetPost(database *sql.DB, id string) Post {
	query := "SELECT username, title, categories, content, created_at, upvotes, downvotes FROM posts WHERE id = ?"
	rows, _ := database.Query(query, id)
	defer rows.Close()

	var post Post
	post.Id, _ = strconv.Atoi(id)
	for rows.Next() {
		var catString string
		rows.Scan(&post.Username, &post.Title, &catString, &post.Content, &post.CreatedAt, &post.UpVotes, &post.DownVotes)
		post.Categories = strings.Split(catString, ",")
	}
	return post
}

// GetComments get comments by post id
func GetComments(database *sql.DB, id string) []Comment {
	query := "SELECT id, username, content, created_at FROM comments WHERE post_id = ?"
	rows, _ := database.Query(query, id)
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		rows.Scan(&comment.Id, &comment.Username, &comment.Content, &comment.CreatedAt)
		comments = append(comments, comment)
	}
	return comments
}

// GetPostsByCategory returns all posts in a given category
func GetPostsByCategory(database *sql.DB, category string) []Post {
	query := "SELECT id, username, title, categories, content, created_at, upvotes, downvotes FROM posts WHERE categories LIKE ?"
	rows, _ := database.Query(query, "%"+category+"%")
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var catString string
		rows.Scan(&post.Id, &post.Username, &post.Title, &catString, &post.Content, &post.CreatedAt, &post.UpVotes, &post.DownVotes)
		post.Categories = strings.Split(catString, ",")
		posts = append(posts, post)
	}
	return posts
}

// GetPostsByCategories returns all posts for all categories
func GetPostsByCategories(database *sql.DB) [][]Post {
	categories := GetCategories(database)
	var posts [][]Post
	for _, category := range categories {
		categoryPosts := GetPostsByCategory(database, category)
		posts = append(posts, categoryPosts)
	}
	return posts
}

// GetPostsByUser returns all posts by a user
func GetPostsByUser(database *sql.DB, username string) []Post {
	query := "SELECT id, username, title, categories, content, created_at, upvotes, downvotes FROM posts WHERE username = ?"
	rows, _ := database.Query(query, username)
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var catString string
		rows.Scan(&post.Id, &post.Username, &post.Title, &catString, &post.Content, &post.CreatedAt, &post.UpVotes, &post.DownVotes)
		post.Categories = strings.Split(catString, ",")
		posts = append(posts, post)
	}
	return posts
}

// GetLikedPosts gets posts that user has liked
func GetLikedPosts(database *sql.DB, username string) []Post {
	query := "SELECT id, username, title, categories, content, created_at, upvotes, downvotes FROM posts WHERE id IN (SELECT post_id FROM votes WHERE username = ? AND vote = 1)"
	rows, _ := database.Query(query, username)
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var catString string
		rows.Scan(&post.Id, &post.Username, &post.Title, &catString, &post.Content, &post.CreatedAt, &post.UpVotes, &post.DownVotes)
		post.Categories = strings.Split(catString, ",")
		posts = append(posts, post)
	}
	return posts
}

// GetCategories returns all categories
func GetCategories(database *sql.DB) []string {
	query := "SELECT name FROM categories"
	rows, _ := database.Query(query)
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var name string
		rows.Scan(&name)
		categories = append(categories, name)
	}
	return categories
}

// GetCategoriesIcons returns all categories' icons
func GetCategoriesIcons(database *sql.DB) []string {
	query := "SELECT icon FROM categories"
	rows, _ := database.Query(query)
	defer rows.Close()

	var icons []string
	for rows.Next() {
		var icon string
		rows.Scan(&icon)
		icons = append(icons, icon)
	}
	return icons
}

// GetCategoryIcon returns the icon for a category
func GetCategoryIcon(database *sql.DB, category string) string {
	query := "SELECT icon FROM categories WHERE name = ?"
	rows, _ := database.Query(query, category)
	defer rows.Close()

	var icon string
	for rows.Next() {
		rows.Scan(&icon)
	}
	return icon
}

// CreatePost
func CreatePost(database *sql.DB, username string, title string, categories string, content string, createdAt time.Time) {
	createdAtString := createdAt.Format("2006-01-02 15:04:05")
	query := "INSERT INTO posts (username, title, categories, content, created_at, upvotes, downvotes) VALUES (?, ?, ?, ?, ?, ?, ?)"
	statement, _ := database.Prepare(query)
	statement.Exec(username, title, categories, content, createdAtString, 0, 0)
}

// AddComment adds a comment to a post
func AddComment(database *sql.DB, username string, postId int, content string, createdAt time.Time) {
	createdAtString := createdAt.Format("2006-01-02 15:04:05")
	query := "INSERT INTO comments (username, post_id, content, created_at) VALUES (?, ?, ?, ?)"
	statement, _ := database.Prepare(query)
	statement.Exec(username, postId, content, createdAtString)
}
