package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hashmap/transientwatch/core"
	"log"
	"net/http"
	"os"
	//for extracting service credentials from VCAP_SERVICES
	//"github.com/cloudfoundry-community/go-cfenv"
)

const (
	DEFAULT_PORT = "8080"
	DEFAULT_HOST = "localhost"
)

var out chan core.Record = make(chan core.Record)
var ds *core.DataStore = core.NewDataStore(5)

func GetFeed(c *gin.Context) {
	c.JSON(http.StatusOK, ds.Slice())
}

func crawl() {
	go core.PollFeed("http://www.astronomerstelegram.org/?rss", 5, nil, out)
}

func receive() {
	for {
		rec := <-out
		ds.Add(&rec)
		log.Println("GOT: %s", rec.Title)
	}
}

func main() {
	var port string
	if port = os.Getenv("VCAP_APP_PORT"); len(port) == 0 {
		port = DEFAULT_PORT
	}

	var host string
	if host = os.Getenv("VCAP_APP_HOST"); len(host) == 0 {
		host = DEFAULT_HOST
	}

	crawl()
	go receive()

	r := gin.Default()
	r.GET("/feed", GetFeed)

	r.Run(host + ":" + port)
	log.Printf("Starting app on %+v:%+v\n", host, port)
}
