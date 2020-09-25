package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

var db *sql.DB
var err error

type Url struct{
	LUrl string `json:"lurl"`
	SUrl string `json:"surl"`
	ID   int    `json:"id"`
}
//Для быстрого кодирования строки используется алгоритм base64
func EncodeB64(message string) string {
	base64Text := make([]byte, base64.StdEncoding.EncodedLen(len(message)))
	base64.StdEncoding.Encode(base64Text, []byte(message))
	return string(base64Text)
}
//Так как в некоторых случаях кодирование приводит к появлению ==
//а в этом алгоритме == ипользуется только для дополнения пустых битов
//решено было заменить "=" на ""
func getShortUrl(url *Url){
	url.SUrl = EncodeB64(string(url.ID))
	url.SUrl = strings.Replace(url.SUrl,"=","",1000)
}

/*
	Метод проверяет, обращались ли раньше для сокращения этой ссылки
	Если нет, то создает ее экземпляр, добавляет в БД и отправляет
	json ответ пользователю
 */
func addUrl(w http.ResponseWriter, r *http.Request){
	var isInUrls bool
	w.Header().Set("Content-Type", "application/json")

	var url Url
	_ = json.NewDecoder(r.Body).Decode(&url)
	rows, err := db.Query("SELECT surl FROM avito_urls.urls WHERE lurl = ?" , url.LUrl)

	var str string
	for rows.Next(){
		rows.Scan(&str)
	}
	defer rows.Close()

	if err != nil{
		panic(err)
	}
	isInUrls = (str != "")

	if (!isInUrls) {
		var count int
		cnt, err := db.Query("SELECT MAX(ID) FROM avito_urls.urls")

		for cnt.Next(){
			cnt.Scan(&count)
		}
		defer cnt.Close()

		if err != nil{
			panic(err)
		}
		url.ID = count + 1
		isInUrls = true

		if (url.SUrl == "") {
			getShortUrl(&url)
		}

		_, err = db.Exec("insert into avito_urls.urls (surl, lurl) values (?, ?)", url.SUrl, url.LUrl)
	} else{
		url.SUrl = str
	}
	_ = json.NewEncoder(w).Encode(url)
}
//Функция поиска длинной ссылки по короткой в БД для переадресации
func findUrl(Surl string) string {
	rows, err := db.Query("SELECT lurl FROM avito_urls.urls WHERE surl = ?" , Surl)
	var str string
	for rows.Next(){
		rows.Scan(&str)
	}
	defer rows.Close()
	if err != nil{
		panic(err)
	}
	return str
}
//Фукнция переадресации
func redirect(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var targetLUrl string = findUrl(params["str"])
	http.Redirect(w, r, targetLUrl, 301)
}
/*
	Сервер обрабатывает несколько http запросов
	1) http://localhost:8000/url   POST - json запрос на получение короткой ссылки
	2) http://localhost:8000/сurl  POST - json запрос на добавление кастомизированной ссылки
	3) http://localhost:8000/{url} GET	- запрос на переадресацию по короткой ссылке
	Все данные хранятся на БД MySQL
 */
func main() {
	Db, err := sql.Open("mysql", "root:Root2000@/avito_urls")
	db = Db
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err != nil{
		panic(err)
	}

	fmt.Println("Server is running...")
	r := mux.NewRouter()
	r.HandleFunc("/url", addUrl).Methods("POST")
	r.HandleFunc("/curl", addUrl).Methods("POST")
	r.HandleFunc("/{str}", redirect).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", r))
}
