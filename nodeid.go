package dht

import (
	"encoding/hex"
	"math/rand"
)

type i160 [IDLength]byte

// Less asserts if i160 is lower than other i160
func (oneI160 i160) Less(otheri160 i160) bool {
	for i := range oneI160 {
		if oneI160[i] != otheri160[i] {
			return oneI160[i] < otheri160[i]
		}
	}
	return false
}

// LeadingZeros returns the number of zero bits that the nodeID starts with.
// Each byte can add up to 8 zeros
func (oneI160 i160) LeadingZeros() (ret int) {
	for i := 0; i < IDLength; i++ {
		for j := 0; j < 8; j++ {
			if (oneI160[i]>>uint8(7-j))&0x1 != 0 {
				return i*8 + j
			}
		}
	}
	return IDLength*8 - 1
}

// Equals asserts if two nodeIDs have the same value
func (oneI160 i160) Equals(otherID NodeID) bool {
	for i := range oneI160 {
		if oneI160[i] != otherID[i] {
			return false
		}
	}
	return true
}

// IDLength
const IDLength = 20

type NodeID i160

// Equals asserts if two nodeIDs have the same value
func (nodeID NodeID) Equals(otherID NodeID) bool {
	for i := range nodeID {
		if nodeID[i] != otherID[i] {
			return false
		}
	}
	return true
}

// LeadingZeros returns the number of zero bits that the nodeID starts with.
// Each byte can add up to 8 zeros
func (nodeID NodeID) LeadingZeros() (ret int) {
	for i := 0; i < IDLength; i++ {
		for j := 0; j < 8; j++ {
			if (nodeID[i]>>uint8(7-j))&0x1 != 0 {
				return i*8 + j
			}
		}
	}
	return IDLength*8 - 1
}

// Less asserts if i160 is lower than other i160
func (nodeID NodeID) Less(otherID NodeID) bool {
	for i := range nodeID {
		if nodeID[i] != otherID[i] {
			return nodeID[i] < otherID[i]
		}
	}
	return false
}

func NewNodeID(data string) (ret NodeID) {
	decoded, err := hex.DecodeString(data)
	if err != nil {
		// TODO: handle error
	}
	for i := 0; i < IDLength; i++ {
		ret[i] = decoded[i]
	}
	return
}

func NewRandomNodeID() (ret NodeID) {
	for i := 0; i < IDLength; i++ {
		ret[i] = uint8(rand.Intn(256))
	}
	return
}

func (nodeID NodeID) String() string {
	return hex.EncodeToString(nodeID[0:IDLength])
}

type distance struct {
	i160
}

// Xor computes the XOR distance between nodeIDs
func (nodeID NodeID) Xor(otherID NodeID) (d distance) {
	d160 := &d.i160
	for i := range nodeID {
		d160[i] = nodeID[i] ^ otherID[i]
	}
	return
}

func XORDistance(id1, id2 NodeID) (d distance) {
	return id1.Xor(id2)
}

// // LowerThan asserts if nodeID is lower than otherID
// func (nodeID NodeID) Less(otherID NodeID) bool {
// 	for i := 0; i < IDLength; i++ {
// 		if nodeID[i] != otherID[i] {
// 			return nodeID[i] < otherID[i]
// 		}
// 	}
// 	return false
// }
