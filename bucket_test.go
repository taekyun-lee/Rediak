package KangDB

import (
	"testing"
	"time"
)

var (

)

func TestBucket_GetInterval(t *testing.T) {

	bzero := Bucket{
		name:"testname",
		expiringInterval:0,

	}
	rzero := bzero.GetInterval()
	if rzero != 0{
		t.Errorf("Get interval error orig : %d, test: %s\n ",bzero.expiringInterval,rzero.String())
	}

	b := Bucket{
		name:"testname",
		expiringInterval:10*time.Second,

	}
	r := bzero.GetInterval()
	if r !=10*time.Second {
		t.Errorf("Get interval error orig : %d, test: %s\n ",b.expiringInterval,r.String())
	}
}
func TestBucket_CONTAINS(t *testing.T) {
	
}