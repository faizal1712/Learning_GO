package main

import (
	"log"
	"time"
)

var ch = make(chan bool)

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func apiCall(i int, ch chan bool) {
	log.Println("API call for", i, "started")
	time.Sleep(100 * time.Millisecond)
	ch <- true
}

func main() {
	// concurrencyLimit := 1001
	numArray := makeRange(0, 1000)

	start := time.Now()

	for i := range numArray {
		go apiCall(i, ch)
	}
	for i := range numArray {
		_ = i
		<-ch
	}
	elapsed := time.Since(start)
	log.Printf("Time taken %s", elapsed)
}
