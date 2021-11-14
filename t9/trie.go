package t9

type trie struct {
	root *branch
}

func newTrie() *trie {
	return &trie{
		root: newBranch('0'),
	}
}

func (t *trie) addWord(key []rune, word string) {
	t._getBranch(t.root, key, true).addWord(word)
}

func (t *trie) getWords(key []rune, exact bool) []string {
	exactMatchBranch := t._getBranch(t.root, key, false)
	if exactMatchBranch == nil {
		return nil
	} else if exact {
		return exactMatchBranch.getWords()
	}

	// prefix match then
	return t._getDescendantWords(exactMatchBranch)
}

func (t *trie) _getBranch(b *branch, key []rune, addIfNecessary bool) *branch {
	if len(key) == 0 {
		return b
	}

	nextDigit, nextKey := key[0], key[1:]

	var nextBranch *branch
	if addIfNecessary {
		nextBranch = b.getOrAddChild(nextDigit)
	} else {
		nextBranch = b.getChild(nextDigit)
	}

	if nextBranch == nil {
		return nil
	}

	return t._getBranch(nextBranch, nextKey, addIfNecessary)
}

func (t *trie) _getDescendantWords(b *branch) []string {
	words := b.getWords()

	for _, child := range b.getChildren() {
		words = append(words, t._getDescendantWords(child)...)
	}

	return words
}
