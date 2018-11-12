package dht

import (
	// "container/vector"
	// "exp/iterable"
	"sort"
)

// IDLengthBits is the NodeID length in bits (160 in the paper)
const IDLengthBits = IDLength * 8

// RoutingTable is the data structure that keeps track of contact nodes in the network
type RoutingTable struct {
	// selfNode is the node that owns this routing table
	selfNode Node
	// each bit of the nodeID
	// buckets    [IDLength * 8]*list.List
	newBuckets [IDLength * 8]Bucket
}

// NewRoutingTable instantiates a new *RoutingTable
func NewRoutingTable(selfNode *Node) *RoutingTable {
	rt := &RoutingTable{
		selfNode: *selfNode,
	}
	// for i := 0; i < IDLength*8; i++ {
	// 	rt.buckets[i] = list.New()
	// }
	for i := 0; i < IDLength*8; i++ {
		rt.newBuckets[i] = *NewBucket()
	}

	return rt
}

func (table *RoutingTable) getBucketByNodeID(nodeID NodeID) *Bucket {
	prefixLength := XORDistance(table.selfNode.id, nodeID).LeadingZeros()
	return &table.newBuckets[prefixLength]
}

func (table *RoutingTable) getBucketByDistanceToTarget(targetID NodeID) *Bucket {
	bucketNum := XORDistance(table.selfNode.id, targetID).LeadingZeros()
	return &table.newBuckets[bucketNum]
}

func (table *RoutingTable) isSelfInBucket(bucket Bucket) bool {
	for contactID := range bucket.recordsMap {
		if contactID == table.selfNode.id {
			return true
		}
	}
	return false
}

func (table *RoutingTable) Update(contact *Node) {
	bucket := table.getBucketByNodeID(contact.id)
	if bucket.Len() <= BucketCapacity {
		bucket.InsertAtFront(contact)
	} else {
		// TODO: Handle insertion when the list is full by evicting old elements if
		// they don't respond to a ping.
	}

	// // find if self is in the destination bucket
	// element := iterable.Find(bucket, func(x interface{}) bool {
	// 	return x.(*Node).id.Equals(table.node.id)
	// })

	// // if the element is not in the bucket yet
	// if element == nil {
	// 	// add it at the front (best / most recent position) if the bucket is not full
	// 	if bucket.Len() <= BucketCapacity {
	// 		bucket.PushFront(contact)
	// 	}
	// } else {
	// 	bucket.MoveToFront(element.(*list.Element))
	// }
}

// // return bucket.map(node => record{node, distance(node,target)})
// func copyToVector(start, end *list.Element, vec *vector.Vector, target NodeID) {
// 	for elt := start; elt != end; elt = elt.Next() {
// 		contact := elt.Value.(*Node)
// 		vec.Push(&NodeRecord{contact, contact.id.Xor(target)})
// 	}
// }

// FindKClosest returns the K (bucketCapacity) closest known nodes to the one that was requested (target)
func (table *RoutingTable) FindKClosest(targetID NodeID) (kClosest NodeRecords) {
	bucketNum := XORDistance(table.selfNode.id, targetID).LeadingZeros()
	bucket := table.newBuckets[bucketNum]
	// append elements from the bucket, but with the distance from the target for sorting
	kClosest = append(kClosest, bucket.mapToDistanceToTarget(targetID)...)

	// if less than k to return, append nearby buckets
	for i := 1; (bucketNum-i >= 0 || bucketNum+i < IDLength*8) && len(kClosest) < BucketCapacity; i++ {
		if bucketNum-i >= 0 {
			bucket = table.newBuckets[bucketNum-i]
			kClosest = append(kClosest, bucket.mapToDistanceToTarget(targetID)...)
		}
		if bucketNum+i < IDLength*8 {
			bucket = table.newBuckets[bucketNum+i]
			kClosest = append(kClosest, bucket.mapToDistanceToTarget(targetID)...)
		}
	}

	// sort then trim if longer than k
	sort.Sort(kClosest)
	return kClosest[:BucketCapacity]
}
