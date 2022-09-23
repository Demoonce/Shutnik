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

var (
	KEY string

	Jokes     []string
	Client    *vk.Client
	JokesChan = make(chan string)
)

func getRandomJoke(jokes []string) string {
	num, err := rand.Int(rand.Reader, big.NewInt(int64(len(jokes))))
	if err != nil {
		log.Fatalln(err)
	}
	return jokes[num.Int64()]
}

func getJokesPage(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
		return
	}
	jokes, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Println(err)
		return
	}
	jokes.Find("div.text").Each(func(i int, s *goquery.Selection) {
		Jokes = append(Jokes, s.Text())
	})
}

func ParseJokes() {
	for a := 1; a < 500; a++ {
		time.Sleep(time.Millisecond * 10)
		go getJokesPage(fmt.Sprintf("https://nekdo.ru/page/%d", a))
	}
loop:
	for {
		select {
		case <-time.NewTimer(time.Second * 5).C:
			close(JokesChan)
			break loop
		case a := <-JokesChan:
			Jokes = append(Jokes, a)
		}
	}
}

func WriteJokes() {
	file, err := os.Create("jokes.gob")
	defer file.Close()
	if err != nil {
		log.Fatalln(err)
	}
	encoder := gob.NewEncoder(file)
	encoder.Encode(Jokes)
}

func init() {
	KEY = os.Getenv("KEY")
	if KEY == "" {
		log.Fatalln("Specify a key for application")
	}
	var err error
	Client, err = vk.NewClientWithOptions(vk.WithToken(KEY))
	if err != nil {
		log.Fatalln(err)
	}
}

func sendJoke() {
	err := Client.CallMethod("wall.post", vk.RequestParams{
		"owner_id":   -160130110,
		"from_group": 1,
		"message":    getRandomJoke(Jokes),
	}, nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Shutnik started"))
}

func main() {
	http.HandleFunc("/", Index)
	go http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	ParseJokes()
	WriteJokes()
	file, err := os.Open("jokes.gob")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	decoder.Decode(&Jokes)
	post_ticker := time.NewTicker(time.Hour * 3)
	sendJoke()
	for {
		select {
		case <-post_ticker.C:
			sendJoke()
		}
		time.Sleep(time.Second * 5)
	}
}
