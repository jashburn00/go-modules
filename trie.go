// Package trie provides a node-based trie data structure for text autocomplete.
package trie

import "sort"

// TrieNode represents a single node in the trie.
type TrieNode struct {
	children    map[rune]*TrieNode
	isEndOfWord bool
}

// newTrieNode creates and returns a new TrieNode.
func newTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[rune]*TrieNode),
	}
}

// Trie is a node-based prefix tree that supports word insertion and
// autocomplete-style prefix searches.
type Trie struct {
	root *TrieNode
}

// NewTrie creates and returns a new, empty Trie.
func NewTrie() *Trie {
	return &Trie{root: newTrieNode()}
}

// Insert adds the given word to the trie. Inserting a word that already
// exists is a no-op.
func (t *Trie) Insert(word string) {
	node := t.root
	for _, ch := range word {
		if _, ok := node.children[ch]; !ok {
			node.children[ch] = newTrieNode()
		}
		node = node.children[ch]
	}
	node.isEndOfWord = true
}

// Search returns true if word exists in the trie as a complete word.
func (t *Trie) Search(word string) bool {
	node := t.traverse(word)
	return node != nil && node.isEndOfWord
}

// StartsWith returns true if any word in the trie begins with prefix.
func (t *Trie) StartsWith(prefix string) bool {
	return t.traverse(prefix) != nil
}

// AutoComplete returns all words stored in the trie that begin with prefix.
// The returned slice is in lexicographic order. If no words match, an empty
// slice is returned.
func (t *Trie) AutoComplete(prefix string) []string {
	node := t.traverse(prefix)
	if node == nil {
		return []string{}
	}
	results := []string{}
	collectWords(node, prefix, &results)
	return results
}

// traverse walks the trie following the characters in s and returns the node
// reached, or nil if no such path exists.
func (t *Trie) traverse(s string) *TrieNode {
	node := t.root
	for _, ch := range s {
		next, ok := node.children[ch]
		if !ok {
			return nil
		}
		node = next
	}
	return node
}

// collectWords performs a depth-first traversal from node, accumulating
// complete words into results. current holds the prefix string that led to
// node.
func collectWords(node *TrieNode, current string, results *[]string) {
	if node.isEndOfWord {
		*results = append(*results, current)
	}
	// Collect and sort the children keys for deterministic (lexicographic) output.
	keys := make([]rune, 0, len(node.children))
	for k := range node.children {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, ch := range keys {
		collectWords(node.children[ch], current+string(ch), results)
	}
}

