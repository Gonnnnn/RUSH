package array

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	assert.Equal(t, []int{2, 3, 4}, Map([]int{1, 2, 3}, func(x int) int { return x + 1 }))
	assert.Equal(t, []string{"2", "3", "4"}, Map([]int{1, 2, 3}, func(x int) string { return strconv.Itoa(x + 1) }))
	assert.Equal(t, []string{"", "", ""}, Map([]int{1, 2, 3}, func(x int) string { return "" }))

	assert.Equal(t, []int{}, Map([]int{}, func(x int) int { return x + 1 }))
}

func TestFilter(t *testing.T) {
	assert.Equal(t, []int{2, 4}, Filter([]int{1, 2, 3, 4, 5}, func(x int) bool { return x%2 == 0 }))
	assert.Equal(t, []int{}, Filter([]int{1, 3, 5}, func(x int) bool { return x%2 == 0 }))
	assert.Equal(t, []int{1, 3, 5}, Filter([]int{1, 2, 3, 4, 5}, func(x int) bool { return x%2 != 0 }))

	assert.Equal(t, []int{}, Filter([]int{}, func(x int) bool { return x%2 != 0 }))
}

func TestContains(t *testing.T) {
	assert.True(t, Contains([]int{1, 2, 3}, 2))
	assert.False(t, Contains([]int{1, 2, 3}, 4))
	assert.False(t, Contains([]int{}, 4))
}
