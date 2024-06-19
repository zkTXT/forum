package forumGO

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// CreateUsersTable crée la table des utilisateurs
func CreateUsersTable(database *sql.DB) {
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT, email TEXT, password TEXT, cookie TEXT, expires TEXT)")
	statement.Exec()
}

// CreatePostTable crée la table des posts
func CreatePostTable(database *sql.DB) {
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS posts (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, title TEXT, categories TEXT, content TEXT, created_at TEXT, upvotes INTEGER, downvotes INTEGER)")
	statement.Exec()
}

// CreateCommentTable crée la table des commentaires
func CreateCommentTable(database *sql.DB) {
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS comments (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, post_id INTEGER, content TEXT, created_at TEXT)")
	statement.Exec()
}

// CreateVoteTable crée la table des votes
func CreateVoteTable(database *sql.DB) {
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS votes (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, post_id INTEGER, vote INTEGER)")
	statement.Exec()
}

// CreateCategoriesTable crée la table des catégories
func CreateCategoriesTable(database *sql.DB) {
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS categories (id INTEGER PRIMARY KEY, name TEXT, icon TEXT)")
	statement.Exec()
}

// CreateCategories crée les catégories par défaut
func CreateCategories(database *sql.DB) {
	statement, _ := database.Prepare("INSERT INTO categories (name) SELECT ? WHERE NOT EXISTS (SELECT 1 FROM categories WHERE name = ?)")
	categories := []string{"Healthcare and Medicine", "Technology and IT", "Finance and Accounting", "Science and Research", "Gaming", "Music", "Education and Teaching", "Movies", "TV", "Gastronomy", "Pilot and Flight Operations", "Photography", "Art", "Writing", "Programming", "Other"}
	for _, category := range categories {
		statement.Exec(category, category)
	}
}

// CreateCategoriesIcons met à jour les icônes des catégories
func CreateCategoriesIcons(database *sql.DB) {
	statement, _ := database.Prepare("UPDATE categories SET icon = ? WHERE name = ?")
	icons := map[string]string{
		"Healthcare and Medicinel":    "fa-globe",
		"Technology and IT":           "fa-laptop",
		"Finance and Accounting":      "fa-laptop",
		"Science and Research":        "fa-flask",
		"Gaming":                      "fa-gamepad",
		"Music":                       "fa-music",
		"Education and Teaching":      "fa-book",
		"Movies":                      "fa-film",
		"TV":                          "fa-tv",
		"Gastronomy":                  "fa-cutlery",
		"Pilot and Flight Operations": "fa-plane",
		"Photography":                 "fa-camera",
		"Art":                         "fa-paint-brush",
		"Writing":                     "fa-pencil",
		"Programming":                 "fa-code",
		"Other":                       "fa-question",
	}
	for name, icon := range icons {
		statement.Exec(icon, name)
	}
}
