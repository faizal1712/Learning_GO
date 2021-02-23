package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func getAllUser(w http.ResponseWriter, r *http.Request) {
	db := conn()
	defer db.Close()
	fetchAll, err := db.Query("SELECT * FROM user_details")

	if err != nil {
		panic(err.Error())
	}
	var scoreDataArray []scoreData
	for fetchAll.Next() {
		var score scoreData
		err = fetchAll.Scan(&score.Id, &score.Name, &score.GRE, &score.TOEFL, &score.CGPA)
		scoreDataArray = append(scoreDataArray, score)
	}
	res, _ := json.Marshal(scoreDataArray)
	w.Write(res)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	db := conn()
	defer db.Close()
	bodydata, err := ioutil.ReadAll(r.Body)

	var score scoreData

	_ = json.Unmarshal(bodydata, &score)

	insert, err := db.Query("INSERT INTO user_details (name, GRE, TOEFL, CGPA) VALUES (?,?,?,?)", score.Name, score.GRE, score.TOEFL, score.CGPA)

	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()

	data := struct {
		Is_error bool
		Msg      string
	}{false, "User Added Successfully"}
	res, _ := json.Marshal(data)
	w.Write(res)
}

func getOneUser(w http.ResponseWriter, r *http.Request) {
	db := conn()
	defer db.Close()
	id, _ := strconv.Atoi(r.FormValue("Id"))
	var score scoreData
	err := db.QueryRow("SELECT * FROM user_details WHERE id = ?", id).Scan(&score.Id, &score.Name, &score.GRE, &score.TOEFL, &score.CGPA)
	if err != nil {
		panic(err.Error())
	}
	res, _ := json.Marshal(score)
	w.Write(res)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	db := conn()
	defer db.Close()

	GREscore, _ := strconv.Atoi(r.FormValue("GRE_score"))
	TOEFLscore, _ := strconv.Atoi(r.FormValue("TOEFL_score"))
	CGPAscore, _ := strconv.Atoi(r.FormValue("CGPA_score"))
	idscore, _ := strconv.Atoi(r.FormValue("ID"))
	update, err := db.Query("UPDATE user_details SET name = ?, GRE = ?, TOEFL = ?, CGPA = ? WHERE id = ?", r.FormValue("user_name"), GREscore, TOEFLscore, CGPAscore, idscore)

	if err != nil {
		panic(err.Error())
	}
	defer update.Close()
	res, _ := json.Marshal(struct {
		Is_error bool
		Msg      string
	}{false, "User Updated Successfully"})
	w.Write(res)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	db := conn()
	defer db.Close()

	id := r.URL.Query().Get("Id")
	delete, err := db.Query("DELETE FROM user_details WHERE id = ?", id)
	if err != nil {
		panic(err.Error())
	}
	defer delete.Close()
	res, _ := json.Marshal(struct {
		Is_error bool
		Msg      string
	}{false, "User Deleted Successfully"})
	w.Write(res)
}

func main() {
	fmt.Println("")
	http.HandleFunc("/user", getAllUser)
	http.HandleFunc("/addUser", addUser)
	http.HandleFunc("/fetchUser", getOneUser)
	http.HandleFunc("/updateUser", updateUser)
	http.HandleFunc("/deleteUser", deleteUser)

	http.ListenAndServe(":8080", nil)
}
