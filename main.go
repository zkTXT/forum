package main

import (
	"fmt"
	forum "forum/go"
	"net/http"
)

func main() {
	http.HandleFunc("/home", forum.HomePage)
	http.HandleFunc("/", forum.ErrHandler)
	http.HandleFunc("/principal", forum.PrincipalPage)
	http.HandleFunc("/acc", forum.AccPage)

	fmt.Println("server started...")
	fmt.Println("http://localhost:7070/home")
	http.ListenAndServe(":7070", nil)

}
