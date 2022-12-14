package main

import (
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-vk-api/vk"
)

func getRandomJoke(jokes []string) string {
	log.Println(jokes)
	num, err := rand.Int(rand.Reader, big.NewInt(int64(len(jokes))))
	if err != nil {
		log.Fatalln(err)
	}
	return jokes[num.Int64()]
}

func getJokesPage(url string, jokes_chan chan<- string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	jokes, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Println(err)
		return
	}
	jokes.Find("div.text").Each(func(i int, s *goquery.Selection) {
		jokes_chan <- s.Text()
	})
}

func ParseJokes() []string {
	jokes := make([]string, 0, 1000)
	for a := 1; a < 500; a++ {
		time.Sleep(time.Millisecond * 10)
		go getJokesPage(fmt.Sprintf("https://nekdo.ru/short/%d", a), JokesChan)
	}
loop:
	for {
		select {
		case <-time.After(time.Second * 5):
			close(JokesChan)
			break loop
		case a := <-JokesChan:
			jokes = append(jokes, a)
		}
	}
	return jokes
}

func WriteJokes(jokes []string) {
	file, err := os.Create("jokes.gob")
	defer file.Close()
	if err != nil {
		log.Fatalln(err)
	}
	encoder := gob.NewEncoder(file)
	encoder.Encode(jokes)
}

func sendJoke() {
	err := Client.CallMethod("wall.post", vk.RequestParams{
		"owner_id":   -217202035,
		"from_group": 1,
		"message":    getRandomJoke(Jokes),
	}, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
