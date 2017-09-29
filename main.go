package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gempir/go-twitch-irc"
	elastic "gopkg.in/olivere/elastic.v5"
)

type msg struct {
	Text     string    `json:"text"`
	Username string    `json:"username"`
	Time     time.Time `json:"time"`
}

// Streams asd
type Streams struct {
	Streams []Stream `json:"streams"`
}

// Stream asd
type Stream struct {
	Channel Channel `json:"channel"`
}

// Channel asd
type Channel struct {
	Name string `json:"name"`
}

func main() {
	time.Sleep(time.Second * 5)
	ctx := context.Background()

	client, err := elastic.NewClient(elastic.SetURL("http://" + getEnv("ESHOST", "127.0.0.1:9200")))
	if err != nil {
		// Handle error
		panic(err)
	}

	tclient := twitch.NewClient("justinfan123123", "oauth:123123123")

	tclient.OnNewMessage(func(channel string, user twitch.User, message twitch.Message) {
		// Add a document to the index
		esMessage := msg{
			Text:     message.Text,
			Username: user.Username,
			Time:     message.Time,
		}

		_, err := client.Index().
			Index(channel).
			Type("doc").
			BodyJson(esMessage).
			Refresh("true").
			Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
		}
	})

	go func() {
		for {
			top := getTopChannels()
			for _, channel := range top {
				fmt.Printf("Joining: %s\r\n", channel.Channel.Name)
				go tclient.Join(channel.Channel.Name)
			}
			time.Sleep(time.Hour)
		}
	}()

	go tclient.Join("pajlada")
	go tclient.Join("gempbot")
	go tclient.Join("gempir")
	go tclient.Join("forsenlol")
	go tclient.Join("jaxerie")
	go tclient.Join("nuuls")
	go tclient.Join("imsoff")
	go tclient.Join("nymn")
	go tclient.Join("nanilul")
	go tclient.Join("xfsn_saber")

	tclient.Connect()
}

func getTopChannels() []Stream {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.twitch.tv/kraken/streams?limit=30", nil)
	req.Header.Set("Client-Id", "cazb1iyx9igrhk42ruep3e6dit84id")
	res, _ := client.Do(req)

	contents, _ := ioutil.ReadAll(res.Body)
	var streams Streams
	json.Unmarshal(contents, &streams)

	return streams.Streams
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
