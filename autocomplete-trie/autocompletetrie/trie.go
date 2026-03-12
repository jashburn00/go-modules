package main

import (
	"strings"
)

type BadTrieInputError struct {
	Message string
}

func (BadTrieInputError) Error() string {
	return "Bad input given for trie operation"
}

type AutocompleteTrie struct {
	Children map[byte]*AutocompleteTrie
	Values   []string
}

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
		return strings.Join(words[0:len(words)-1], " ") + " ", []string{query}
	}
	return strings.Join(words[0:len(words)-1], " ") + " ", node.Values
}
