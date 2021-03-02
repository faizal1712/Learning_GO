package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "faizu1712"
	dbname   = "go_postgre"
)

type userStr struct {
	id    int
	name  string
	food  string
	sport string
}

func dbConn() (db *sql.DB) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
	return db
}

func ctreateTable() {
	//CREATE
	db := dbConn()
	defer db.Close()
	query := `create table users(id serial primary key,name varchar(20) not null,food varchar(10) not null,sport varchar(20) not null)`
	_, err := db.Exec(query)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Successfully created!")
}

func insertTab() {
	//Insert
	db := dbConn()
	query := `insert into users(name,food,sport) values($1,$2,$3)`
	_, err := db.Exec(query, "Jhanvi", "Lasagnia", "Secret")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Successfully inserted!")
}

func selectTab() {
	//select
	db := dbConn()
	query := `select * from users`
	rows, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}

	var users []userStr
	for rows.Next() {
		var u userStr
		err = rows.Scan(&u.id, &u.name, &u.food, &u.sport)
		if err != nil {
			panic(err.Error())
		}
		users = append(users, u)
	}
	fmt.Println(users)
}

func selectRow() {
	//select one
	db := dbConn()
	query := `select * from users where id=$1`
	var u userStr
	err := db.QueryRow(query, 2).Scan(&u.id, &u.name, &u.food, &u.sport)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(u)
}

func updateTab() {
	//update
	db := dbConn()
	query := `update users set food=$1 where id=$2`
	_, err := db.Exec(query, "sabChalega", 2)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Successfully updated!")
}

func deleteTab() {
	//delete
	db := dbConn()
	query := `delete from users where id=$1`
	result, err := db.Exec(query, 1)
	if err != nil {
		panic(err.Error())
	}
	count, _ := result.RowsAffected()
	fmt.Println("Successfully deleted!", count)
}

func main() {

}
