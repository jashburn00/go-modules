package main

import (
	"strings"
)

// BadTrieInputError will be returned when arguments prevent trie operations from working,
// such as adding an empty string to the trie.
type BadTrieInputError struct {
	Message string
}

func (BadTrieInputError) Error() string {
	return "Bad input given for trie operation"
}

// AutocompleteTrie follows a typical trie structure, having a map of child nodes and a list of values.
type AutocompleteTrie struct {
	Children map[byte]*AutocompleteTrie
	Values   []string
}

// GetNodeAt returns an AutocompleteTrie pointer found by traversing the trie with the specified bytes.
func (t *AutocompleteTrie) GetNodeAt(path []byte) *AutocompleteTrie { //abc
	traverser := t
	for _, b := range path {
		val, ok := traverser.Children[b]
		if !ok {
			return nil
		} else {
			traverser = val
		}
	}
	return traverser
}

// Creates an AutocompleteTrie and populates it with the given array of strings, returning a pointer to the root node.
//
// Strings will be added in the order that they appear in the array argument, meaning that
// if two words are available to autocomplete at the same node, the one that appeared first in
// the array argument will be at the front of the .Values array. This impacts results of the
// Autocomplete function.
func MakeAutocompleteTrie(words []string) (*AutocompleteTrie, error) {
	if len(words) < 1 {
		return nil, BadTrieInputError{"list of words to build trie upon was empty"}
	}

	trie := &AutocompleteTrie{Children: map[byte]*AutocompleteTrie{}, Values: []string{}}

	for _, w := range words {
		err := trie.AddWord(w)
		if err != nil {
			return nil, err
		}
	}

	return trie, nil
}

// Adds a single word to the trie (usually the root node). Words added appear at the end of each node's
// .Values array, meaning they will not be selected by the Autocomplete function if there is an existing
// value at any applicable nodes.
func (trie *AutocompleteTrie) AddWord(w string) error {
	if len(w) < 1 {
		return BadTrieInputError{"Word to add to trie was empty"}
	}

	var b byte
	length := len([]byte(w))
	traverser := trie

	for i := 0; i < length; i++ {
		b = []byte(w)[i]
		if traverser.GetNodeAt([]byte{b}) == nil {
			traverser.Children[b] = &AutocompleteTrie{Children: map[byte]*AutocompleteTrie{}, Values: []string{w}}
		} else {
			traverser.Children[b].Values = append(traverser.Children[b].Values, w)
		}
		traverser = traverser.GetNodeAt([]byte{b})
	}

	return nil
}

// Autocomplete returns the string parameter with the last word in the string
// completed to the first value found at the applicable node. If no completion
// is found, the parameter is returned unchanged.
//
// This is where the order in which strings are added to the trie matters -
// the earliest added completion will be used at each node.
func (trie *AutocompleteTrie) Autocomplete(query string) string {
	if len(query) < 1 {
		return ""
	}

	words := strings.Split(query, " ")
	if len(words) == 1 {
		node := trie.GetNodeAt([]byte(words[0]))
		if node == nil {
			return query
		}
		if len(node.Values) < 1 || len(node.Values[0]) < 1 {
			panic("Want to autocomplete " + query + " but trie node Values was empty or had a bad string")
		}
		return node.Values[0]
	}

	node := trie.GetNodeAt([]byte(words[len(words)-1]))
	if node != nil {
		if len(node.Values) < 1 || len(node.Values[0]) < 1 {
			panic("Want to autocomplete " + query + " but trie node Values was empty or had a bad string")
		}
		words[len(words)-1] = node.Values[0]
	}

	var result string
	for i, v := range words {
		if i == len(words)-1 {
			result = result + v
		} else {
			result = result + v + " "
		}
	}

	return result
}

// GetAllCompletions returns a prefix containing any words in the parameter before the last word (if any),
// along with the array of completions for the last word in the string parameter. If no completions
// are found, the string parameter is returned such that `prefix + completions[0]` equals the original string.
//
// The prefix is formatted so that `prefix + completions[i]` will effectively autocomplete the last
// word in the string parameter and leave the rest of the string unchanged.
func (trie *AutocompleteTrie) GetAllCompletions(query string) (string, []string) {
	if len(query) < 1 {
		return "", []string{""}
	}

	words := strings.Split(query, " ")
	if len(words) == 1 {
		node := trie.GetNodeAt([]byte(words[0]))
		if node == nil {
			return "", []string{query}
		}
		return "", node.Values
	}

	node := trie.GetNodeAt([]byte(words[len(words)-1]))
	if node == nil {
		return strings.Join(words[0:len(words)-1], " ") + " ", []string{words[len(words)-1]}
	}
	return strings.Join(words[0:len(words)-1], " ") + " ", node.Values
}
