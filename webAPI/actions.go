package webAPI

import (
	"FORUM-GO/forumGO"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Vote struct {
	PostId int
	Vote   int
}

// CreatePostApi manages the creation of posts
func CreatePostApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("ParseForm() error: %v", err), http.StatusInternalServerError)
		return
	}
	if !isLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	cookie, _ := r.Cookie("SESSION")
	username := forumGO.GetUser(database, cookie.Value)
	title := r.FormValue("title")
	content := r.FormValue("content")
	categories := r.Form["categories[]"]
	validCategories := forumGO.GetCategories(database)

	for _, category := range categories {
		if !inArray(category, validCategories) {
			http.Error(w, "Invalid category: "+category, http.StatusBadRequest)
			return
		}
	}

	stringCategories := strings.Join(categories, ",")
	now := time.Now()
	forumGO.CreatePost(database, username, title, stringCategories, content, now)
	fmt.Printf("Post created by %s with title %s at %s\n", username, title, now.Format("2006-01-02 15:04:05"))
	http.Redirect(w, r, "/filter?by=myposts", http.StatusFound)
}

// CommentsApi manages the creation of comments
func CommentsApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("ParseForm() error: %v", err), http.StatusInternalServerError)
		return
	}
	if !isLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	cookie, _ := r.Cookie("SESSION")
	username := forumGO.GetUser(database, cookie.Value)
	postId := r.FormValue("postId")
	content := r.FormValue("content")
	now := time.Now()
	postIdInt, _ := strconv.Atoi(postId)
	forumGO.AddComment(database, username, postIdInt, content, now)
	fmt.Printf("Comment created by %s on post %s at %s\n", username, postId, now.Format("2006-01-02 15:04:05"))
	http.Redirect(w, r, "/post?id="+postId, http.StatusFound)
}

// VoteApi manages votes on posts
func VoteApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !isLoggedIn(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("ParseForm() error: %v", err), http.StatusInternalServerError)
		return
	}
	cookie, _ := r.Cookie("SESSION")
	username := forumGO.GetUser(database, cookie.Value)
	postId := r.FormValue("postId")
	postIdInt, _ := strconv.Atoi(postId)
	vote := r.FormValue("vote")
	voteInt, _ := strconv.Atoi(vote)
	now := time.Now().Format("2006-01-02 15:04:05")

	if voteInt == 1 {
		if forumGO.HasUpvoted(database, username, postIdInt) {
			forumGO.RemoveVote(database, postIdInt, username)
			forumGO.DecreaseUpvotes(database, postIdInt)
			fmt.Printf("Removed upvote from %s on post %s at %s\n", username, postId, now)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Vote removed"))
			return
		}
		if forumGO.HasDownvoted(database, username, postIdInt) {
			forumGO.DecreaseDownvotes(database, postIdInt)
			forumGO.IncreaseUpvotes(database, postIdInt)
			forumGO.UpdateVote(database, postIdInt, username, 1)
			fmt.Printf("%s upvoted on post %s at %s\n", username, postId, now)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Upvote added"))
			return
		}
		forumGO.IncreaseUpvotes(database, postIdInt)
		forumGO.AddVote(database, postIdInt, username, 1)
		fmt.Printf("%s upvoted on post %s at %s\n", username, postId, now)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Upvote added"))
		return
	}

	if voteInt == -1 {
		if forumGO.HasDownvoted(database, username, postIdInt) {
			forumGO.RemoveVote(database, postIdInt, username)
			forumGO.DecreaseDownvotes(database, postIdInt)
			fmt.Printf("Removed downvote from %s on post %s at %s\n", username, postId, now)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Vote removed"))
			return
		}
		if forumGO.HasUpvoted(database, username, postIdInt) {
			forumGO.DecreaseUpvotes(database, postIdInt)
			forumGO.IncreaseDownvotes(database, postIdInt)
			forumGO.UpdateVote(database, postIdInt, username, -1)
			fmt.Printf("%s downvoted on post %s at %s\n", username, postId, now)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Downvote added"))
			return
		}
		forumGO.IncreaseDownvotes(database, postIdInt)
		forumGO.AddVote(database, postIdInt, username, -1)
		fmt.Printf("%s downvoted on post %s at %s\n", username, postId, now)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Downvote added"))
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Invalid vote"))
}
