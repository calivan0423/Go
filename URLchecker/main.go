package main

import (
	"errors"
	"fmt"
	"net/http"
)

type Requestresult struct {
	url    string
	status string
}

var errRequestFailed = errors.New("Request failed")

func main() {
	results := make(map[string]string)
	channel := make(chan Requestresult)

	urls := []string{
		"https://www.google.co.kr/",
		"https://www.naver.com/",
		"https://www.airbnb.co.kr/",
		"https://www.facebook.com/",
		"https://www.instagram.com/?hl=ko",
		"http://portal.seoultech.ac.kr/",
		"https://www.youtube.com/",
	}

	for _, url := range urls {

		go hitURL(url, channel)

	}
	for i := 0; i < len(urls); i++ {
		result := <-channel
		results[result.url] = result.status
	}

	for url, status := range results {
		fmt.Println(url, status)
	}

}

func hitURL(url string, c chan<- Requestresult) {
	resp, err := http.Get(url)
	status := "OK"
	if err != nil || resp.StatusCode >= 400 {
		status = "FAILED"
	}
	c <- Requestresult{url: url, status: status}

}
