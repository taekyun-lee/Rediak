package KangDB

import (
	"fmt"
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
	r := b.GetInterval()
	if r !=10*time.Second {
		t.Errorf("Get interval error orig : %d, test: %s\n ",b.expiringInterval,r.String())
	}
}
func TestBucket_SET(t *testing.T) {
	bzero := Bucket{
		name:"testname",
		expiringInterval:0,

	}
	//No exp, kv exists

	bzero.SET("1","starbucks",0,1)

	v,ok:= bzero.items.Load("1")
	if ok != true {
		t.Errorf("Set failed all ")
	}

	va:= v.(item)
	fmt.Println(va)

	if va.ttl != 0 || va.dtype != 1 || va.value !="starbucks"{
		t.Errorf("Set value failed  ")
	}

	bone := Bucket{
		name:"testname",
		expiringInterval:5,

	}
	//Has exp, kv exists

	bone.SET("1","starbucks",3,1)

	v,ok = bzero.items.Load("1")
	if ok != true {
		t.Errorf("Set failed all ")
	}

	va = v.(item)
	fmt.Println(va)

	if va.ttl != 0 || va.dtype != 1 || va.value !="starbucks"{
		t.Errorf("Set value failed  ")
	}


}

func TestBucket_GET(t *testing.T) {
	//no exp
	b := Bucket{
		name:"testname",
		expiringInterval:0,
	}

	stritem := item{
		value:"asdfasdfasdf",
		ttl:time.Now().Add(3*time.Second).UnixNano(),
		dtype:1,
	}

	listitem := item{
		value:[]string{"asdf","qwer","zxcv"},
		ttl:time.Now().Add(3*time.Second).UnixNano(),
		dtype:1,
	}

	b.items.Store("1",stritem)
	b.items.Store("2",listitem)


	strget , strok := b.GET("1")
	listget , listok := b.GET("1")


	if (strok != nil) || (strget.(item) != stritem) {
		t.Errorf("get failed all ")
	}
	if (listok != nil) || (listget.(item).value != stritem.value) {
		t.Errorf("get failed all ")
	}

	time.Sleep(3*time.Second)



	strget , strok = b.GET("1")
	listget , listok = b.GET("1")


	if (strok == nil)  {
		t.Errorf("get passive expiration failed ")
	}
	if (listok == nil) {
		t.Errorf("get passive expiration failed ")
	}





//active Expiration

	//no exp
	b = Bucket{
		name:"testname",
		expiringInterval:5*time.Second,
	}

	stritem = item{
		value:"asdfasdfasdf",
		ttl:time.Now().Add(3*time.Second).UnixNano(),
		dtype:1,
	}

	listitem = item{
		value:[]string{"asdf","qwer","zxcv"},
		ttl:time.Now().Add(3*time.Second).UnixNano(),
		dtype:1,
	}

	b.items.Store("1",stritem)
	b.items.Store("2",listitem)


	strget , strok = b.GET("1")
	listget , listok = b.GET("1")


	if (strok != nil) || (strget.(item) != stritem) {
		t.Errorf("get failed all ")
	}
	if (listok != nil) || (listget.(item).value != stritem.value) {
		t.Errorf("get failed all ")
	}

	time.Sleep(5*time.Second)



	strget , strok = b.GET("1")
	listget , listok = b.GET("1")


	if (strok == nil)  {
		t.Errorf("get active expiration failed ")
	}
	if (listok == nil) {
		t.Errorf("get active expiration failed ")
	}


}