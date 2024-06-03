package forumGO

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// HasUpvoted vérifie si un utilisateur a upvoté un post
func HasUpvoted(database *sql.DB, username string, postId int) bool {
	rows, _ := database.Query("SELECT vote FROM votes WHERE username = ? AND post_id = ? AND vote = 1", username, postId)
	defer rows.Close()
	vote := 0
	for rows.Next() {
		rows.Scan(&vote)
	}
	return vote != 0
}

// HasDownvoted vérifie si un utilisateur a downvoté un post
func HasDownvoted(database *sql.DB, username string, postId int) bool {
	rows, _ := database.Query("SELECT vote FROM votes WHERE username = ? AND post_id = ? AND vote = -1", username, postId)
	defer rows.Close()
	vote := 0
	for rows.Next() {
		rows.Scan(&vote)
	}
	return vote != 0
}

// RemoveVote retire un vote d'un post
func RemoveVote(database *sql.DB, postId int, username string) {
	statement, _ := database.Prepare("DELETE FROM votes WHERE post_id = ? AND username = ?")
	statement.Exec(postId, username)
}

// DecreaseUpvotes diminue le nombre d'upvotes d'un post de 1
func DecreaseUpvotes(database *sql.DB, postId int) {
	statement, _ := database.Prepare("UPDATE posts SET upvotes = upvotes - 1 WHERE id = ?")
	statement.Exec(postId)
}

// DecreaseDownvotes diminue le nombre de downvotes d'un post de 1
func DecreaseDownvotes(database *sql.DB, postId int) {
	statement, _ := database.Prepare("UPDATE posts SET downvotes = downvotes - 1 WHERE id = ?")
	statement.Exec(postId)
}

// IncreaseUpvotes augmente le nombre d'upvotes d'un post de 1
func IncreaseUpvotes(database *sql.DB, postId int) {
	statement, _ := database.Prepare("UPDATE posts SET upvotes = upvotes + 1 WHERE id = ?")
	statement.Exec(postId)
}

// IncreaseDownvotes augmente le nombre de downvotes d'un post de 1
func IncreaseDownvotes(database *sql.DB, postId int) {
	statement, _ := database.Prepare("UPDATE posts SET downvotes = downvotes + 1 WHERE id = ?")
	statement.Exec(postId)
}

// AddVote ajoute un vote à la base de données
func AddVote(database *sql.DB, postId int, username string, vote int) {
	statement, _ := database.Prepare("INSERT INTO votes (username, post_id, vote) VALUES (?, ?, ?)")
	statement.Exec(username, postId, vote)
}

// UpdateVote met à jour le vote d'un utilisateur pour un post
func UpdateVote(database *sql.DB, postId int, username string, vote int) {
	statement, _ := database.Prepare("UPDATE votes SET vote = ? WHERE post_id = ? AND username = ?")
	statement.Exec(vote, postId, username)
}
