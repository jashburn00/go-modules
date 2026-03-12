package autocompletetrie

import (
	"testing"
)

func TestMakeAutocompleteTrie(t *testing.T) {
	t.Run("creates a trie without error", func(t *testing.T) {
		//setup
		words := []string{"Hello", "World"}
		//exercise
		_, err := MakeAutocompleteTrie(words)
		//verify
		if err != nil {
			t.Errorf("Error creating trie: %s", err.Error())
		}
		//teardown
	})

	t.Run("generates a usable autocomplete trie", func(t *testing.T) {
		//setup
		words := []string{"Hello", "World", "Hematology", "Hercules"}
		trie, err := MakeAutocompleteTrie(words)
		if err != nil {
			t.Fatalf("Error creating trie: %s", err.Error())
		}
		//exercise
		completion := trie.Autocomplete("He")
		//verify
		if completion != "Hello" {
			t.Errorf("Expected autocomplete to return 'Hello' but got '%s'", completion)
		}
		//teardown
	})
}

func TestAddWord(t *testing.T) {
	t.Run("adds a word to the trie", func(t *testing.T) {
		//setup
		trie, err := MakeAutocompleteTrie([]string{"frog"})
		if err != nil {
			t.Fatalf("Error creating trie: %s", err.Error())
		}
		//exercise
		err = trie.AddWord("lemur")
		if err != nil {
			t.Fatalf("Error adding word to trie: %s", err.Error())
		}
		//verify
		completion := trie.Autocomplete("le")
		if completion != "lemur" {
			t.Errorf("Expected autocomplete to return 'lemur' but got '%s'", completion)
		}
	})
}

func TestGetAllCompletions(t *testing.T) {
	t.Run("returns all completions for a given query", func(t *testing.T) {
		//setup
		trie, err := MakeAutocompleteTrie([]string{"frog", "bat", "frond"})
		if err != nil {
			t.Fatalf("Error creating trie: %s", err.Error())
		}
		//exercise
		prefix, completions := trie.GetAllCompletions("I see a fr")
		//verify
		if prefix != "I see a " {
			t.Errorf("Expected prefix to be 'I see a ' but got %q", prefix)
		}
		expected := []string{"frog", "frond"}
		if len(completions) != len(expected) {
			t.Errorf("Expected %d completions but got %d", len(expected), len(completions))
		}
		for i, v := range expected {
			if completions[i] != v {
				t.Errorf("Expected completions %v but got %v", expected, completions)
			}
		}
	})

	t.Run("returns the input string as the only completion if there are no matches", func(t *testing.T) {
		//setup
		trie, err := MakeAutocompleteTrie([]string{"frog", "bat", "frond"})
		if err != nil {
			t.Fatalf("Error creating trie: %s", err.Error())
		}
		//exercise
		prefix, completions := trie.GetAllCompletions("North Ame")
		//verify
		if prefix != "North " {
			t.Errorf("Expected prefix to be 'North ' but got %q", prefix)
		}
		if len(completions) != 1 || completions[0] != "Ame" {
			t.Errorf("Expected completions to be ['Ame'] but got %v", completions)
		}
	})
}
