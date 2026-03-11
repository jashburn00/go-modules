package trie

import (
	"reflect"
	"testing"
)

func TestNewTrie(t *testing.T) {
	tr := NewTrie()
	if tr == nil {
		t.Fatal("NewTrie returned nil")
	}
	if tr.root == nil {
		t.Fatal("NewTrie root is nil")
	}
}

func TestInsertAndSearch(t *testing.T) {
	tr := NewTrie()

	// Words not yet inserted should not be found.
	if tr.Search("hello") {
		t.Error("expected Search(\"hello\") to be false before insertion")
	}

	tr.Insert("hello")

	if !tr.Search("hello") {
		t.Error("expected Search(\"hello\") to be true after insertion")
	}
	// A prefix alone is not a complete word.
	if tr.Search("hell") {
		t.Error("expected Search(\"hell\") to be false; it is only a prefix")
	}
	// The empty string is not a word unless explicitly inserted.
	if tr.Search("") {
		t.Error("expected Search(\"\") to be false")
	}
}

func TestInsertEmptyString(t *testing.T) {
	tr := NewTrie()
	tr.Insert("")
	if !tr.Search("") {
		t.Error("expected Search(\"\") to be true after inserting empty string")
	}
}

func TestInsertDuplicate(t *testing.T) {
	tr := NewTrie()
	tr.Insert("apple")
	tr.Insert("apple")
	if !tr.Search("apple") {
		t.Error("expected Search(\"apple\") to be true after duplicate insertion")
	}
}

func TestStartsWith(t *testing.T) {
	tr := NewTrie()
	tr.Insert("hello")
	tr.Insert("world")

	cases := []struct {
		prefix string
		want   bool
	}{
		{"hel", true},
		{"hello", true},
		{"he", true},
		{"h", true},
		{"wor", true},
		{"world", true},
		{"xyz", false},
		{"hellop", false},
	}

	for _, tc := range cases {
		got := tr.StartsWith(tc.prefix)
		if got != tc.want {
			t.Errorf("StartsWith(%q) = %v; want %v", tc.prefix, got, tc.want)
		}
	}
}

func TestAutoComplete(t *testing.T) {
	tr := NewTrie()
	words := []string{"apple", "app", "application", "apply", "banana", "band"}
	for _, w := range words {
		tr.Insert(w)
	}

	cases := []struct {
		prefix string
		want   []string
	}{
		{"app", []string{"app", "apple", "application", "apply"}},
		{"appl", []string{"apple", "application", "apply"}},
		{"ban", []string{"banana", "band"}},
		{"banana", []string{"banana"}},
		{"xyz", []string{}},
		{"", []string{"app", "apple", "application", "apply", "banana", "band"}},
	}

	for _, tc := range cases {
		got := tr.AutoComplete(tc.prefix)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("AutoComplete(%q) = %v; want %v", tc.prefix, got, tc.want)
		}
	}
}

func TestAutoCompleteNoMatches(t *testing.T) {
	tr := NewTrie()
	tr.Insert("hello")
	got := tr.AutoComplete("xyz")
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}

func TestUnicodeWords(t *testing.T) {
	tr := NewTrie()
	tr.Insert("café")
	tr.Insert("cafétéria")
	tr.Insert("cat")

	if !tr.Search("café") {
		t.Error("expected Search(\"café\") to be true")
	}
	got := tr.AutoComplete("café")
	want := []string{"café", "cafétéria"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("AutoComplete(\"café\") = %v; want %v", got, want)
	}
}
