package main
//
//import (
//	"fmt"
//	"hash/crc64"
//	"math/rand"
//	"time"
//
//	"github.com/buraksezer/consistent"
//
//)
//
//// In your code, you probably have a custom data type
//// for your cluster members. Just add a String function to implement
//// consistent.Member interface.
//type Member string
//
//func (m Member) String() string {
//	return string(m)
//}
//
//// consistent package doesn't provide a default hashing function.
//// You should provide a proper one to distribute keys/members uniformly.
//type hasher struct{}
//
//func (h hasher) Sum64(data []byte) uint64 {
//	// you should use a proper hash function for uniformity.
//	return crc64.Checksum(data,crc64.MakeTable(crc64.ECMA))
//}
//func init() {
//	rand.Seed(time.Now().UTC().UnixNano())
//}
//
//func main() {
//	// Create a new consistent instance.
//	members := []consistent.Member{}
//	for i := 0; i < 8; i++ {
//		member := Member(fmt.Sprintf("node%d.olricmq", i))
//		members = append(members, member)
//	}
//	// Modify PartitionCount, ReplicationFactor and Load to increase or decrease
//	// relocation ratio.
//	cfg := consistent.Config{
//		PartitionCount:    271,
//		ReplicationFactor: 20,
//		Load:              1.25,
//		Hasher:            hasher{},
//	}
//	c := consistent.New(members, cfg)
//
//	// Store current layout of partitions
//	owners := make(map[int]string)
//	for partID := 0; partID < cfg.PartitionCount; partID++ {
//		owners[partID] = c.GetPartitionOwner(partID).String()
//	}
//
//	// Add a new member
//	m := Member(fmt.Sprintf("node%d.olricmq", 9))
//	c.Add(m)
//
//	// Get the new layout and compare with the previous
//	var changed int
//	for partID, member := range owners {
//		owner := c.GetPartitionOwner(partID)
//		if member != owner.String() {
//			changed++
//			fmt.Printf("partID: %3d moved to %s from %s\n", partID, owner.String(), member)
//		}
//	}
//	fmt.Printf("\n%d%% of the partitions are relocated\n", (100*changed)/cfg.PartitionCount)
//}