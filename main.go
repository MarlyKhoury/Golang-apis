package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type Tag struct {
	ID int `json:"id"`
}

func main() {

	db, err := sql.Open("mysql", "root@tcp(127.0.0.1)/facebookdb")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	http.HandleFunc("/", logIn)
	http.HandleFunc("/signup", signUp)

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func logIn(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {

	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		email := r.FormValue("email")
		password := r.FormValue("password")
		fmt.Fprintf(w, "Email = %s\n", email)
		fmt.Fprintf(w, "Password = %s\n", password)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

	// Execute the query
	var tag Tag
	email := r.FormValue("email")
	password := r.FormValue("password")
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1)/facebookdb")
	print(err)
	err = db.QueryRow("SELECT id FROM user_account WHERE email = ? AND password = ?", email, password).Scan(&tag.ID)

	