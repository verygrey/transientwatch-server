package core

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func PollGCN(timeout, limit int, out chan Record) {
	startMark := "<!XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX>"
	endMark := "<!YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY>"
	for {
		page, err := GetPage("http://gcn.gsfc.nasa.gov/gcn/gcn3_archive.html")
		if err != nil {
			log.Println("HTTP error ", err)
		} else {
			if strings.TrimSpace(page) == "" {
				continue
			}
			log.Println("End mark ", endMark)
			lines := ExtractArea(page, startMark, endMark, limit)
			if lines[0] != "" {
				endMark = lines[0]
				processLines(lines, out)
			}
		}
		time.Sleep(time.Duration(timeout) * time.Minute)
	}
}

func processLines(lines []string, out chan Record) {
	hrefRx := regexp.MustCompile("HREF=(.*?)>")
	idRx := regexp.MustCompile(`>(\d+?)<`)
	descrRx := regexp.MustCompile(`</A> (.+?)<br>`)
	for _, line := range lines {
		line := line
		defer func() {
			out <- Record{"GCN Circular " + extractData(idRx, line),
				extractData(descrRx, line), "GCN Circulars",
				"http://gcn.gsfc.nasa.gov/" + extractData(hrefRx, line), ""}
		}()
	}
}

func extractData(rx *regexp.Regexp, line string) (out string) {
	matches := rx.FindStringSubmatch(line)
	if matches != nil {
		out = matches[1]
	} else {
		out = ""
	}
	return
}

func ExtractArea(page, startMark, endMark string, limit int) []string {
	page = strings.TrimSpace(page)
	log.Println("Page length ", len(page), " start ", startMark, " end ", endMark)
	lines := strings.Split(page[strings.Index(page, startMark)+len(startMark):strings.Index(page, endMark)], "\n")
	if len(lines) < limit {
		return lines
	}
	return lines[1:limit]
}

func GetPage(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Print("HTTP error on ", url, err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("Error reading response body on ", url, err)
		return "", err
	}
	return string(body), nil
}
