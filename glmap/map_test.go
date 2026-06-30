package glmap

import (
	"reflect"
	"sort"
	"testing"
)

func TestKeys(t *testing.T) {
	got := Keys(map[string]int{"b": 2, "a": 1})
	sort.Strings(got)

	want := []string{"a", "b"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Keys() = %#v, want %#v", got, want)
	}
}

func TestValues(t *testing.T) {
	got := Values(map[string]int{"a": 2, "b": 1})
	sort.Ints(got)

	want := []int{1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Values() = %#v, want %#v", got, want)
	}
}

func TestHasKey(t *testing.T) {
	items := map[string]int{"a": 1}

	if !HasKey(items, "a") {
		t.Fatal("HasKey() = false, want true")
	}
	if HasKey(items, "b") {
		t.Fatal("HasKey() = true, want false")
	}
}

func TestMergeRightOverwritesLeft(t *testing.T) {
	left := map[string]int{"a": 1, "b": 2}
	right := map[string]int{"b": 20, "c": 3}
	want := map[string]int{"a": 1, "b": 20, "c": 3}

	got := Merge(left, right)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Merge() = %#v, want %#v", got, want)
	}
	if left["b"] != 2 {
		t.Fatalf("Merge() mutated left map: left[b] = %d, want 2", left["b"])
	}
}
