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
var ds *core.DataStore = core.NewDataStore(100)

func GetFeed(c *gin.Context) {
	c.JSON(http.StatusOK, ds.Slice())
}

func crawl() {
	go core.PollFeed("http://www.astronomerstelegram.org/?rss+Neutron%20Star", 5, nil, out, 10, "The Astronomer's Telegram")
	go core.PollFeed("http://www.astronomerstelegram.org/?rss+Soft%20Gamma-ray%20Repeater", 5, nil, out, 10, "The Astronomer's Telegram")
	go core.PollFeed("http://www.astronomerstelegram.org/?rss+Black%20Hole", 5, nil, out, 10, "The Astronomer's Telegram")
	go core.PollFeed("http://www.astronomerstelegram.org/?rss+Gamma-Ray%20Burst", 5, nil, out, 10, "The Astronomer's Telegram")
	go core.PollFeed("http://www.astronomerstelegram.org/?rss+Transient", 5, nil, out, 10, "The Astronomer's Telegram")
	go core.PollFeed("http://www.cbat.eps.harvard.edu/unconf/tocp.xml", 5, nil, out, 10, "Central Bureau for Astronomical Telegrams")
	go core.PollGCN(1, 30, out)
}

func receive() {
	for {
		rec := <-out
		ds.Add(&rec)
		log.Println("GOT: %s", rec.Title, rec.Url)
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
