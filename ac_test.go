package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"
)

type ACAutoMation struct {
	root *TrieNode
}

type TrieNode struct {
	children map[rune]*TrieNode

	length int

	fail *TrieNode
}

func NewACAutoMation() *ACAutoMation {
	return &ACAutoMation{
		root: &TrieNode{},
	}
}

func (am *ACAutoMation) BuildTrie(rd io.Reader) error {
	bfRd := bufio.NewReader(rd)

	for {
		line, prefix, err := bfRd.ReadLine()
		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}
		if prefix {
			return fmt.Errorf(`The word is too long, start by "%s"`, line)
		}
		am.addWord(string(line))
	}
}

func (am *ACAutoMation) BuildFailNode() {
	queue := []*TrieNode{}
	for _, node := range am.root.children {
		node.fail = am.root
		queue = append(queue, node)
	}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		for r, node := range cur.children {
			fail := cur.fail
			queue = append(queue, node)
			for {
				if fail == nil {
					node.fail = am.root
					break
				}
				if failChild := fail.children[r]; failChild != nil {
					node.fail = fail.children[r]
					break
				} else {
					fail = fail.fail
				}
			}
		}
	}
}

func (am *ACAutoMation) addWord(word string) error {
	cur := am.root
	if cur.children == nil {
		cur.children = make(map[rune]*TrieNode)
	}
	length := 0
	for _, r := range word {
		if next := cur.children[r]; next != nil {
			cur = next
		} else {
			next = &TrieNode{
				children: make(map[rune]*TrieNode),
			}
			cur.children[r] = next
			cur = next
		}
		length++
	}
	cur.length = length
	return nil
}

type DirtyKey struct {
	endIndex int
	length   int
}

func (am *ACAutoMation) filter(content string) string {
	keys := []DirtyKey{}
	cur := am.root
	rs := []rune(content)
	for i := 0; i < len(rs); {
		addI := false
		if next := cur.children[rs[i]]; next != nil {
			cur = next
			addI = true
		} else if cur != am.root {
			cur = cur.fail
		} else {
			addI = true
		}
		ccur := cur
		for ccur != am.root {
			if ccur.length > 0 {
				keys = append(keys, DirtyKey{i, ccur.length})
				break
			}
			ccur = ccur.fail
		}
		if addI {
			i++
		}
	}
	// 填充**
	for _, key := range keys {
		fmt.Println(key)
		for i := key.endIndex - key.length + 1; i <= key.endIndex; i++ {
			rs[i] = rune('*')
		}
	}
	return string(rs)
}

func TestAC(t *testing.T) {
	ac := NewACAutoMation()
	f, err := os.Open("words.txt")
	if err != nil {
		t.Fatal(err)
	}
	if err := ac.BuildTrie(f); err != nil {
		t.Fatal(err)
	}
	ac.addWord("x")
	ac.BuildFailNode()

	t.Log(ac.filter("办理各种本科 罢工门的事 xxxx"))
}
