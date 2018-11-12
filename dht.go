package dht

import (
	"sync"

	// "container/vector"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Kademlia struct {
	routes    *RoutingTable
	NetworkID string
}

// NewKademlia instaciates a new Kademlia DHT with the chosen parameters
func NewKademlia(self *Node, networkID string) *Kademlia {
	return &Kademlia{
		routes:    NewRoutingTable(self),
		NetworkID: networkID,
	}
}

func (k *Kademlia) Serve() (err error) {
	rpc.Register(&KademliaCore{k})

	rpc.HandleHTTP()
	if l, err := net.Listen("tcp", k.routes.selfNode.address); err == nil {
		go http.Serve(l, nil)
	}
	return
}

func (k *Kademlia) Call(contact *Node, method string, args, reply interface{}) (err error) {
	if client, err := rpc.DialHTTP("tcp", contact.address); err == nil {
		err = client.Call(method, args, reply)
		if err == nil {
			k.routes.Update(contact)
		}
	}
	return
}

func (k *Kademlia) sendQuery(node *Node, target NodeID, done chan []Node) {
	args := FindNodeRequest{RPCHeader{&k.routes.selfNode, k.NetworkID}, target}
	reply := FindNodeResponse{}

	if err := k.Call(node, "KademliaCore.FindNode", &args, &reply); err == nil {
		done <- reply.contacts
	} else {
		done <- []Node{}
	}
}

type seenNodes struct {
	*sync.Mutex
	seenMap map[NodeID]struct{}
}

func newSeenNodes() *seenNodes {
	return &seenNodes{
		Mutex:   &sync.Mutex{},
		seenMap: make(map[NodeID]struct{}),
	}
}

func (seen seenNodes) Add(nodeID NodeID) {
	seen.Lock()
	defer seen.Unlock()

	seen.seenMap[nodeID] = struct{}{}
}

func (seen seenNodes) alreadySeen(nodeID NodeID) bool {
	seen.Lock()
	defer seen.Unlock()

	_, contains := seen.seenMap[nodeID]
	return contains
}

func (k *Kademlia) IterativeFindNode(target NodeID, delta int) (rRecords []NodeRecord) {
	done := make(chan []Node)

	// A heap of not-yet-queried *Node structs
	// frontier := new(vector.Vector).Resize(0, BucketSize)
	// var nodesToQuery []Node

	// A map of client values we've seen so far
	// seen := make(map[string]struct{})
	seen := newSeenNodes()
	var wg sync.WaitGroup

	// Initialize the return list, frontier heap, and seen list with local nodes
	for _, nodeRecord := range k.routes.FindKClosest(target) {
		// ret.Push(record)
		// heap.Push(frontier, record.node)
		// seen[record.node.id.String()] = true

		seen.Add(nodeRecord.node.id)
		wg.Add(1)
		go k.sendQuery(nodeRecord.node, target, done)
	}
	wg.Wait()

	// // Start off delta queries
	// pending := 0
	// for i := 0; i < delta && frontier.Len() > 0; i++ {
	// 	pending++
	// 	go k.sendQuery(frontier.Pop().(*Node), target, done)
	// }

	// // TODO: use waitgroups instead
	// // Iteratively look for closer nodes
	// for pending > 0 {
	// 	nodes := <-done
	// 	pending--
	// 	for _, node := range nodes {
	// 		// If we haven't seen the node before, add it
	// 		if _, ok := seen[node.id.String()]; ok == false {
	// 			ret.Push(&NodeRecord{&node, node.id.Xor(target)})
	// 			heap.Push(frontier, node)
	// 			seen[node.id.String()] = true
	// 		}
	// 	}

	// 	for pending < delta && frontier.Len() > 0 {
	// 		go k.sendQuery(frontier.Pop().(*Node), target, done)
	// 		pending++
	// 	}
	// }

	// sort.Sort(ret)
	// if ret.Len() > BucketSize {
	// 	ret.Cut(BucketSize, ret.Len())
	// }

	return
}

type RPCHeader struct {
	Sender    *Node
	NetworkID string
}

func (k *Kademlia) HandleRPC(request, response *RPCHeader) error {
	if request.NetworkID != k.NetworkID {
		return fmt.Errorf("Expected network ID %s, got %s",
			k.NetworkID, request.NetworkID)
	}
	if request.Sender != nil {
		k.routes.Update(request.Sender)
	}
	response.Sender = &k.routes.selfNode
	return nil
}

type KademliaCore struct {
	kad *Kademlia
}

type PingRequest struct {
	RPCHeader
}

type PingResponse struct {
	RPCHeader
}

func (kc *KademliaCore) Ping(args *PingRequest, response *PingResponse) error {
	if err := kc.kad.HandleRPC(&args.RPCHeader, &response.RPCHeader); err != nil {
		return err
	}
	log.Printf("Ping from %s\n", args.RPCHeader)
	return nil
}

type FindNodeRequest struct {
	RPCHeader
	target NodeID
}

type FindNodeResponse struct {
	RPCHeader
	contacts []Node
}

func (kc *KademliaCore) FindNode(args *FindNodeRequest, response *FindNodeResponse) (err error) {
	if err = kc.kad.HandleRPC(&args.RPCHeader, &response.RPCHeader); err == nil {
		kClosestContacts := kc.kad.routes.FindKClosest(args.target)

		response.contacts = make([]Node, len(kClosestContacts))
		for _, nodeRecord := range kClosestContacts {
			response.contacts = append(response.contacts, *nodeRecord.node)
		}

		// for i := 0; i < contacts.Len(); i++ {
		// 	response.contacts[i] = *contacts.At(i).(*NodeRecord).node
		// }
	}
	return
}
