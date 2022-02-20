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

	if err != nil {
		panic(err.Error())
	}
	validToken, err := GetJWT(tag.ID)
	fmt.Println(validToken)
	if err != nil {
		fmt.Println("Failed to generate token")
	}

	fmt.Fprintf(w, "validToken")
	fmt.Fprint(w, validToken)
	if err != nil {
		panic(err.Error()) //error handling
	}
}

// JWT TOKEN GENERATION

func GetJWT(id int) (string, error) {
	var mySigningKey = []byte("unicorns")
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["id"] = id

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

//-------------------------------------------------------------------------------------------------------

func signUp(w http.ResponseWriter, r *http.Request) {


	db, err := sql.Open("mysql", "root@tcp(127.0.0.1)/facebookdb")
	print(err)
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
		firstname := r.FormValue("firstname")
		lastname := r.FormValue("lastname")
		email := r.FormValue("email")
		password := r.FormValue("password")
		fmt.Fprintf(w, "Firstname = %s\n", firstname)
		fmt.Fprintf(w, "Lastname = %s\n", lastname)
		fmt.Fprintf(w, "Email = %s\n", email)
		fmt.Fprintf(w, "Password = %s\n", password)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

	// Execute the query

	firstname := r.FormValue("firstname")
	lastname := r.FormValue("lastname")
	email := r.FormValue("email")
	password := r.FormValue("password")

	var user string

	err = db.QueryRow("SELECT firstname FROM user_account WHERE firstname=?", firstname).Scan(&user)
	print(err)

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Server error, unable to create your account.", 500)
			return
		}
		