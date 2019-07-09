package main

import (
	"crypto/rand"
	"fmt"
	"sync"
)

type msbucket struct {
	name      string
	kv        map[string]string
	mu        *sync.RWMutex
	rwlock    bool
	isSlave   bool
	slavelist []msbucket
	slavenum  int
}

func INITMSBUCKET(name string, rwlockable, slavecheck bool) *msbucket {
	v := msbucket{
		name:      name,
		kv:        make(map[string]string),
		mu:        new(sync.RWMutex),
		rwlock:    rwlockable,
		isSlave:   slavecheck,
		slavelist: make([]msbucket, 5),
		slavenum:  0,
	}
	return &v
}

func (b *msbucket) GET(k string) (string, error) {
	fmt.Println("name ", b.name, "GET with key = ", k)
	b.mu.RLock()
	randnode := rand.Int(b.slavenum) - 1
	if randnode == -1 {
		val, ok := b.kv[k]

	} else {
		val, ok := b.slavelist[randnode].kv[k]
	}
	b.mu.RUnlock()

	if ok == false {
		fmt.Print(k, " is not here\n")
		return "", fmt.Errorf("NotExistError")
	}
	fmt.Println(k, val)
	return val, nil
}

func (b *msbucket) SET(k string, v string) error {
	fmt.Println("name ", b.name, "SET with key = ", k, "and Value = ", v)

	if b.isSlave {
		return fmt.Errorf("SlaveInstanceError")
	}

	if k == "" {
		return fmt.Errorf("KeyIsNilError")
	}

	b.mu.RLock()
	_, ok := b.kv[k]
	b.mu.RUnlock()
	if ok == true {
		return fmt.Errorf("KeyAlreadyExistError")
	}

	b.mu.Lock()
	b.kv[k] = v
	b.mu.Unlock()

	return nil
}

func (b *msbucket) DEL(k string) error {
	fmt.Println("name ", b.name, "DEL with key = ")

	if k == "" {
		return fmt.Errorf("KeyIsNilError")
	}

	_, e := b.GET(k)
	if e == nil {
		return fmt.Errorf("KeyAlreadyExistError")
	}

	b.mu.Lock()
	delete(b.kv, k)
	b.mu.Unlock()

	return nil
}

func (b *msbucket) MODIFY(k string, v string) error {
	fmt.Println("name ", b.name, "MODIFY with key = ", k, "and Value = ", v)

	if k == "" {
		return fmt.Errorf("KeyIsNilError")
	}

	_, e := b.GET(k)

	if e != nil {
		return fmt.Errorf("KeyNotExistError")
	}

	b.mu.Lock()
	b.kv[k] = v
	b.mu.Unlock()

	return nil
}

func (b *msbucket) GETALL() error {
	fmt.Println("name ", b.name, "GET ALL ")

	b.mu.RLock()
	for k, v := range b.kv {
		fmt.Println(k, v)
	}
	b.mu.RUnlock()

	return nil
}

func main() {
	bb := INITMSBUCKET("master", true, false)
	for i := 1; i <= 3; i++ {
		bb.slavelist = append(bb.slavelist, *INITMSBUCKET(fmt.Sprintf("slave:%d", i), true, true))
		bb.slavenum++
	}

	bb.SET("1", "a1")
	bb.SET("2", "a2")
	bb.SET("3", "a3")
	bb.GET("1")
	bb.GET("2")
	bb.GET("3")
	bb.GETALL()

	bb.MODIFY("1", "b1")

	bb.GETALL()

}
