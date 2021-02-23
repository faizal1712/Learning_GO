package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type scoreData struct {
	Name  string `json:"Name"`
	GRE   int    `json:"GRE"`
	TOEFL int    `json:"TOEFL"`
	CGPA  int    `json:"CGPA"`
}

func main() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/go_mysql")
	// fmt.Println(err)
	if err != nil || db.Ping() != nil {
		panic(err.Error())
	}
	defer db.Close()

	// insert, err := db.Query("INSERT INTO user_details VALUES (1, 'Faizu', 100, 5, 2)")

	// insert, err := db.Query("INSERT INTO user_details (name, GRE, TOEFL, CGPA) VALUES ('Zuhi', 250, 100, 9)")

	// insert, err := db.Query("INSERT INTO user_details (name, GRE, TOEFL, CGPA) VALUES ('Jeet', 200, 50, 10)")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// defer insert.Close()

	// update, err := db.Query("UPDATE user_details SET GRE = 300 WHERE name = 'Zuhi'")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// defer update.Close()

	// delete, err := db.Query("DELETE FROM user_details WHERE name = 'Jeet'")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// defer delete.Close()

	view, err := db.Query("SELECT * FROM user_details")
	var scoreDataArray []scoreData
	for view.Next() {
		var rowData scoreData
		var id int
		err := view.Scan(&id, &rowData.Name, &rowData.GRE, &rowData.TOEFL, &rowData.CGPA)
		if err != nil {
			panic(err.Error())
		}
		scoreDataArray = append(scoreDataArray, rowData)
		// fmt.Println(rowData)
	}
	fmt.Println(scoreDataArray)

	// var rowData scoreData
	// var id int
	// if err = db.QueryRow("SELECT * FROM user_details WHERE id = 1").Scan(&id, &rowData.Name, &rowData.GRE, &rowData.TOEFL, &rowData.CGPA); err != nil {
	// 	panic(err.Error())
	// }
	// fmt.Println(rowData)
}
