package trie

import "dq/linkedlist"

type Node struct {
	Value interface{}
	key   rune
	nodes *linkedlist.List
}
