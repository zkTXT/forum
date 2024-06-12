//Comment
function showComment() {
    var commentArea = document.getElementById("comment-area");
    commentArea.classList.remove("hide");
}


function upvote(Id) {
    console.log("Upvoting post with ID:", Id);
    fetch("/api/vote", {
        headers: {
            "Content-Type": "application/x-www-form-urlencoded",
        },
        body: "postId=" + Id + "&vote=1",
        method: "POST",
        credentials: "include"
    }).then(response => {
        if (response.ok) {
            console.log("Upvote successful for post ID:", Id);
            location.reload();
        } else {
            console.error("Upvote failed for post ID:", Id, response.statusText);
        }
    }).catch(error => {
        console.error("Error during upvote for post ID:", Id, error);
    });
}

function downvote(Id) {
    console.log("Downvoting post with ID:", Id);
    fetch("/api/vote", {
        headers: {
            "Content-Type": "application/x-www-form-urlencoded",
        },
        body: "postId=" + Id + "&vote=-1",
        method: "POST",
        credentials: "include"
    }).then(response => {
        if (response.ok) {
            console.log("Downvote successful for post ID:", Id);
            location.reload();
        } else {
            console.error("Downvote failed for post ID:", Id, response.statusText);
        }
    }).catch(error => {
        console.error("Error during downvote for post ID:", Id, error);
    });
}