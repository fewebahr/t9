package t9

import (
	"sync"
)

type branch struct {
	digit rune

	childrenSlice []*branch
	childrenMap   map[rune]*branch
	childrenMutex *sync.RWMutex

	words      []string
	wordsMutex *sync.RWMutex
}

func newBranch(digit rune, words ...string) *branch {

	return &branch{
		digit:         digit,
		childrenSlice: make([]*branch, 0),
		childrenMap:   make(map[rune]*branch),
		childrenMutex: new(sync.RWMutex),
		words:         words,
		wordsMutex:    new(sync.RWMutex),
	}
}

func (b *branch) addChild(child *branch) {

	b.childrenMutex.Lock()
	b.childrenSlice = append(b.childrenSlice, child)
	b.childrenMap[child.digit] = child
	b.childrenMutex.Unlock()
}

func (b *branch) getChildren() []*branch {

	b.childrenMutex.RLock()
	children := b.childrenSlice
	b.childrenMutex.RUnlock()

	return children
}

func (b *branch) getChild(digit rune) *branch {

	b.childrenMutex.RLock()
	child := b.childrenMap[digit]
	b.childrenMutex.RUnlock()

	return child
}

func (b *branch) getOrAddChild(digit rune) *branch {

	b.childrenMutex.Lock()
	child := b.childrenMap[digit]
	if child == nil {
		child = newBranch(digit)
		b.childrenSlice = append(b.childrenSlice, child)
		b.childrenMap[digit] = child
	}
	b.childrenMutex.Unlock()

	return child
}

func (b *branch) addWord(word string) {

	b.wordsMutex.Lock()
	b.words = append(b.words, word)
	b.wordsMutex.Unlock()
}

func (b *branch) getWords() []string {

	b.wordsMutex.RLock()
	words := b.words
	b.wordsMutex.RUnlock()

	return words
}
