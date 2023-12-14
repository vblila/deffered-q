package trie

import "dq/linkedlist"

type Trie struct {
	rootNode *Node
	length   uint32
}

func (t *Trie) Length() uint32 {
	return t.length
}

func (t *Trie) Get(key []rune) *Node {
	levelNode := t.rootNode
	if levelNode == nil {
		return nil
	}

	for i, k := range key {
		if levelNode.nodes == nil {
			return nil
		}

		var lastListNode *linkedlist.Node

		for {
			listNode := levelNode.nodes.Next(lastListNode)
			if listNode == nil {
				break
			}

			trieNode := listNode.Value.(*Node)

			if trieNode.key == k {
				if i == len(key)-1 {
					return trieNode
				} else {
					levelNode = trieNode
					break
				}
			}

			lastListNode = listNode
		}
	}

	return nil
}

func (t *Trie) Put(key []rune, value interface{}) {
	if t.rootNode == nil {
		t.rootNode = &Node{nodes: &linkedlist.List{}}
	}

	levelNode := t.rootNode

	for i, k := range key {
		var nextLevelNode *Node

		var lastListNode *linkedlist.Node
		for {
			listNode := levelNode.nodes.Next(lastListNode)
			if listNode == nil {
				break
			}

			trieNode := listNode.Value.(*Node)

			if trieNode.key == k {
				if i == len(key)-1 {
					trieNode.Value = value
					return
				}

				nextLevelNode = trieNode
				break
			}

			lastListNode = listNode
		}

		if nextLevelNode != nil {
			levelNode = nextLevelNode
		} else {
			nextLevelNode = &Node{key: k, Value: nil, nodes: &linkedlist.List{}}
			levelNode.nodes.Push(nextLevelNode)

			if i == len(key)-1 {
				t.length++
				nextLevelNode.Value = value
				return
			}

			levelNode = nextLevelNode
		}
	}
}

func (t *Trie) Delete(key []rune) bool {
	levelNode := t.rootNode
	if levelNode == nil {
		return false
	}

	for i, k := range key {
		if levelNode.nodes == nil {
			return false
		}

		var lastListNode *linkedlist.Node
		for {
			listNode := levelNode.nodes.Next(lastListNode)
			if listNode == nil {
				break
			}

			trieNode := listNode.Value.(*Node)

			if trieNode.key == k {
				if i == len(key)-1 {
					levelNode.nodes.Delete(listNode)
					t.length--
					return true
				} else {
					levelNode = trieNode
					break
				}
			}

			lastListNode = listNode
		}
	}

	return false
}
