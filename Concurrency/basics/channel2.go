package main

import (
	"fmt"
	"time"
)

// go routine main
func main() {
	fmt.Println("main func")
	concurrencyLimit := 5
	ch := make(chan bool, concurrencyLimit)

	now := time.Now()

	for i := 0; i < concurrencyLimit; i++ {
		go waitSec(ch, 1)
	}

	for i := 0; i < concurrencyLimit; i++ {
		<-ch
	}
	fmt.Println(time.Since(now))
}

func waitSec(ch chan bool, sec time.Duration) {
	time.Sleep(time.Second * sec)
	fmt.Print(sec)
	fmt.Print(" done!!")
	fmt.Println()
	ch <- true
}
