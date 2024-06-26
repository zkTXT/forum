package forumGO

import (
	"database/sql"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// GetUser retourne le nom d'utilisateur associé à un cookie donné
func GetUser(database *sql.DB, cookie string) string {
	rows, _ := database.Query("SELECT username FROM users WHERE cookie = ?", cookie)
	defer rows.Close()
	var username string
	for rows.Next() {
		rows.Scan(&username)
	}
	return username
}

// GetUserInfo retourne le nom d'utilisateur, l'email et le mot de passe hashé associés à un email donné
func GetUserInfo(database *sql.DB, submittedEmail string) (string, string, string) {
	var user, email, password string
	rows, _ := database.Query("SELECT username, email, password FROM users WHERE email = ?", submittedEmail)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&user, &email, &password)
	}
	return user, email, password
}

// GetUserByUsername récupère les informations d'un utilisateur par son nom d'utilisateur
func GetUserByUsername(database *sql.DB, username string) (string, string) {
	var email, password string
	row := database.QueryRow("SELECT email, password FROM users WHERE username = ?", username)
	err := row.Scan(&email, &password)
	if err != nil {
		return "", ""
	}
	return email, password
}

// GetLatestPostByUser retrieves the most recent post by the given username
func GetLatestPostByUser(database *sql.DB, username string) (Post, error) {
	var post Post
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

// UpdateUsername updates the username of a user in the database
func UpdateUsername(database *sql.DB, oldUsername, newUsername string) error {
	statement, err := database.Prepare("UPDATE users SET username = ? WHERE username = ?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(newUsername, oldUsername)
	return err
}

// GetPasswordByUsername retourne le mot de passe hashé associé à un nom d'utilisateur donné
func GetPasswordByUsername(database *sql.DB, username string) string {
	var password string
	row := database.QueryRow("SELECT password FROM users WHERE username = ?", username)
	row.Scan(&password)
	return password
}

// UpdatePassword met à jour le mot de passe de l'utilisateur dans la base de données
func UpdatePassword(database *sql.DB, username string, newPassword string) error {
	statement, err := database.Prepare("UPDATE users SET password = ? WHERE username = ?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(newPassword, username)
	return err
}
