package lb

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRandomized(t *testing.T) {
	t.Run("panics for empty endpoints", func(t *testing.T) {
		assert.Panics(t, func() {
			Randomized([]string{})
		})
	})

	t.Run("returns the only endpoint", func(t *testing.T) {
		selector := Randomized([]string{"only"})

		assert.Equal(t, 1, selector.Count())
		assert.Equal(t, "only", selector.Get())
		assert.Equal(t, "only", selector.Get())
	})

	t.Run("copies the endpoint slice", func(t *testing.T) {
		endpoints := []string{"a", "b"}
		selector := Randomized(endpoints)
		endpoints[0] = "mutated"

		for i := 0; i < 20; i++ {
			assert.NotEqual(t, "mutated", selector.Get())
		}
	})

	t.Run("returns only configured endpoints under concurrent use", func(t *testing.T) {
		selector := Randomized([]string{"a", "b", "c"})
		allowed := map[string]bool{"a": true, "b": true, "c": true}
		results := make(chan string, 100)
		var wg sync.WaitGroup

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				results <- selector.Get()
			}()
		}
		wg.Wait()
		close(results)

		for result := range results {
			require.True(t, allowed[result], "unexpected endpoint %q", result)
		}
	})
}
