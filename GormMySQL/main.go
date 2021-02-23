package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/gorilla/mux"
)

type scoreData struct {
	gorm.Model
	ID    int    `json:ID`
	Name  string `json:"Name"`
	GRE   int    `json:"GRE"`
	TOEFL int    `json:"TOEFL"`
	CGPA  int    `json:"CGPA"`
}

func getdatabase() *gorm.DB {
	conn, err := gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/gorm_mysql?parseTime=true"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqldb, err := conn.DB()
	if err != nil {
		panic(err)
	}
	if err = sqldb.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("connected to database")
	return conn
}

func migrate() {
	db := getdatabase()
	db.AutoMigrate(&scoreData{})
	defer close(db)
}

func close(conn *gorm.DB) {
	sqldb, _ := conn.DB()
	sqldb.Close()
}

func router() {
	var muxrouter *mux.Router = mux.NewRouter()
	muxrouter.HandleFunc("/user", getAllUser).Methods("GET")
	muxrouter.HandleFunc("/user", addUser).Methods("POST")
	muxrouter.HandleFunc("/user/{id}", getOneUser).Methods("GET")
	muxrouter.HandleFunc("/user/{id}", updateUser).Methods("POST")
	muxrouter.HandleFunc("/user/{id}", deleteUser).Methods("DELETE")
	http.ListenAndServe(":8080", muxrouter)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := getdatabase()
	defer close(db)

	bodydata, _ := ioutil.ReadAll(r.Body)
	var score scoreData
	_ = json.Unmarshal(bodydata, &score)

	insert := db.Create(&score)
	if insert.Error != nil {
		panic(insert.Error)
	}

	res, _ := json.Marshal(struct {
		Is_error bool
		Msg      string
	}{false, "User Added Successfully"})
	w.Write(res)
}

func getAllUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := getdatabase()
	defer close(db)
	var scoreDataArray []scoreData

	db.Find(&scoreDataArray)
	// fmt.Println(scoreDataArray)
	res, _ := json.Marshal(scoreDataArray)
	w.Write(res)
}

func getOneUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := getdatabase()
	defer close(db)
	id := mux.Vars(r)["id"]
	var score scoreData
	db.Where("id = ?", id).Find(&score)

	res, _ := json.Marshal(score)
	w.Write(res)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := getdatabase()
	defer close(db)
	id := mux.Vars(r)["id"]
	var score scoreData
	db.First(&score, id)
	bodydata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bodydata, &score)
	if err != nil {
		panic(err)
	}
	db.Save(&score)
	res, _ := json.Marshal(struct {
		Is_error bool
		Msg      string
	}{false, "User Updated Successfully"})
	w.Write(res)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := getdatabase()
	defer close(db)
	id := mux.Vars(r)["id"]
	var score scoreData
	db.Delete(&score, id)
	res, _ := json.Marshal(struct {
		Is_error bool
		Msg      string
	}{false, "User Deleted Successfully"})
	w.Write(res)
}

func main() {
	migrate()
	router()
}
