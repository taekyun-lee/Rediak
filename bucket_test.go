package main

import (
	"testing"
	"time"
)

const (
	TESTSHARDNUM =32
)





func TestDBInstance_Set(t *testing.T) {
	anew := New(true, 10*time.Second)
	pnew := New(false,10*time.Second)


	key1exp := Item{
		v:"expired in 1 sec",
		ttl:time.Now().Add(1*time.Second).UnixNano(),
	}
	key5exp := Item{
		v:"expired in 5 sec",
		ttl:time.Now().Add(5*time.Second).UnixNano(),
	}

	keynotexp := Item{
		v:"not exp",
		ttl:0,
	}

	anew.Set("key1exp", "expired in 1 sec",time.Now().Add(1*time.Second).UnixNano())
	anew.Set("key5exp", "expired in 5 sec",time.Now().Add(5*time.Second).UnixNano())
	anew.Set("keynotexp", "not exp",0)

	if hv,_:= anew.bucket.Load("key1exp");hv.(Item).v != key1exp.v{
		t.Errorf("key1exp not equal")
	}

	if hv,_:= anew.bucket.Load("key5exp");hv.(Item).v != key5exp.v{
		t.Errorf("key5exp not equal")
	}

	if hv,_:= anew.bucket.Load("keynotexp");hv.(Item).v != keynotexp.v{
		t.Errorf("keynotexp not equal")
	}


	pnew.Set("key1exp", "expired in 1 sec",time.Now().Add(1*time.Second).UnixNano())
	pnew.Set("key5exp", "expired in 5 sec",time.Now().Add(5*time.Second).UnixNano())
	pnew.Set("keynotexp", "not exp",0)

	if hv,_:= pnew.bucket.Load("key1exp");hv.(Item).v != key1exp.v{
		t.Errorf("key1exp not equal")
	}

	if hv,_:= pnew.bucket.Load("key5exp");hv.(Item).v != key5exp.v{
		t.Errorf("key5exp not equal")
	}

	if hv,_:= pnew.bucket.Load("keynotexp");hv.(Item).v != keynotexp.v{
		t.Errorf("keynotexp not equal")
	}



}


func TestDBInstance_Get(t *testing.T) {

	anew := New(true, 1*time.Second)
	pnew := New(false,1*time.Second)


	//key1exp := Item{
	//	v:"expired in 1 sec",
	//	ttl:time.Now().Add(1*time.Second).UnixNano(),
	//}
	key5exp := Item{
		v:"expired in 5 sec",
		ttl:time.Now().Add(5*time.Second).UnixNano(),
	}

	keynotexp := Item{
		v:"not exp",
		ttl:0,
	}

	anew.Set("key1exp", "expired in 1 sec",time.Now().Add(1*time.Second).UnixNano())
	anew.Set("key5exp", "expired in 5 sec",time.Now().Add(5*time.Second).UnixNano())
	anew.Set("keynotexp", "not exp",0)
	time.Sleep(5*time.Second)

	if _,ok:=anew.Get("key1exp");ok==nil{
		t.Errorf("key1exp not expired w/ error ")
	}

	if v,ok := anew.Get("key5exp"); ok!=nil || v.v!= key5exp.v {
		t.Errorf("key5exp not equal")
	}

	if v,ok := anew.Get("keynotexp"); ok!=nil || v.v!= keynotexp.v {
		t.Errorf("keynotexp not equal")
	}


	pnew.Set("key1exp", "expired in 1 sec",time.Now().Add(1*time.Second).UnixNano())
	pnew.Set("key5exp", "expired in 5 sec",time.Now().Add(5*time.Second).UnixNano())
	pnew.Set("keynotexp", "not exp",0)

	time.Sleep(3*time.Second)

	if _,ok:=pnew.Get("key1exp");ok==nil{
		t.Errorf("key1exp not expired")
	}


	if v,ok := pnew.Get("key5exp"); ok!=nil || v.v!= key5exp.v {
		t.Errorf("key5exp not equal or expired too early")
	}

	if v,ok := pnew.Get("keynotexp"); ok!=nil || v.v!= keynotexp.v {
		t.Errorf("keynotexp not equal or expired too early")
	}


}



func TestDBInstance_Delete(t *testing.T) {
	anew := New(true, 1*time.Second)
	pnew := New(false,1*time.Second)




	anew.Set("key1exp", "expired in 1 sec",time.Now().Add(1*time.Second).UnixNano())
	anew.Set("key5exp", "expired in 5 sec",time.Now().Add(5*time.Second).UnixNano())
	anew.Set("keynotexp", "not exp",0)
	time.Sleep(3*time.Second)

	anew.Delete("key1exp")
	anew.Delete("key5exp")
	anew.Delete("keynotexp")


	if _,ok:=anew.Get("key1exp");ok==nil{
		t.Errorf("key1exp not expired w/ error ")
	}

	if _,ok := anew.Get("key5exp"); ok==nil  {
		t.Errorf("key5exp not equal")
	}

	if _,ok := anew.Get("keynotexp"); ok==nil {
		t.Errorf("keynotexp not equal")
	}


	pnew.Set("key1exp", "expired in 1 sec",time.Now().Add(1*time.Second).UnixNano())
	pnew.Set("key5exp", "expired in 5 sec",time.Now().Add(5*time.Second).UnixNano())
	pnew.Set("keynotexp", "not exp",0)

	time.Sleep(3*time.Second)

	pnew.Delete("key1exp")
	pnew.Delete("key5exp")
	pnew.Delete("keynotexp")


	if _,ok:=pnew.Get("key1exp");ok==nil{
		t.Errorf("key1exp not expired")
	}


	if _,ok := pnew.Get("key5exp"); ok==nil  {
		t.Errorf("key5exp not deleted")
	}

	if _,ok := pnew.Get("keynotexp"); ok==nil {
		t.Errorf("keynotexp not deleted")
	}

}

func TestDBInstance_IsExists(t *testing.T) {
	anew := New(true, 1*time.Second)
	pnew := New(false,1*time.Second)




	anew.Set("key1exp", "expired in 1 sec",time.Now().Add(1*time.Second).UnixNano())
	anew.Set("key5exp", "expired in 5 sec",time.Now().Add(5*time.Second).UnixNano())
	anew.Set("keynotexp", "not exp",0)
	time.Sleep(3*time.Second)

	anew.Delete("key5exp")


	if ok:=anew.IsExists("key1exp");ok{
		t.Errorf("key1exp existsterror w/ error ")
	}

	if ok:=anew.IsExists("key5exp");ok  {
		t.Errorf("key5exp existsterror")
	}

	if ok:=anew.IsExists("keynotexp");!ok {
		t.Errorf("keynotexp existsterror")
	}


	pnew.Set("key1exp", "expired in 1 sec",time.Now().Add(1*time.Second).UnixNano())
	pnew.Set("key5exp", "expired in 5 sec",time.Now().Add(5*time.Second).UnixNano())
	pnew.Set("keynotexp", "not exp",0)

	time.Sleep(3*time.Second)


	pnew.Delete("key5exp")



	if ok:=pnew.IsExists("key1exp");ok{
		t.Errorf("key1exp existsterror w/ error ")
	}

	if ok:=pnew.IsExists("key5exp");ok  {
		t.Errorf("key5exp existsterror")
	}

	if ok:=pnew.IsExists("keynotexp");!ok {
		t.Errorf("keynotexp existsterror")
	}
}


