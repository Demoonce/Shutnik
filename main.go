package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-vk-api/vk"
)

var (
	Key string

	Jokes     []string
	Client    *vk.Client
	JokesChan = make(chan string)
)

func init() {
	Key = os.Getenv("KEY")
	if Key == "" {
		log.Fatalln("Specify a key for application")
	}
	var err error
	Client, err = vk.NewClientWithOptions(vk.WithToken(Key))
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

	// Check if file already exists
	if _, err := os.Stat("jokes.gob"); err != nil {
		Jokes = ParseJokes()
		WriteJokes(Jokes)
	}
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
			_, err := http.Get("https://shutnik.herokuapp.com/")
			if err != nil {
				log.Println(err)
			}
			time.Sleep(time.Minute)
			sendJoke()
		}
		time.Sleep(time.Second * 5)
	}
}
