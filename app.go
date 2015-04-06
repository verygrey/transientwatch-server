package main

import (
	"github.com/gin-gonic/gin"
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

func GetFeed(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id": 1, "url": "http://www.ibm.com"})
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

	r := gin.Default()
	r.GET("/feed", GetFeed)

	r.Run(host + ":" + port)
	log.Printf("Starting app on %+v:%+v\n", host, port)
}
