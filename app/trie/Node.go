package trie

import "dq/linkedlist"

type Node struct {
	Value   interface{}
	nodeKey rune
	nodes   *linkedlist.List
}
