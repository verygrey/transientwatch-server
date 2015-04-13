package core

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

type Record struct {
	Id      string
	Title   string
	Body    string
	Source  string
	Url     string
	PubDate string
}

type Node struct {
	Data *Record
	Prev *Node
	Next *Node
}

type DataStore struct {
	Head     *Node
	Tail     *Node
	Capacity int
	Size     int
}

var alpha = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func Srand() string {
	return SrandN(20)
}

// generates a random string of fixed size
func SrandN(size int) string {
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		buf[i] = alpha[rand.Intn(len(alpha))]
	}
	return string(buf)
}

func SendPushNotification(text string) {
	url := "https://mobile.ng.bluemix.net:443/push/v1/apps/c178f6c2-be4b-4a14-a16d-ecb00da18089/messages"
	secret := "96cdaf6a61c377a6301c7022288236d30c41c290"
	body := fmt.Sprintf(`{"message": {"alert":"%s"}}`, text)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		log.Println("Error creating request", err)
	}
	req.Header.Add("IBM-Application-Secret", secret)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request", err)
	}
	if resp.StatusCode == 202 {
		log.Println("Push message was sent: ", text)
	} else {
		log.Println("Error sending push ", resp.StatusCode, " ", text)
	}
}

func NewDataStore(capacity int) *DataStore {
	head, tail := &Node{nil, nil, nil}, &Node{nil, nil, nil}
	head.Next = tail
	tail.Prev = head
	return &DataStore{head, tail, capacity, 0}
}

func (ds *DataStore) Add(record *Record) {
	if ds.Size >= ds.Capacity {
		ds.Tail = ds.Tail.Prev
		ds.Tail.Next = nil
		ds.Size--
	}
	newNode := &Node{record, ds.Head, ds.Head.Next}
	ds.Head.Next.Prev = newNode
	ds.Head.Next = newNode
	ds.Size++
}

func (ds *DataStore) String() string {
	return fmt.Sprintf("Size %d capacity %d", ds.Size, ds.Capacity)
}

func (ds *DataStore) Slice() []*Record {
	out := make([]*Record, ds.Size)
	curr := ds.Head.Next
	for i := range out {
		out[i] = curr.Data
		curr = curr.Next
	}
	return out
}

func (ds *DataStore) dump() {
	curr := ds.Head.Next
	for {
		fmt.Printf(" %s ->", curr.Data.Title)
		if curr.Next == ds.Tail {
			break
		}
		curr = curr.Next

	}
}
