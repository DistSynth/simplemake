package main

import (
	"fmt"

	"golang.org/x/sync/errgroup"
)

type fType[T any] func(item T) error

func forEach[T any](items []T, p int, f fType[T]) error {

	gr := new(errgroup.Group)
	gr.SetLimit(p)

	for _, el := range items {
		el := el
		gr.Go(func() error {
			return f(el)
		})
	}
	return gr.Wait()
}

func main() {
	arr := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}

	forEach(arr, 2, func(item string) error {
		fmt.Println(item)
		return nil
	})

}
