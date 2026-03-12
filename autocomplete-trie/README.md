# go-autocomplete-trie
this module provides a trie type and methods to be used for last-word autocomplete

# Use
This trie is designed to autocomplete the last word in a string (space delimited) after being initialized with a list of strings.

Each node in the trie contains a map of child nodes and a list of strings which are the completions available at that node.

The contributed `AutocompleteTrie` type has the following method functions:
- `CreateAutocompleteTrie([]string)` returns a pointer to the created trie's root node and an error. The string array is used to populate the trie.
- `t.AddWord(string)` (pointer receiver) updates the trie with the given word, and returns an error
- `t.GetAllCompletions(string)` (pointer receiver) returns a string prefix of the query and a list of strings containing available completions, formatted so that `prefix + completions[i]` will return the desired result
- `t.Autocomplete(string)` (pointer receiver) accepts a query and returns the query with the last word autocompleted to the first completion found at the traversed node. This is intended for use with tries that were initialized with words in order of importance.
- `t.GetNodeAt([]byte)` (pointer receiver) which takes in an array of bytes to use to traverse the trie. It returns a node found at that path or `nil` if it didn't exist.

# Examples
For a trie initialized with the strings "foo", "bar", and "baz", `myTrie.GetNodeAt([]byte{'b','a'})` would return a trie node with node.Values being `["bar", "baz"]` and a map node.Children containing a mapping 'r' and a mapping 'z'

This is why input strings are recommended to be listed in order of importance, because Values[0] would return the best option. 

The trie is designed to autocomplete one word at a time. For example, a trie initialized with the words "ham", "cheese", "sandwich" would be able to autocomplete the following strings:

"h" -> "ham"

"ham an" -> "ham an" (because "and" is not part of the trie)

"ham and c" -> "ham and cheese"

"ham and cr" -> "ham and cr" (because "cr" does not match any words in the trie)

However, words used to initialize the trie or which are added to the trie can contain spaces. 

If you want more control over autocomplete results, use the method GetAllCompletions. This method will return a prefix (empty string for one-word queries) and an array of strings which are the available completions. The prefix is formatted so that prefix + completions[i] will be the autocompleted query. For example:

 if the trie was initialized with ["anaconda","apples","action"], then you could use:
 
`prefix, completions := t.GetAllCompletions("I really love a")`

and `prefix + completions[1]` would be "I really love apples"

<br>

*Note: the trie is case-sensitive due to the nature of 'A' != 'a'*