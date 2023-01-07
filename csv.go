package khone

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"sync"
)

// https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#pointer-method-example
type csvUnmarshaler[T any] interface {
	UnmarshalCSV(map[string]string) error
	*T
}

func StreamCSV[T csvUnmarshaler[V], V any](ctx context.Context, r io.Reader, cb func(T, int), options ...option) error {
	config := newConfig(options...)

	csvReader := csv.NewReader(r)

	headers, err := csvReader.Read()
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	ch := make(chan struct{}, config.concurrency)
	defer func() { close(ch) }()

	i := 0

readLoop:
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case ch <- struct{}{}:
			record, err := csvReader.Read()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break readLoop
				}

				return err
			}

			m := map[string]string{}

			for i, fieldValue := range record {
				m[headers[i]] = fieldValue
			}

			var val T = new(V)
			err = val.UnmarshalCSV(m)
			if err != nil {
				return err
			}

			wg.Add(1)
			go func(idx int) {
				defer wg.Done()

				cb(val, idx)
				<-ch
			}(i)
			i++

		}
	}

	wg.Wait()

	return nil
}
