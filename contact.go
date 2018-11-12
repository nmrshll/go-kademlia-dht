package dht

import (
	"fmt"
)

type Node struct {
	id      NodeID
	address string
}

func (contact *Node) String() string {
	return fmt.Sprintf("Node(\"%s\", \"%s\")", contact.id, contact.address)
}
