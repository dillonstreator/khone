package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
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

	res, err := http.Get("https://jsonplaceholder.typicode.com/todos")
	if err != nil {
		log.Fatal(err)
	}

	err = khone.StreamJSON(ctx, res.Body, func(t *Todo, _ int) {
		fmt.Println(t)
	}, khone.WithConcurrency(5))
	if err != nil {
		log.Fatal(err)
	}
}
