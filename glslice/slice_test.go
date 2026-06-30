package glslice

import (
	"reflect"
	"testing"
)

func TestContains(t *testing.T) {
	items := []string{"go", "php", "mysql"}

	if !Contains(items, "php") {
		t.Fatal("Contains() = false, want true")
	}
	if Contains(items, "redis") {
		t.Fatal("Contains() = true, want false")
	}
}

func TestUniquePreservesFirstSeenOrder(t *testing.T) {
	items := []int{3, 1, 3, 2, 1}
	want := []int{3, 1, 2}

	if got := Unique(items); !reflect.DeepEqual(got, want) {
		t.Fatalf("Unique() = %#v, want %#v", got, want)
	}
}

func TestFilter(t *testing.T) {
	items := []int{1, 2, 3, 4}
	want := []int{2, 4}

	got := Filter(items, func(v int) bool {
		return v%2 == 0
	})
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Filter() = %#v, want %#v", got, want)
	}
}

func TestCompactStrings(t *testing.T) {
	items := []string{" go ", "", "\t", "php"}
	want := []string{"go", "php"}

	if got := CompactStrings(items); !reflect.DeepEqual(got, want) {
		t.Fatalf("CompactStrings() = %#v, want %#v", got, want)
	}
}
