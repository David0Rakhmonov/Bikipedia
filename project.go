package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

type Article struct {
	Id                   uint16
	Title, Idea, Article string
}

var posts = []Article{}
var watchPost = Article{}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("C:/www/templates/index.html", "C:/www/templates/header.html", "C:/www/templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	res, err := db.Query("SELECT * FROM articles")

	if err != nil {
		panic(err)
	}

	posts = []Article{}

	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Idea, &post.Article)
		if err != nil {
			panic(err)
		}

		posts = append(posts, post)
	}

	t.ExecuteTemplate(w, "index", posts)
}

// func contacts(w http.ResponseWriter, r *http.Request) {
// 	t, err := template.ParseFiles("C:/www/templates/contacts.html", "C:/www/templates/header.html", "C:/www/templates/footer.html")

// 	if err != nil {
// 		fmt.Fprintf(w, err.Error())
// 	}

// 	t.ExecuteTemplate(w, "contacts", nil)
// }

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("C:/www/templates/create.html", "C:/www/templates/header.html", "C:/www/templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	idea := r.FormValue("idea")
	article := r.FormValue("article")

	// if title == "" || idea == "" || article == "" {
	// 	fmt.Fprintf(w, "Пожалуйста, заполните все поля")
	// } else {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	insert, err := db.Query(fmt.Sprintf("INSERT INTO `articles` (`title`, `idea`, `aricle`) VALUES ('%s', '%s', '%s')", title, idea, article))

	if err != nil {
		panic(err)
	}

	defer insert.Close()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// }

func watch_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	t, err := template.ParseFiles("C:/www/templates/watch.html", "C:/www/templates/header.html", "C:/www/templates/footer.html")

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	res, err := db.Query(fmt.Sprintf("SELECT * FROM articles WHERE id = '%s'", vars["id"]))

	if err != nil {
		panic(err)
	}

	watchPost = Article{}

	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Idea, &post.Article)
		if err != nil {
			panic(err)
		}

		watchPost = post
	}

	t.ExecuteTemplate(w, "watch", watchPost)
}

func search(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("query")

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	res, err := db.Query("SELECT * FROM articles WHERE title LIKE ? OR idea LIKE ? OR aricle LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	if err != nil {
		panic(err)
	}

	searchResults := []Article{}

	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Idea, &post.Article)
		if err != nil {
			panic(err)
		}

		searchResults = append(searchResults, post)
	}

	t, err := template.ParseFiles("C:/www/templates/search.html", "C:/www/templates/header.html", "C:/www/templates/footer.html")

	t.ExecuteTemplate(w, "search", searchResults)
}

func allArticles(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("C:/www/templates/allArticles.html", "C:/www/templates/header.html", "C:/www/templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	res, err := db.Query("SELECT * FROM articles")

	if err != nil {
		panic(err)
	}

	allPosts := []Article{}

	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Idea, &post.Article)
		if err != nil {
			panic(err)
		}

		allPosts = append(allPosts, post)
	}

	t.ExecuteTemplate(w, "allArticles", allPosts)
}

func handleFunc() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	// rtr.HandleFunc("/contacts", contacts).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/save_article", save_article).Methods("POST")
	rtr.HandleFunc("/post/{id:[0-9]+}", watch_post).Methods("GET")
	rtr.HandleFunc("/search", search).Methods("GET")
	rtr.HandleFunc("/all-articles", allArticles).Methods("GET")
	http.Handle("/", rtr)

	http.ListenAndServe(":8080", nil)
}

func main() {
	handleFunc()
}
