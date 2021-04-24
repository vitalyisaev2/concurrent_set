package set

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

func BenchmarkSet(b *testing.B) {
	kinds := []setKind{
		coarseGrained,
		fineGrained,
	}

	const inputLength = 2 << 9

	dataSources := []*dataSource{
		{name: "ascending_array_input", data: makeAscendingArray(inputLength)},
		{name: "descending_array_input", data: makeDescendingArray(inputLength)},
		{name: "shuffled_array_input", data: makeShuffledArray(inputLength)},
	}

	threadNumbers := []int{1, 4, 16}

	// combination of parameters
	for _, threadNumber := range threadNumbers {
		threadNumber := threadNumber

		b.Run(fmt.Sprintf("%v_threads", threadNumber), func(b *testing.B) {
			for _, ds := range dataSources {
				ds := ds

				b.Run(ds.name, func(b *testing.B) {
					for _, kind := range kinds {
						kind := kind

						b.Run(kind.String(), func(b *testing.B) {
							params := &benchParams{kind: kind, threads: threadNumber, dataSource: ds}

							b.Run("concurrent insertion", func(b *testing.B) {
								benchConcurrentInsertion(b, params)
							})
						})
					}
				})
			}
		})
	}
}

type dataSource struct {
	name string
	data []int
}

func makeAscendingArray(length int) []int {
	output := make([]int, length)
	for i := 0; i < length; i++ {
		output[i] = i
	}

	return output
}

func makeDescendingArray(length int) []int {
	output := make([]int, length)
	for i := 0; i < length; i++ {
		output[length-i-1] = i
	}

	return output
}

func makeShuffledArray(length int) []int {
	output := makeAscendingArray(length)
	rand.Shuffle(length, func(i, j int) {
		output[i], output[j] = output[j], output[i]
	})

	return output
}

type benchParams struct {
	dataSource *dataSource
	threads    int
	kind       setKind
}

func benchConcurrentInsertion(b *testing.B, params *benchParams) {
	b.Helper()

	f := factory{}

	set := f.new(params.kind)

	wg := sync.WaitGroup{}
	wg.Add(params.threads)

	for i := 0; i < params.threads; i++ {
		go func() {
			defer wg.Done()

			for j := 0; j < b.N; j++ {
				ix := j % len(params.dataSource.data)
				val := params.dataSource.data[ix]
				set.Insert(val)
			}
		}()
	}
	wg.Wait()
}
