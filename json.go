package khone

import (
	"context"
	"encoding/json"
	"io"
	"sync"
)

func StreamJSON[T any](ctx context.Context, r io.Reader, cb func(*T, int), options ...option) error {
	config := newConfig(options...)

	dec := json.NewDecoder(r)

	_, err := dec.Token()
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	ch := make(chan struct{}, config.concurrency)
	defer func() { close(ch) }()

	i := 0

	for dec.More() {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case ch <- struct{}{}:
			var val T
			err := dec.Decode(&val)
			if err != nil {
				return err
			}

			wg.Add(1)
			go func(idx int) {
				defer wg.Done()

				cb(&val, idx)
				<-ch
			}(i)
			i++

		}
	}

	wg.Wait()

	return nil
}
