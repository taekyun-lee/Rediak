package KangDB

import "sync"

type instance struct{
	mu sync.Mutex

	b Bucket


	// Persistence
	persistInstance ZippedSnapshot
	snapshotfolder string
	snapshotname string


	//Metadata of instance
	ipaddr string
	portnum uint

}




