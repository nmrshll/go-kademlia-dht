package dht

import "testing"

func TestRoutingTable(t *testing.T) {
	n1 := NewNodeID("FFFFFFFF00000000000000000000000000000000")
	n2 := NewNodeID("FFFFFFF000000000000000000000000000000000")
	n3 := NewNodeID("1111111100000000000000000000000000000000")
	rt := NewRoutingTable(&Node{n1, "localhost:8000"})
	rt.Update(&Node{n2, "localhost:8001"})
	rt.Update(&Node{n3, "localhost:8002"})

	// nodeRecords := rt.FindKClosest(NewNodeID("2222222200000000000000000000000000000000"))
	// if nodeRecords.Len() != 1 {
	// 	t.Fail()
	// 	return
	// }
	// if !nodeRecords[0].node.id.Equals(n3) {
	// 	t.Error(nodeRecords[0])
	// }

	// nodeRecords = rt.FindKClosest(n2)
	// if nodeRecords.Len() != 2 {
	// 	t.Error(nodeRecords.Len())
	// 	return
	// }
	// if !nodeRecords[0].node.id.Equals(n2) {
	// 	t.Error(nodeRecords[0])
	// }
	// if !nodeRecords[1].node.id.Equals(n3) {
	// 	t.Error(nodeRecords[1])
	// }
}
