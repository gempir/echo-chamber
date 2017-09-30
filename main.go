package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gempir/go-twitch-irc"
	"github.com/labstack/echo"

	"gopkg.in/olivere/elastic.v5"
)

type msg struct {
	Text     string    `json:"text"`
	Username string    `json:"username"`
	Time     time.Time `json:"time"`
}

// Streams struct
type Streams struct {
	Streams []Stream `json:"streams"`
}

// Stream struct
type Stream struct {
	Channel Channel `json:"channel"`
}

// Channel struct
type Channel struct {
	Name string `json:"name"`
}

var (
	esClient *elastic.Client
)

func main() {
	// wait for ES to start
	time.Sleep(time.Second * 20)

	e := echo.New()

	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	e.Renderer = t

	e.GET("/", home)
	e.Static("/static", "static")
	e.GET("/search", search)

	ctx := context.Background()

	var err error
	esClient, err = elastic.NewClient(elastic.SetURL("http://elasticsearch:9200"))
	if err != nil {
		panic(err)
	}

	tclient := twitch.NewClient("justinfan123123", "oauth:123123123")

	tclient.OnNewMessage(func(channel string, user twitch.User, message twitch.Message) {
		esMessage := msg{
			Text:     message.Text,
			Username: user.Username,
			Time:     message.Time,
		}

		_, err := esClient.Index().
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

	go tclient.Connect()

	e.Logger.Fatal(e.Start(":1323"))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func home(c echo.Context) error {

	return c.Render(http.StatusOK, "index", "world")
}

func search(c echo.Context) error {
	ctx := context.Background()
	query := c.QueryParam("q")

	simpleQueryStringQuery := elastic.NewSimpleQueryStringQuery(query)
	randomFunction := elastic.NewRandomFunction()
	randomFunction.Weight(100000000)

	fnScoreQuery := elastic.NewFunctionScoreQuery()
	fnScoreQuery.Add(simpleQueryStringQuery, randomFunction)

	result, err := esClient.Search().Query(fnScoreQuery).Size(1).Do(ctx)
	if err != nil {
		c.Error(err)
	}

	return c.JSON(http.StatusOK, result)
}

func getTopChannels() []Stream {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.twitch.tv/kraken/streams?limit=30", nil)
	req.Header.Set("Client-Id", getEnv("CLIENTID"))
	res, _ := client.Do(req)

	contents, _ := ioutil.ReadAll(res.Body)
	var streams Streams
	json.Unmarshal(contents, &streams)

	return streams.Streams
}

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	panic("Missing env var: " + key)
}
