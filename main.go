package main

import (
	"fmt"
	"net/http"
	"html/template"
	"github.com/gorilla/mux"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Article struct {
	Id uint16
	Title string
	Anons string
	FullText string
}
var posts = []Article{}
var showPost = Article{}

func index(w http.ResponseWriter, r *http.Request){
	//парсим темплейты индекс хедер и футер
	 t, err := template.ParseFiles("templates/index.html","templates/header.html","templates/footer.html")
	//проверка на ошибку
	 if err != nil {
		fmt.Fprintf(w, err.Error())
	 }

	 		//соединение с сервером
		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golang")
		if err != nil {
			http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		res, err := db.Query("SELECT * FROM `article`")
		if err != nil {
			panic(err)
		}
		posts = []Article{}
		for res.Next() {
			var post Article
			err = res.Scan(&post.Id,&post.Title,&post.Anons,&post.FullText,)
			if err != nil {
				panic(err)
			}
			posts = append(posts, post)
			fmt.Println(fmt.Sprintf("Post: %s with id: %d", post.Title, post.Id))
		}


	 // вывод tтмлейта на страницу
	 t.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, r *http.Request){
	t, err := template.ParseFiles("templates/create.html","templates/header.html","templates/footer.html")

	if err != nil {
	   fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "create", nil)
}


func save_article(w http.ResponseWriter, r *http.Request) {
	//сохранение в переменные данных которые мы получаем из формы (все имена указаны в тегах)
	title := r.FormValue("title")
	anons := r.FormValue("anons")
 	full_text := r.FormValue("full_text")
	if title == "" || anons == "" || full_text == "" {
		fmt.Fprintf(w, "Не все данные заполнены!!! ")
	}else {
		//соединение с сервером
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golang")
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	//утановка данных в базу
	insert, err := db.Query(fmt.Sprintf("INSERT INTO `article` (`title`, `anons`,`full_text`) VALUES('%s', '%s', '%s')", title, anons, full_text))
	if err != nil {
		http.Error(w, "Ошибка при выполнении запроса к базе данных", http.StatusInternalServerError)
		return
	}
	defer insert.Close()

	//перенаправляем на основную страницу
	http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
func show_post(w http.ResponseWriter, r *http.Request){
 vars := mux.Vars(r)

 t, err := template.ParseFiles("templates/show.html","templates/header.html","templates/footer.html")
	//проверка на ошибку
	 if err != nil {
		fmt.Fprintf(w, err.Error())
	 }

 db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golang")
 if err != nil {
	 http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
	 return
 }
 defer db.Close()

 res, err := db.Query(fmt.Sprintf("SELECT * FROM `article` WHERE `id` = '%s'", vars["id"]))
		if err != nil {
			panic(err)
		}
		showPost = Article{}
		for res.Next() {
			var post Article
			err = res.Scan(&post.Id,&post.Title,&post.Anons,&post.FullText,)
			if err != nil {
				panic(err)
			}
			showPost = post
			fmt.Println(fmt.Sprintf("Post: %s with id: %d", post.Title, post.Id))
		}


	 // вывод tтмлейта на страницу
	 t.ExecuteTemplate(w, "show", showPost)


}


func handleFunc(){
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create/", create).Methods("GET")
	rtr.HandleFunc("/save_article/", save_article).Methods("POST")
	rtr.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")

	http.Handle("/", rtr)
	http.ListenAndServe(":8080", nil)
}

func main()  {
handleFunc()
}