package main

import (
	"fmt"
	"time"
)

// go routine main
func main() {
	fmt.Println("main func")
	ch1 := make(chan bool)
	ch2 := make(chan bool)

	now := time.Now()

	go wait10sec(ch2)
	go wait5sec(ch1)
	<-ch2
	<-ch1

	fmt.Println(time.Since(now))
}

// seperate go routine
func wait5sec(ch chan bool) {
	time.Sleep(time.Second * 5)
	fmt.Println("5 sec thya gai")
	ch <- true
}

// seperate go routine
func wait10sec(ch chan bool) {
	time.Sleep(time.Second * 10)
	fmt.Println("10 sec thya gai")
	ch <- true
}
