package forumGO

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// AddUser insère un nouvel utilisateur dans la base de données
func AddUser(database *sql.DB, username string, email string, password string, cookie string, expires string) {
	password, _ = generateHash(password)
	statement, _ := database.Prepare("INSERT INTO users (username, email, password, cookie, expires) VALUES (?, ?, ?, ?, ?)")
	statement.Exec(username, email, password, cookie, expires)
	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println("Added user: " + username + " with email: " + email + " at " + now)
}

// AddModerator crée un nouveau modérateur en mettant à jour le rôle dans la base de données
func AddModerator(database *sql.DB, username string) {
	statement, _ := database.Prepare("UPDATE users SET role = 'moderator' WHERE username = ?")
	statement.Exec(username)
}

// AddAdmin crée un nouvel administrateur en mettant à jour le rôle dans la base de données
func AddAdmin(database *sql.DB, username string) {
	statement, _ := database.Prepare("UPDATE users SET role = 'admin' WHERE username = ?")
	statement.Exec(username)
}

// EmailNotTaken vérifie si un email est déjà utilisé
func EmailNotTaken(database *sql.DB, email string) bool {
	rows, _ := database.Query("SELECT email FROM users WHERE email = ?", email)
	defer rows.Close()
	var emailExists string
	for rows.Next() {
		rows.Scan(&emailExists)
	}
	return emailExists == ""
}

// UsernameNotTaken vérifie si un nom d'utilisateur est déjà utilisé
func UsernameNotTaken(database *sql.DB, username string) bool {
	rows, _ := database.Query("SELECT username FROM users WHERE username = ?", username)
	defer rows.Close()
	var usernameExists string
	for rows.Next() {
		rows.Scan(&usernameExists)
	}
	return usernameExists == ""
}

// CheckCookie vérifie si un cookie est valide
func CheckCookie(database *sql.DB, cookie string) bool {
	var result bool
	err := database.QueryRow("SELECT IIF(COUNT(*), 'true', 'false') FROM users WHERE cookie = ?", cookie).Scan(&result)
	if err != nil {
		return false
	}
	return result
}

// GetExpires retourne la date d'expiration d'un cookie
func GetExpires(database *sql.DB, cookie string) string {
	var expires string
	rows, _ := database.Query("SELECT expires FROM users WHERE cookie = ?", cookie)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&expires)
	}
	return expires
}

// Logout déconnecte un utilisateur
func Logout(database *sql.DB, username string) {
	statement, _ := database.Prepare("UPDATE users SET cookie = '', expires = '' WHERE username = ?")
	statement.Exec(username)
}

// UpdateCookie met à jour le cookie d'un utilisateur
func UpdateCookie(database *sql.DB, token string, expiration time.Time, email string) {
	statement, _ := database.Prepare("UPDATE users SET cookie = ?, expires = ? WHERE email = ?")
	statement.Exec(token, expiration.Format("2006-01-02 15:04:05"), email)
}

// generateHash génère un hash pour le mot de passe
func generateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hash), err
}
