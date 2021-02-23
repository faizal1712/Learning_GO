package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type productCategory struct {
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
		ID         int `json:"id"`
		ProductID  int `json:"product_id"`
		CategoryID int `json:"category_id"`
	} `json:"data"`
}

func main() {
	url := "https://gorest.co.in/public-api/product-categories"
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	} else {
		resbyteValue, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		} else {
			var data productCategory
			json.Unmarshal(resbyteValue, &data)
			fmt.Println(data)
			file, err := json.MarshalIndent(data, "", "")
			err = ioutil.WriteFile("productCategory.json", file, 0644)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
