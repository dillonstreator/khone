package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dillonstreator/khone"
)

type Todo struct {
	UserID    int
	ID        int
	Title     string
	Completed bool
}

func (p *Todo) UnmarshalCSV(m map[string]string) error {
	userID, err := strconv.Atoi(m["userId"])
	if err != nil {
		return err
	}
	p.UserID = userID

	ID, err := strconv.Atoi(m["id"])
	if err != nil {
		return err
	}
	p.ID = ID

	p.Title = m["title"]

	completed, err := strconv.ParseBool(m["completed"])
	if err != nil {
		return err
	}
	p.Completed = completed

	return nil
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	file := bytes.NewReader([]byte(`userId,id,title,completed
1,1,"delectus aut autem",false
1,2,"quis ut nam facilis et officia qui",false
2,3,"fugiat veniam minus",true`))

	err := khone.StreamCSV(ctx, file, func(t *Todo, _ int) {
		fmt.Println(t)
	}, khone.WithConcurrency(3))
	if err != nil {
		log.Fatal(err)
	}
}
