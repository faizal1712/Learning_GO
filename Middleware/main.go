package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type userData struct {
	gorm.Model
	ID       int    `json:ID`
	Name     string `json:name`
	Password string `json:password`
}

type tokenData struct {
	Token string `json:token`
}

const (
	host     = "localhost"
	port     = "5432"
	user     = "postgres"
	password = "faizu1712"
	dbname   = "project1"
)

var myKey = []byte("faizu17121999")

func getdatabase() *gorm.DB {
	var dsn string = "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable" // or "root:@tcp(127.0.0.1:3306)/gorm_mysql?parseTime=true"
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
	conn := getdatabase()
	conn.AutoMigrate(&userData{})
	defer close(conn)
}

func close(conn *gorm.DB) {
	psqldb, _ := conn.DB()
	psqldb.Close()
}

func router() {
	var muxrouter *mux.Router = mux.NewRouter()
	muxrouter.HandleFunc("/signin", fetchUser).Methods("POST")
	muxrouter.HandleFunc("/signup", addUser).Methods("POST")
	muxrouter.HandleFunc("/", homePage).Methods("POST")
	muxrouter.Use(loggingMiddleware)
	http.ListenAndServe(":8080", muxrouter)
}

func main() {
	migrate()
	router()
}

func homePage(w http.ResponseWriter, r *http.Request) {
	res, _ := json.Marshal(struct {
		IsError bool
		Msg     string
	}{false, "Hello"})
	w.Write(res)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		bodydata, _ := ioutil.ReadAll(r.Body)
		var tokendata tokenData
		_ = json.Unmarshal(bodydata, &tokendata)

		token, err := jwt.Parse(string(tokendata.Token), func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return myKey, nil
		})

		if err != nil {
			fmt.Fprintf(w, err.Error())
		}

		fmt.Println(r.RequestURI)
		if token.Valid {
			// Call the next handler, which can be another middleware in the chain, or the final handler.
			next.ServeHTTP(w, r)
		} else {
			fmt.Println("not logged in")
		}

		// Do stuff here
	})
}

func fetchUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := getdatabase()
	defer close(db)

	bodydata, _ := ioutil.ReadAll(r.Body)
	var user userData
	_ = json.Unmarshal(bodydata, &user)

	var fetchUser userData

	var res []byte

	result := db.Where("name = ? ", user.Name).Find(&fetchUser).Select("name")
	if result.Error != nil {
		res, _ = json.Marshal(struct {
			IsError bool
			Msg     error
		}{true, result.Error})
	} else {
		err := bcrypt.CompareHashAndPassword([]byte(fetchUser.Password), []byte(user.Password))

		if err != nil {
			fmt.Println("Pwd don't match")
			res, _ = json.Marshal(struct {
				IsError bool
				Msg     string
			}{true, "Password Didn't match"})
		} else {
			token := jwt.New(jwt.SigningMethodHS256)
			// var claimData tokenData
			// claimData.Name = fetchUser.Name
			// claimData.ID = fetchUser.ID
			// marshal := json.Marshal(claimData)
			// var claims map[string]interface{}
			// if err := json.Unmarshal(marshal, &claims); err != nil {
			// 	fmt.Println("Couldn't parse claims JSON: %v", err)
			// }
			claims := token.Claims.(jwt.MapClaims)
			claims["Name"] = fetchUser.Name
			claims["Id"] = fetchUser.ID

			tokenstring, err := token.SignedString(myKey)

			if err != nil {
				fmt.Println(err)
			}
			// r.Header.Set("Token", tokenstring)
			res, _ = json.Marshal(struct {
				IsError bool
				Count   int
				Data    userData
				Token   string
			}{false, 1, fetchUser, tokenstring})
		}
	}
	w.Write(res)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := getdatabase()
	defer close(db)

	bodydata, _ := ioutil.ReadAll(r.Body)
	var user userData
	_ = json.Unmarshal(bodydata, &user)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(passwordHash)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(user)

	insert := db.Create(&user)
	if insert.Error != nil {
		panic(insert.Error)
	}

	res, _ := json.Marshal(struct {
		IsError bool
		Msg     string
	}{false, "User Added Successfully"})
	w.Write(res)
}
