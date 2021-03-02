package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type userData struct {
	gorm.Model
	ID       int    `json:ID`
	Name     string `json:name`
	Password string `json:password`
	Role     string `json:role`
}

type tokenData struct {
	Token string `json:token`
}

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func getdatabase() *gorm.DB {
	host := goDotEnvVariable("host")
	user := goDotEnvVariable("user")
	password := goDotEnvVariable("password")
	dbname := goDotEnvVariable("dbname")
	port := goDotEnvVariable("port")
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
	muxrouter.HandleFunc("/superadmin", superAdminHomePage).Methods("GET")
	muxrouter.HandleFunc("/admin", adminHomePage).Methods("GET")
	muxrouter.HandleFunc("/user", userHomePage).Methods("GET")
	muxrouter.HandleFunc("/", loggingMiddleware(homePage)).Methods("POST")
	// muxrouter.Use(loggingMiddleware)
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
	}{false, "Hello public"})
	w.Write(res)
}

func userHomePage(w http.ResponseWriter, r *http.Request) {
	res, _ := json.Marshal(struct {
		IsError bool
		Msg     string
	}{false, "Hello user"})
	w.Write(res)
}

func superAdminHomePage(w http.ResponseWriter, r *http.Request) {
	res, _ := json.Marshal(struct {
		IsError bool
		Msg     string
	}{false, "Hello superadmin"})
	w.Write(res)
}

func adminHomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("add")
	res, _ := json.Marshal(struct {
		IsError bool
		Msg     string
	}{false, "Hello admin"})
	w.Write(res)
}

func loggingMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		bodydata, _ := ioutil.ReadAll(r.Body)
		var tokendata tokenData
		_ = json.Unmarshal(bodydata, &tokendata) //err:=json.NewDecoder(r.body).Decode(&tokendata)
		if tokendata.Token == "" {
			handler.ServeHTTP(w, r)
			return
		}
		// cookie, _ := r.Cookie("Token")
		// tokendata.Token = cookie.Value
		token, err := jwt.Parse(string(tokendata.Token), func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			myKeyString := goDotEnvVariable("myKey")
			myKey := []byte(myKeyString)
			return myKey, nil
		})

		if err != nil {
			fmt.Fprintf(w, err.Error())
		}

		// fmt.Println(r.RequestURI, token.Claims, token.Header, token.Method, token.Raw, token.Signature, token.Valid)
		if token.Valid {
			role := token.Claims.(jwt.MapClaims)["Role"].(string)
			if role == "SuperAdmin" {
				http.Redirect(w, r, "/superadmin", http.StatusSeeOther)
				return
			} else if role == "Admin" {
				http.Redirect(w, r, "/admin", http.StatusSeeOther)
				return
			} else if role == "User" {
				http.Redirect(w, r, "/user", http.StatusSeeOther)
				return
			}
			// next.ServeHTTP(w, r)
		} else {
			fmt.Println("not logged in")
		}
	}
}

func fetchUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := getdatabase()
	defer close(db)

	bodydata, _ := ioutil.ReadAll(r.Body)
	var user userData
	_ = json.Unmarshal(bodydata, &user)
	fmt.Println(user)
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
			claims["Role"] = fetchUser.Role
			myKeyString := goDotEnvVariable("myKey")
			myKey := []byte(myKeyString)
			tokenstring, err := token.SignedString(myKey)

			if err != nil {
				fmt.Println(err)
			}
			// http.SetCookie(w, &http.Cookie{
			// 	Name:  "Token",
			// 	Value: tokenstring,
			// })
			res, _ = json.Marshal(struct {
				IsError bool
				Msg     string
				Count   int
				Data    userData
				Token   string
			}{false, "Signed In Successfully", 1, fetchUser, tokenstring})
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
