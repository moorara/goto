package sort

import (
	"testing"

	. "github.com/moorara/go-box/dt"
	"github.com/moorara/go-box/util"
	"github.com/stretchr/testify/assert"
)

func TestShellSortInt(t *testing.T) {
	tests := []struct {
		compare Compare
		items   []Generic
	}{
		{CompareInt, []Generic{}},
		{CompareInt, []Generic{20, 10, 30}},
		{CompareInt, []Generic{30, 20, 10, 40, 50}},
		{CompareInt, []Generic{90, 80, 70, 60, 50, 40, 30, 20, 10}},
	}

	for _, test := range tests {
		ShellSort(test.items, test.compare)

		assert.True(t, util.IsSorted(test.items, test.compare))
	}
}

func TestShellSortString(t *testing.T) {
	tests := []struct {
		compare Compare
		items   []Generic
	}{
		{CompareString, []Generic{}},
		{CompareString, []Generic{"Milad", "Mona"}},
		{CompareString, []Generic{"Alice", "Bob", "Alex", "Jackie"}},
		{CompareString, []Generic{"Docker", "Kubernetes", "Go", "JavaScript", "Elixir", "React", "Redux", "Vue"}},
	}

	for _, test := range tests {
		ShellSort(test.items, test.compare)

		assert.True(t, util.IsSorted(test.items, test.compare))
	}
}

func BenchmarkShellSortInt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		items := util.GenerateIntSlice(1000, -1000, 1000)
		ShellSort(items, CompareInt)
	}
}

func BenchmarkShellSortString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		items := util.GenerateStringSlice(1000, 10, 50)
		ShellSort(items, CompareString)
	}
}
