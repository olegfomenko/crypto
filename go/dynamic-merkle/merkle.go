package main

import (
	"bytes"

	"github.com/ethereum/go-ethereum/crypto"
)

type Node struct {
	Hash        []byte
	Priority    uint64
	MerkleHash  []byte
	Left, Right *Node
}

type ITreap interface {
	Remove(key []byte)
	Insert(key []byte, priority uint64)
	MerklePath(key []byte) [][]byte
	MerkleRoot() []byte
}

type Treap struct {
	Root *Node
}

// Implements ITreap
var _ ITreap = &Treap{}

func New() ITreap {
	return &Treap{}
}

func (t *Treap) Remove(key []byte) {
	if t.Root == nil {
		return
	}

	t1, t2 := split(t.Root, key)

	if bytes.Compare(t2.Hash, key) == 0 {
		t.Root = merge(t1, t2.Right)
		return
	}

	node := t2
	for {
		if bytes.Compare(node.Left.Hash, key) == 0 {
			node.Left = node
			updateNode(node)
			break
		}

		node = node.Left
	}

	t.Root = merge(t1, t2)
}

func (t *Treap) Insert(key []byte, priority uint64) {
	node := &Node{
		Hash:       key,
		MerkleHash: key,
		Priority:   priority,
	}

	if t.Root == nil {
		t.Root = node
		return
	}

	t1, t2 := split(t.Root, key)
	t.Root = merge(merge(t1, node), t2)
}

func (t *Treap) MerklePath(key []byte) [][]byte {
	node := t.Root
	result := make([][]byte, 0, 64)

	for node != nil {
		if bytes.Compare(node.Hash, key) == 0 {
			result = append(result, hashNodes(node.Left, node.Right))
			return result
		}

		if bytes.Compare(node.Hash, key) > 0 {
			result = append(result, node.Hash)
			if node.Right != nil {
				result = append(result, node.Right.MerkleHash)
			}
			node = node.Left
			continue
		}

		result = append(result, node.Hash)
		if node.Left != nil {
			result = append(result, node.Left.MerkleHash)
		}
		node = node.Right
	}

	return nil
}

func (t *Treap) MerkleRoot() []byte {
	return t.Root.MerkleHash
}

func split(root *Node, key []byte) (*Node, *Node) {
	if root == nil {
		return nil, nil
	}

	if bytes.Compare(root.Hash, key) < 0 {
		t1, t2 := split(root.Right, key)
		root.Right = t1
		updateNode(root)
		return root, t2
	}

	t1, t2 := split(root.Left, key)
	root.Left = t2
	updateNode(root)
	return t1, root
}

func merge(t1, t2 *Node) *Node {
	if t1 == nil {
		return t2
	}

	if t2 == nil {
		return t1
	}

	if t1.Priority > t2.Priority {
		t1.Right = merge(t1.Right, t2)
		updateNode(t1)
		return t1
	}

	t2.Left = merge(t1, t2.Left)
	updateNode(t2)
	return t2
}

func updateNode(node *Node) {
	childrenHash := hashNodes(node.Left, node.Right)
	if childrenHash == nil {
		node.MerkleHash = node.Hash
		return
	}

	node.MerkleHash = hash(childrenHash, node.Hash)
}

func hashNodes(a, b *Node) []byte {
	var left []byte = nil
	var right []byte = nil

	if a != nil {
		left = a.MerkleHash
	}

	if b != nil {
		right = b.MerkleHash
	}

	return hash(left, right)
}

func hash(a, b []byte) []byte {
	if len(a) == 0 {
		return b
	}

	if len(b) == 0 {
		return a
	}

	if bytes.Compare(a, b) < 0 {
		return crypto.Keccak256(a, b)
	}

	return crypto.Keccak256(b, a)
}
