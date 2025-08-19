package sensitive

import (
	"os"
	"regexp"
	"sync"

	"strings"
)

var (
	sensitiveTrieOnce sync.Once
)

var st *sensitiveTrie

type sensitiveTrie struct {
	replaceChar rune
	root        *trieNode
}

func newSensitiveTrie() *sensitiveTrie {
	return &sensitiveTrie{
		replaceChar: '*',
		root:        &trieNode{End: false},
	}
}

func (st *sensitiveTrie) filterSpecialChar(text string) string {
	text = strings.ToLower(text)
	text = strings.Replace(text, " ", "", -1)
	otherCharReg := regexp.MustCompile("[^\u4e00-\u9fa5a-zA-Z0-9]")
	return otherCharReg.ReplaceAllString(text, "")
}

func (st *sensitiveTrie) addWord(sensitiveWord string) {
	sensitiveWord = st.filterSpecialChar(sensitiveWord)
	tireNode := st.root
	sensitiveChars := []rune(sensitiveWord)
	for _, charInt := range sensitiveChars {
		tireNode = tireNode.addChild(charInt)
	}
	tireNode.End = true
	tireNode.Data = sensitiveWord
}

func (st *sensitiveTrie) addWords(sensitiveWords []string) {
	for _, sensitiveWord := range sensitiveWords {
		st.addWord(sensitiveWord)
	}
}

func (st *sensitiveTrie) match(text string) (sensitiveWords []string, replaceText string) {
	if st.root == nil {
		return nil, text
	}
	filteredText := st.filterSpecialChar(text)
	sensitives := make(map[string]*struct{})
	textChars := []rune(filteredText)
	textCharsCopy := make([]rune, len(textChars))
	copy(textCharsCopy, textChars)
	for i, textLen := 0, len(textChars); i < textLen; i++ {
		trieNode := st.root.findChild(textChars[i])
		if trieNode == nil {
			continue
		}
		j := i + 1
		for ; j < textLen && trieNode != nil; j++ {
			if trieNode.End {
				if _, ok := sensitives[trieNode.Data]; !ok {
					sensitiveWords = append(sensitiveWords, trieNode.Data)
				}
				sensitives[trieNode.Data] = nil
				st.replaceRune(textCharsCopy, i, j)
			}
			trieNode = trieNode.findChild(textChars[j])
		}
		if j == textLen && trieNode != nil && trieNode.End {
			if _, ok := sensitives[trieNode.Data]; !ok {
				sensitiveWords = append(sensitiveWords, trieNode.Data)
			}
			sensitives[trieNode.Data] = nil
			st.replaceRune(textCharsCopy, i, textLen)
		}
	}
	if len(sensitiveWords) > 0 {
		replaceText = string(textCharsCopy)
	} else {
		replaceText = text
	}
	return sensitiveWords, replaceText
}

func (st *sensitiveTrie) replaceRune(chars []rune, begin int, end int) {
	for i := begin; i < end; i++ {
		chars[i] = st.replaceChar
	}
}

type trieNode struct {
	children map[rune]*trieNode
	Data     string
	End      bool
}

func (tn *trieNode) addChild(c rune) *trieNode {
	if tn.children == nil {
		tn.children = make(map[rune]*trieNode)
	}
	if trieNode, ok := tn.children[c]; ok {
		return trieNode
	}
	tn.children[c] = &trieNode{
		children: nil,
		End:      false,
	}
	return tn.children[c]
}

func (tn *trieNode) findChild(c rune) *trieNode {
	if tn.children == nil {
		return nil
	}
	if trieNode, ok := tn.children[c]; ok {
		return trieNode
	}
	return nil
}

func Match(text string) string {
	_, result := st.match(text)
	return result
}

func IsMatch(text string) bool {
	words, _ := st.match(text)
	return len(words) > 0
}

func NewSensitiveTrie(filePath string) {
	sensitiveTrieOnce.Do(func() {
		if filePath == "" {
			// 敏感词文件路径为空
			panic("The sensitive words file path is empty.")
		}
		bytes, err := os.ReadFile(filePath)
		if err != nil {
			panic(err)
		}
		st = newSensitiveTrie()
		contents := strings.Split(string(bytes), "\n")
		st.addWords(contents)
	})
}
