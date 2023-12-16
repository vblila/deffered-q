package trie

import "dq/linkedlist"

type Trie struct {
	rootNode *Node
	length   uint32
}

func (t *Trie) Length() uint32 {
	return t.length
}

func (t *Trie) Get(key []rune) interface{} {
	levelNode := t.rootNode
	if levelNode == nil {
		return nil
	}

	for i, k := range key {
		if levelNode.nodes == nil {
			return nil
		}

		isRuneNodeFound := false

		var currListNode *linkedlist.Node
		for {
			currListNode = levelNode.nodes.Next(currListNode)
			if currListNode == nil {
				break
			}

			trieNode := currListNode.Value.(*Node)

			if trieNode.nodeKey == k {
				if i == len(key)-1 {
					return trieNode.Value
				} else {
					isRuneNodeFound = true
					levelNode = trieNode
					break
				}
			}
		}

		if isRuneNodeFound == false {
			break
		}
	}

	return nil
}

func (t *Trie) Put(key []rune, value interface{}) {
	if t.rootNode == nil {
		t.rootNode = &Node{nodes: &linkedlist.List{}}
	}

	runeNode := t.rootNode

	var currRuneNode *Node

	for i, k := range key {
		currRuneNode = nil

		var currListNode *linkedlist.Node
		for {
			currListNode = runeNode.nodes.Next(currListNode)
			if currListNode == nil {
				break
			}

			trieNode := currListNode.Value.(*Node)

			if trieNode.nodeKey == k {
				currRuneNode = trieNode
				break
			}
		}

		if currRuneNode == nil {
			currRuneNode = &Node{nodeKey: k, Value: nil, nodes: &linkedlist.List{}}
			runeNode.nodes.Push(currRuneNode)
		}

		if i == len(key)-1 {
			t.length++
			currRuneNode.Value = value
			return
		}

		runeNode = currRuneNode
	}
}

func (t *Trie) Delete(key []rune) bool {
	return t.deleteNode(key, false)
}

func (t *Trie) deleteNode(key []rune, onlyEmptyNode bool) bool {
	if len(key) == 0 {
		return false
	}

	levelNode := t.rootNode
	if levelNode == nil {
		return false
	}

	for i, k := range key {
		if levelNode.nodes == nil {
			return false
		}

		isRuneNodeFound := false

		var currListNode *linkedlist.Node
		for {
			currListNode = levelNode.nodes.Next(currListNode)
			if currListNode == nil {
				break
			}

			trieNode := currListNode.Value.(*Node)

			if trieNode.nodeKey == k {
				if i == len(key)-1 {
					if onlyEmptyNode == false {
						trieNode.Value = nil

						if trieNode.nodes.Length() == 0 {
							levelNode.nodes.Delete(currListNode)
						}

						t.length--

						t.deleteNode(key[:len(key)-1], true)

						return true
					} else {
						if trieNode.Value == nil && trieNode.nodes.Length() == 0 {
							levelNode.nodes.Delete(currListNode)
							t.deleteNode(key[:len(key)-1], true)
						}

						return false
					}
				} else {
					levelNode = trieNode
					isRuneNodeFound = true
					break
				}
			}
		}

		if isRuneNodeFound == false {
			break
		}
	}

	return false
}
