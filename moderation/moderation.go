package moderation

// todo надо бы сравнить бенчмарк с поиском через регулярку
// todo надо полный список матерных слов
var badWords = []string{"хуй", "хуета", "хуев", "пиписьк", "жопа", "срака", "мудак", "мудила", "пизд", "ебать", "ёбнут", "ебнут", "ебстись", "ебанут", "блядь", "шлюх", "блядск", "fuck", "asshole"}
var root *trieNode

func init() {
	root = buildTrie(badWords)
	buildFailureLinks(root)
}

type trieNode struct {
	children map[rune]*trieNode
	fail     *trieNode
	output   []string
}

func newTrieNode() *trieNode {
	return &trieNode{
		children: make(map[rune]*trieNode),
		fail:     nil,
		output:   []string{},
	}
}

func buildTrie(keywords []string) *trieNode {
	result := newTrieNode()
	for _, keyword := range keywords {
		current := result
		for _, char := range keyword {
			if _, exists := current.children[char]; !exists {
				current.children[char] = newTrieNode()
			}
			current = current.children[char]
		}
		current.output = append(current.output, keyword)
	}
	return result
}

func buildFailureLinks(root *trieNode) {
	var queue []*trieNode
	for _, child := range root.children {
		child.fail = root
		queue = append(queue, child)
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for char, child := range current.children {
			queue = append(queue, child)
			failState := current.fail
			for failState != nil {
				if failChild, exists := failState.children[char]; exists {
					child.fail = failChild
					break
				}
				failState = failState.fail
			}
			if child.fail == nil {
				child.fail = root
			}
			child.output = append(child.output, child.fail.output...)
		}
	}
}

func ahoCorasickSearch(text string, root *trieNode) []string {
	var results []string
	current := root

	for _, char := range text {
		for current != nil && current.children[char] == nil {
			current = current.fail
		}
		if current == nil {
			current = root
		} else {
			current = current.children[char]
		}
		for _, keyword := range current.output {
			results = append(results, keyword)
		}
	}

	return results
}

func SearchBadWord(text string) []string {
	return ahoCorasickSearch(text, root)
}
