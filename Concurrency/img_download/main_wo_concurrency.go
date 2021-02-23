package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	fileName1 := "sample1.jpg"
	fileName2 := "sample2.jpg"
	fileName3 := "sample3.jpg"
	fileName4 := "sample4.jpg"
	fileName5 := "sample5.jpg"
	fileName6 := "sample6.jpg"
	fileName7 := "sample7.jpg"
	fileName8 := "sample8.jpg"
	fileName9 := "sample9.jpg"
	fileName10 := "sample10.jpg"
	fileName11 := "sample11.jpg"
	fileName12 := "sample12.jpg"
	fileName13 := "sample13.jpg"
	fileName14 := "sample14.jpg"
	fileName15 := "sample15.jpg"
	URL1 := "https://i.pinimg.com/originals/af/8d/63/af8d63a477078732b79ff9d9fc60873f.jpg"
	URL2 := "https://images.pexels.com/photos/1591447/pexels-photo-1591447.jpeg?auto=compress&cs=tinysrgb&dpr=1&w=500"
	URL3 := "https://i.pinimg.com/originals/df/07/cb/df07cb4ccb697303462ad7a8b57b852f.jpg"
	URL4 := "https://images.pexels.com/photos/1563356/pexels-photo-1563356.jpeg?auto=compress&cs=tinysrgb&dpr=1&w=500"
	URL5 := "https://i.pinimg.com/736x/37/6e/2d/376e2dab5652d6e1751e25cbcb52f2d5.jpg"
	URL6 := "https://images.unsplash.com/photo-1513151233558-d860c5398176?ixid=MXwxMjA3fDB8MHxzZWFyY2h8M3x8ZnVuJTIwYmFja2dyb3VuZHxlbnwwfHwwfA%3D%3D&ixlib=rb-1.2.1&w=1000&q=80"
	URL7 := "https://i.pinimg.com/originals/c8/2a/f9/c82af9c8a818d8dba545fb896b8a6b2c.jpg"
	URL8 := "https://images.pexels.com/photos/1420440/pexels-photo-1420440.jpeg?auto=compress&cs=tinysrgb&dpr=1&w=500"
	URL9 := "https://wallpapercave.com/wp/wp2462597.jpg"
	URL10 := "https://cdn.pixabay.com/photo/2016/11/29/05/45/astronomy-1867616__340.jpg"
	URL11 := "https://venngage-wordpress.s3.amazonaws.com/uploads/2018/09/Perfect-Sunset-Nature-Background-Image.jpeg"
	URL12 := "https://i.pinimg.com/originals/d5/c8/7c/d5c87c9160550d386791069339bbd762.jpg"
	URL13 := "https://images.pexels.com/photos/235986/pexels-photo-235986.jpeg?auto=compress&cs=tinysrgb&dpr=1&w=500"
	URL14 := "https://cdn.pixabay.com/photo/2015/04/23/22/00/tree-736885__340.jpg"
	URL15 := "https://images.pexels.com/photos/255379/pexels-photo-255379.jpeg?auto=compress&cs=tinysrgb&dpr=1&w=500"

	now := time.Now()

	err := downloadFile(URL1, fileName1)
	if err != nil {
		log.Fatal(err)
	}
	err = downloadFile(URL2, fileName2)
	if err != nil {
		log.Fatal(err)
	}
	err = downloadFile(URL3, fileName3)
	if err != nil {
		log.Fatal(err)
	}
	err = downloadFile(URL4, fileName4)
	if err != nil {
		log.Fatal(err)
	}
	err = downloadFile(URL5, fileName5)
	if err != nil {
		log.Fatal(err)
	}
	err = downloadFile(URL6, fileName6)
	if err != nil {
		log.Fatal(err)
	}
	err = downloadFile(URL7, fileName7)
	if err != nil {
		log.Fatal(err)
	}
	err = downloadFile(URL8, fileName8)
	if err != nil {
		log.Fatal(err)
	}
	err = downloadFile(URL9, fileName9)
	if err != nil {
		log.Fatal(err)
	}
	err = downloadFile(URL10, fileName10)
	if err != nil {
		log.Fatal(err)
	}
	err = downloadFile(URL11, fileName11)
	if err != nil {
		log.Fatal(err)
	}
	err = downloadFile(URL12, fileName12)
	if err != nil {
		log.Fatal(err)
	}
	err = downloadFile(URL13, fileName13)
	if err != nil {
		log.Fatal(err)
	}
	err = downloadFile(URL14, fileName14)
	if err != nil {
		log.Fatal(err)
	}
	err = downloadFile(URL15, fileName15)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(time.Since(now))

	// fmt.Printf("File %s download in current working directory", fileName)
}

func downloadFile(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}