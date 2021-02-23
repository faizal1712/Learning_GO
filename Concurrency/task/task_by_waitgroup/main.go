package main

import (
	"log"
	"sync"
	"time"
)

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func apiCall(i int, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("API call for", i, "started")
	time.Sleep(100 * time.Millisecond)
}

func main() {
	var wg sync.WaitGroup
	numArray := makeRange(0, 1000)

	start := time.Now()

	for i := range numArray {
		wg.Add(1)
		go apiCall(i, &wg)
	}

	wg.Wait()
	elapsed := time.Since(start)
	log.Printf("Time taken %s", elapsed)
}
