// Package t9 provides tools for looking up words that match particular digits based on the T9 algorithm.
package t9

import (
	"errors"
)

// T9 interface enables insertion of any dictionary word and then retrieval based on a prefix or exact-only
// match of t9 digits.
type T9 interface {
	// Inserts a dictionary word into t9 lookup structure.
	InsertWord(word string) error

	// Retrieves any number of dictionary words from t9 lookup structure based on
	// prefix match or exact match, depending on `exact`.
	GetWords(digits string, exact bool) ([]string, error)
}

// New instantiates a new T9 structure.
func New() T9 {
	return &t9{
		trie: newTrie(),
	}
}

type t9 struct {
	trie *trie
}

func (t9 *t9) InsertWord(word string) error {
	if len(word) == 0 {
		return errors.New(`word is empty`)
	}

	key := getDigits(word)
	t9.trie.addWord(key, word)
	return nil
}

func (t9 *t9) GetWords(digits string, exact bool) ([]string, error) {
	if err := CheckDigits(digits); err != nil {
		return nil, err
	}

	key := []rune(digits)
	return t9.trie.getWords(key, exact), nil
}
