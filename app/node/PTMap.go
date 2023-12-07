package node

type mNode struct {
	key   rune
	value interface{}
	nodes []*mNode
}

// PTMap Префиксное дерево
type PTMap struct {
	rootNode *mNode
	length   uint64
}

func (ptmap *PTMap) Length() uint64 {
	return ptmap.length
}

func (ptmap *PTMap) Get(key []rune) interface{} {
	levelNode := ptmap.rootNode
	if levelNode == nil {
		return nil
	}

	for i, k := range key {
		if levelNode.nodes == nil {
			return nil
		}

		for _, node := range levelNode.nodes {
			if node.key == k {
				if i == len(key)-1 {
					return node.value
				} else {
					levelNode = node
					break
				}
			}
		}
	}

	return nil
}

func (ptmap *PTMap) Put(key []rune, value interface{}) {
	if ptmap.rootNode == nil {
		ptmap.rootNode = &mNode{}
	}

	levelNode := ptmap.rootNode

	for i, k := range key {
		if levelNode.nodes == nil {
			levelNode.nodes = make([]*mNode, 0)
		}

		var nextLevelNode *mNode

		for _, node := range levelNode.nodes {
			if node.key == k {
				if i == len(key)-1 {
					node.value = value
					return
				}

				nextLevelNode = node
				break
			}
		}

		if nextLevelNode != nil {
			levelNode = nextLevelNode
		} else {
			nextLevelNode = &mNode{key: k, value: nil}
			levelNode.nodes = append(levelNode.nodes, nextLevelNode)

			if i == len(key)-1 {
				ptmap.length++
				nextLevelNode.value = value
				return
			}

			levelNode = nextLevelNode
		}
	}
}

func (ptmap *PTMap) Delete(key []rune) bool {
	levelNode := ptmap.rootNode
	if levelNode == nil {
		return false
	}

	for i, k := range key {
		if levelNode.nodes == nil {
			return false
		}

		for pos, node := range levelNode.nodes {
			if node.key == k {
				if i == len(key)-1 {
					if pos == 0 {
						levelNode.nodes = levelNode.nodes[1:]
					} else if pos == len(levelNode.nodes)-1 {
						levelNode.nodes = levelNode.nodes[:pos]
					} else {
						levelNode.nodes = append(levelNode.nodes[:pos], levelNode.nodes[pos+1:]...)
					}

					ptmap.length--

					return true
				} else {
					levelNode = node
					break
				}
			}
		}
	}

	return false
}
