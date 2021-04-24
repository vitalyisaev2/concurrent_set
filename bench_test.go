package set

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

func BenchmarkSet(b *testing.B) {
	f := factory{}

	kinds := []kind{
		coarseGrained,
		fineGrained,
	}

	type dataSource struct {
		name string
		data []int
	}

	const inputLength = 2 << 9

	dataSources := []dataSource{
		{
			name: "ascending_array_input",
			data: makeAscendingArray(inputLength),
		},
		{
			name: "descending_array_input",
			data: makeDescendingArray(inputLength),
		},
		{
			name: "shuffled_array_input",
			data: makeShuffledArray(inputLength),
		},
	}

	//threadNumbers := []int{1, 2, 4, 8, 16, 32, 64}
	threadNumbers := []int{1, 4, 16}

	for _, threadNumber := range threadNumbers {
		threadNumber := threadNumber
		b.Run(fmt.Sprintf("%v_threads", threadNumber), func(b *testing.B) {
			for _, ds := range dataSources {
				ds = ds
				b.Run(ds.name, func(b *testing.B) {
					for _, k := range kinds {
						k := k
						b.Run(k.String(), func(b *testing.B) {
							b.Run("concurrent insertion", func(b *testing.B) {
								set := f.new(k)

								wg := sync.WaitGroup{}
								wg.Add(threadNumber)

								for i := 0; i < threadNumber; i++ {
									go func() {
										defer wg.Done()
										for j := 0; j < b.N; j++ {
											ix := j % len(ds.data)
											set.Insert(ds.data[ix])
										}
									}()
								}
								wg.Wait()
							})
						})
					}
				})
			}
		})
	}
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
