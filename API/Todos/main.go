package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type todos struct {
	Code int `json:"code"`
	Meta struct {
		Pagination struct {
			Total int `json:"total"`
			Pages int `json:"pages"`
			Page  int `json:"page"`
			Limit int `json:"limit"`
		} `json:"pagination"`
	} `json:"meta"`
	Data []struct {
		ID        int       `json:"id"`
		UserID    int       `json:"user_id"`
		Title     string    `json:"title"`
		Completed bool      `json:"completed"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"data"`
}

func main() {
	url := "https://gorest.co.in/public-api/todos"
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	} else {
		resbyteValue, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		} else {
			var data todos
			json.Unmarshal(resbyteValue, &data)
			fmt.Println(data)
			file, err := json.MarshalIndent(data, "", "")
			err = ioutil.WriteFile("todos.json", file, 0644)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
