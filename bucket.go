package dht

// BucketCapacity is the max number of items that fit in one bucket
// this factor is called k in the Kademlia paper
const BucketCapacity = 20

// Bucket is a map and a doubly linked list
// it also keeps track of a reference to the first and last item in the linked list
// new contacts get inserted at the back, and the most recently
type Bucket struct {
	recordsMap map[NodeID]NodeRecord
	first      *NodeRecord
	last       *NodeID
	// TODO: it would probably work just as well keeping track of a map + the last seen
}

func NewBucket() *Bucket {
	return &Bucket{
		recordsMap: make(map[NodeID]NodeRecord),
	}
}

func (b *Bucket) Len() int {
	var i int
	for range b.recordsMap {
		i++
	}
	return i
}

type NodeRecords []NodeRecord

func (records NodeRecords) Len() int           { return len(records) }
func (records NodeRecords) Less(i, j int) bool { return records[i].Less(&records[j]) }
func (records NodeRecords) Swap(i, j int)      { records[i], records[j] = records[j], records[i] }

// NodeRecord holds a reference to a node
// and references to the next and previous one (so that the bucket is a doubly linked list)
type NodeRecord struct {
	node     *Node
	sortKey  distance
	previous *NodeID
	next     *NodeID
}

func (nr *NodeRecord) Less(nr2 *NodeRecord) bool {
	return nr.sortKey.Less(nr2.sortKey.i160)
}

// InsertAtBack inserts a node at the last position of the bucket
func (b *Bucket) InsertAtBack(node *Node) {
	nodeRecord := NodeRecord{
		node: node,
	}
	if b.last != nil {
		if lastRecord, ok := b.recordsMap[*b.last]; ok {
			lastRecord.next = &node.id
			nodeRecord.previous = &lastRecord.node.id
		}
	}
	b.recordsMap[node.id] = nodeRecord
	b.last = &node.id
}

// InsertAtBack inserts a node at the first position of the bucket
func (b *Bucket) InsertAtFront(node *Node) {
	nodeRecord := NodeRecord{
		node: node,
	}
	// if b.first != nil {

	// }
	b.recordsMap[node.id] = nodeRecord
	//TODO: incomplete implementation
}

// return bucket.map(node => record{node, distance(node,target)})
func (bucket *Bucket) mapToDistanceToTarget(targetID NodeID) (rRecords []NodeRecord) {
	for _, bucketRecord := range bucket.recordsMap {
		newRecord := bucketRecord
		newRecord.sortKey = XORDistance(bucketRecord.node.id, targetID)
		rRecords = append(rRecords, newRecord)
	}
	return rRecords
}
