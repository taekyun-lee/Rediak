package KangDB

import (
	"sync"
)


type clusters struct {

	hashmap map[string]Nodevalue

	mu sync.Mutex
	rwmu sync.Mutex
	clusterNum uint32



}
