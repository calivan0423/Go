package main

import (
	"fmt"
	"time"
)

func main() {

	channel := make(chan string)

	people := [5]string{"calivan", "lee", "jang", "choi", "kim"}
	for _, person := range people {
		go isFun(person, channel)
	}
	for i := 0; i < len(people); i++ {
		fmt.Println("waiting for", i)
		fmt.Println(<-channel)
	}

}

func isFun(person string, c chan string) {
	time.Sleep(time.Second * 3)
	c <- person + " is fun"
}
