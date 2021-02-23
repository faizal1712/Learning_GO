package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type fileData struct {
	gorm.Model
	Path string `json:"Path"`
}

func getdatabase() *gorm.DB {
	conn, err := gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/go_images_mysql?parseTime=true"), &gorm.Config{})
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
	db.AutoMigrate(&fileData{})
	defer close(db)
}

func close(conn *gorm.DB) {
	sqldb, _ := conn.DB()
	sqldb.Close()
}

func addImages(w http.ResponseWriter, r *http.Request) {
	db := getdatabase()
	r.ParseMultipartForm(10 << 20)
	file, filehandle, err := r.FormFile("imageFile")
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/images", 301)
	}
	fmt.Println(filehandle.Filename)

	imgPath := path.Join("images/", filehandle.Filename)
	fmt.Println(imgPath)
	destination, err := os.Create(imgPath)
	if err != nil {
		panic(err)
	}

	defer destination.Close()
	io.Copy(destination, file)
	var fileRow fileData
	fileRow.Path = imgPath
	insert := db.Create(&fileRow)
	if insert.Error != nil {
		panic(insert.Error)
	}
	http.Redirect(w, r, "/images", 301)
}

func viewImages(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	db := getdatabase()
	var files []fileData
	db.Find(&files)
	tmpl.Execute(w, files)
}

func main() {
	fmt.Println()
	migrate()
	var muxrouter *mux.Router = mux.NewRouter()

	muxrouter.Handle("/images/{rest}", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))

	muxrouter.HandleFunc("/images", addImages).Methods("POST")

	muxrouter.HandleFunc("/images", viewImages).Methods("GET")
	http.ListenAndServe(":8080", muxrouter)
}
