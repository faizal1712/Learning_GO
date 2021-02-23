package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type scoreData struct {
	Id    int    `json:Id`
	Name  string `json:"Name"`
	GRE   int    `json:"GRE"`
	TOEFL int    `json:"TOEFL"`
	CGPA  int    `json:"CGPA"`
}

func conn() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/go_mysql")

	if err != nil || db.Ping() != nil {
		panic(err.Error())
	}

	return db
}

func add(x, y int) int {
	return x + y
}

func viewAll(w http.ResponseWriter, r *http.Request) {
	funcs := template.FuncMap{"add": add}
	temp, err := template.New("view.html").Funcs(funcs).ParseFiles("view.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("view before")

	db := conn()

	fetchAll, err := db.Query("SELECT * FROM user_details")

	if err != nil {
		panic(err.Error())
	}
	var scoreDataArray []scoreData
	fmt.Println("view after")
	for fetchAll.Next() {
		var score scoreData
		err = fetchAll.Scan(&score.Id, &score.Name, &score.GRE, &score.TOEFL, &score.CGPA)
		scoreDataArray = append(scoreDataArray, score)
	}
	temp.Execute(w, scoreDataArray)
}

func insert(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("insert.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("insert before")

	if r.Method != http.MethodPost {
		temp.Execute(w, nil)
		return
	}

	db := conn()
	defer db.Close()
	GREscore, _ := strconv.Atoi(r.FormValue("GRE_score"))
	TOEFLscore, _ := strconv.Atoi(r.FormValue("TOEFL_score"))
	CGPAscore, _ := strconv.Atoi(r.FormValue("CGPA_score"))
	insert, err := db.Query("INSERT INTO user_details (name, GRE, TOEFL, CGPA) VALUES (?,?,?,?)", r.FormValue("user_name"), GREscore, TOEFLscore, CGPAscore)

	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
	fmt.Println("insert after")
	http.Redirect(w, r, "/", 301)
}

func update(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("Id")
	temp, err := template.ParseFiles("update.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("update before")
	db := conn()
	defer db.Close()
	if r.Method != http.MethodPost {
		var rowData scoreData
		if err = db.QueryRow("SELECT * FROM user_details WHERE id = "+id).Scan(&rowData.Id, &rowData.Name, &rowData.GRE, &rowData.TOEFL, &rowData.CGPA); err != nil {
			panic(err.Error())
		}
		temp.Execute(w, rowData)
		return
	}

	GREscore, _ := strconv.Atoi(r.FormValue("GRE_score"))
	TOEFLscore, _ := strconv.Atoi(r.FormValue("TOEFL_score"))
	CGPAscore, _ := strconv.Atoi(r.FormValue("CGPA_score"))
	idscore, _ := strconv.Atoi(r.FormValue("ID"))
	update, err := db.Query("UPDATE user_details SET name = ?, GRE = ?, TOEFL = ?, CGPA = ? WHERE id = ?", r.FormValue("user_name"), GREscore, TOEFLscore, CGPAscore, idscore)

	if err != nil {
		panic(err.Error())
	}
	defer update.Close()
	fmt.Println("update after")
	http.Redirect(w, r, "/", 301)
}

func delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("Id")
	fmt.Println("delete before")
	db := conn()
	defer db.Close()
	delete, err := db.Query("DELETE FROM user_details WHERE id = ?", id)

	/*
		delete, err: = db.Prepare("DELETE FROM user_details WHERE id = ?")
		delete.Exec(id)
	*/

	if err != nil {
		panic(err.Error())
	}
	defer delete.Close()
	fmt.Println("delete after")
	http.Redirect(w, r, "/", 301)
}

func main() {
	http.HandleFunc("/", viewAll) //insert; insert/{id}
	http.HandleFunc("/insert", insert)
	http.HandleFunc("/update", update)
	http.HandleFunc("/delete", delete)
	http.ListenAndServe(":8080", nil)
}
