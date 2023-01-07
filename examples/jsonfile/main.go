package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dillonstreator/khone"
)

type Todo struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	file := bytes.NewReader([]byte(`
[
	{
		"userId": 1,
		"id": 1,
		"title": "delectus aut autem",
		"completed": false
	},
	{
		"userId": 1,
		"id": 2,
		"title": "quis ut nam facilis et officia qui",
		"completed": false
	},
	{
		"userId": 2,
		"id": 3,
		"title": "fugiat veniam minus",
		"completed": true
	}
]
`))

	err := khone.StreamJSON(ctx, file, func(t *Todo, _ int) {
		fmt.Println(t)
	})
	if err != nil {
		log.Fatal(err)
	}
}
