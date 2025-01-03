package main

import (
	"encoding/binary"

	"github.com/stretchr/testify/assert"
)

/*
A node consists of:
1. A fixed-sized header containing the type of the node (leaf node or internal node)
2. The number of keys.
3. A list of pointers to the child nodes. (Used by internal nodes).
4. A list of offsets pointing to each key-value pair.
5. Packed KV pairs.
| type | nkeys | pointers | offsets | key-values
| 2B | 2B | nkeys * 8B | nkeys * 2B | ...

This is the format of the KV pair. Lengths followed by data.
| klen | vlen | key | val |
| 2B | 2B | ... | ... |

*/


type BNode struct {

	data []byte // data of the node, can be flushed to disk
}

const (
	BNODE_NODE = 1
	BNODE_LEAF = 2
)

type BTree struct {
	root uint64
	get func(uint64) BNode
	new func(BNode) uint64
	del func(uint64)
}

const HEADER = 4
const BTREE_PAGE_SIZE = 4096
const BTREE_MAX_KEY_SIZE = 1000
const BTREE_MAX_VALUE_SIZE = 3000

func init() {
	// init the btree
	node1max := HEADER + 8 + 2 + 4 + BTREE_MAX_KEY_SIZE + BTREE_MAX_VALUE_SIZE
	assert.True(nil, node1max <= BTREE_PAGE_SIZE)
}

// function to get the type of a node (leaf or internal)
func (node BNode) btype() uint16 {
	/*
	Remember the format of a node:
	| type | nkeys | pointers | offsets | key-values
	| 2B | 2B | nkeys * 8B | nkeys * 2B | ...
	So the first two bytes of the node.data[0] and node.data[1] 
	represent the type of the node, each has a size of 1 byte or 8 bits.
	So uint16() will return a 8 bits binary to 16 bits binary by adding 0s to the left.
	The bitwise operation << 8 will shift the bits to the left by 8 bits.
	The bitwise operation | (or) will combine the two 16 bits binary to a 16 bits binary.
	Example:
	node.data[0] = 0x34 (binary 00110100)
	node.data[1] = 0x12 (binary 00010010)
	uint16(node.data[0]): 0x34 (binary 00000000 00110100)
	uint16(node.data[1]) << 8: 0x12 shifted left by 8 bits becomes 0x1200 (binary 00010010 00000000)
	so the result will be 0x1234 (binary 0001001000110100)
	*/

	return uint16(node.data[0]) | uint16(node.data[1])<<8
	// or just use the built in function: return binary.LittleEndian.Uint16(node.data[0:2])
}

func (node BNode) nkeys() uint16 {
	// get the number of keys in the node
	return uint16(node.data[2]) | uint16(node.data[3])<<8
}

func (node *BNode) setHeader(btype uint16, nkeys uint16) {
	// set the header of the node
	// Set btype in the first 2 bytes
	/*
	0xFF is 11111111 in binary, so btype & 0xFF will get the last 8 bits of btype.
	(btype >> 8) & 0xFF will get the first 8 bits of btype.
	Because we are storing in little endian, the first 8 bits of btype (the least significant) should be stored in the second byte of the node.data.
	*/
	node.data[0] = byte(btype & 0xFF)       // Low byte
	node.data[1] = byte((btype >> 8) & 0xFF) // High byte

	// Set nkeys in the next 2 bytes
	node.data[2] = byte(nkeys & 0xFF)       // Low byte
	node.data[3] = byte((nkeys >> 8) & 0xFF) // High byte
}


// pointers

// get the pointer at index i
func (node BNode) getPtr(i uint16) uint64 {
	if (i >= node.nkeys()) {
		return 0
	}
	offset := HEADER + 8*i
	return binary.LittleEndian.Uint64(node.data[offset:offset+8])
	// the above line is equivalent to the following:
	/*

	return 
	uint64(node.data[offset]) | 
	uint64(node.data[offset+1])<<8 | 
	uint64(node.data[offset+2])<<16 | 
	uint64(node.data[offset+3])<<24 | 
	uint64(node.data[offset+4])<<32 | 
	uint64(node.data[offset+5])<<40 | 
	uint64(node.data[offset+6])<<48 | 
	uint64(node.data[offset+7])<<56

	It simply combines the 8 bytes into a 64 bit 	integer in littlle endian order.
	*/
}

func (node *BNode) setPtr(i uint16, ptr uint64) {
	if (i >= node.nkeys()) {
		return
	}
	offset := HEADER + 8*i
	binary.LittleEndian.PutUint64(node.data[offset:offset+8], ptr)
}