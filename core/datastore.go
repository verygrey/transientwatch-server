package core

import "fmt"

type Record struct {
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
