package snowman

import "fmt"

var (
	_ Consensus = (*Tree)(nil)
	_ node      = (*unaryNode)(nil)
	_ node      = (*binaryNode)(nil)
)

type Tree struct {
	node
}

func (t *Tree) Initialize(choice int) {
	t.node = &unaryNode{
		tree: t,
	}
}

func (t *Tree) Add(choice int) {

}

type node interface {
	Preference() int
	DecidedPrefix() int
	Add(choice int) node
	// TODO: RecordPoll
	Finalized() bool
	Printable() (string, []node)
}

type unaryNode struct {
	tree *Tree

	child node
}

func (u *unaryNode) Preference() int {
	return -1
}

func (u *unaryNode) DecidedPrefix() int {
	return -1
}

func (u *unaryNode) Add(choice int) node {
	return nil
}

func (u *unaryNode) Finalized() bool {
	return false
}
func (u *unaryNode) Printable() (string, []node) {
	s := fmt.Sprint("unaryNode")
	if u.child == nil {
		return s, nil
	}

	return s, []node{u.child}
}

type binaryNode struct {
	tree        *Tree
	preferences [2]int
	bit         int
	childrens   [2]node
}

func (u *binaryNode) Preference() int {
	return -1
}

func (u *binaryNode) DecidedPrefix() int {
	return -1
}

func (u *binaryNode) Add(choice int) node {
	return nil
}

func (u *binaryNode) Finalized() bool {
	return false
}
func (u *binaryNode) Printable() (string, []node) {
	s := fmt.Sprint("binaryNode")
	if u.childrens[0] == nil {
		return s, nil
	}

	return s, []node{u.childrens[1], u.childrens[0]}
}
