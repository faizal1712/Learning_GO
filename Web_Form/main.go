package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type scoreData struct {
	Name  string `json:"Name"`
	GRE   int    `json:"GRE"`
	TOEFL int    `json:"TOEFL"`
	CGPA  int    `json:"CGPA"`
}

func main() {
	formTemplate := template.Must(template.ParseFiles("form.html")) // or temp, err := template.ParseFiles("view.html"); temp.Execute(...)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			formTemplate.Execute(w, nil)
			return
		}
		GREscore, _ := strconv.Atoi(r.FormValue("GRE_score"))
		TOEFLscore, _ := strconv.Atoi(r.FormValue("TOEFL_score"))
		CGPAscore, _ := strconv.Atoi(r.FormValue("CGPA_score"))
		scoreDetails := scoreData{
			Name:  r.FormValue("user_name"),
			GRE:   GREscore,
			TOEFL: TOEFLscore,
			CGPA:  CGPAscore,
		}

		formDataJSONFile, err := os.Open("formData.json")
		if err != nil {
			fmt.Println(err)
		}
		defer formDataJSONFile.Close()
		formDataByteValue, _ := ioutil.ReadAll(formDataJSONFile)
		var formDataArray []scoreData
		json.Unmarshal(formDataByteValue, &formDataArray)

		formDataArray = append(formDataArray, scoreDetails)

		// f, err := os.OpenFile("formData.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		// if err != nil {
		// 	fmt.Println(err)
		// }

		data, err := json.MarshalIndent(formDataArray, "", "")

		if err = ioutil.WriteFile("formData.json", data, 0644); err != nil {
			fmt.Println(err)
		}

		formTemplate.Execute(w, struct{ Success bool }{true})
	})

	http.ListenAndServe(":8080", nil)

}
