package main

import (
	"crypto/rand"
	"encoding/gob"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/go-vk-api/vk"
)

var (
	KEY    string
	Jokes  []string
	Client *vk.Client
)

func getRandomJoke(jokes []string) string {
	num, err := rand.Int(rand.Reader, big.NewInt(int64(len(jokes))))
	if err != nil {
		log.Fatalln(err)
	}
	return jokes[num.Int64()]
}

func init() {
	KEY = os.Getenv("KEY")
	if KEY == "" {
		log.Fatalln("Specify a key for application")
	}
	file, err := os.Open("jokes.gob")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	decoder.Decode(&Jokes)
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

func main() {
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
