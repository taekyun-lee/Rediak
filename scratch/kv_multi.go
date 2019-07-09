// package main

// import (
// 	"fmt"
// 	"sync"
// )

// type basicbucket struct {
// 	kv     map[string]string
// 	mu     *sync.RWMutex
// 	rwlock bool
// }

// func INITBUCKET(rwlockable bool) *basicbucket {
// 	v := basicbucket{
// 		kv:     make(map[string]string),
// 		mu:     new(sync.RWMutex),
// 		rwlock: rwlockable,
// 	}
// 	return &v
// }

// func (b *basicbucket) GET(k string) (string, error) {
// 	b.mu.RLock()
// 	val, ok := b.kv[k]
// 	b.mu.RUnlock()

// 	if ok == false {
// 		fmt.Print(k, " is not here\n")
// 		return "", fmt.Errorf("NotExistError")
// 	}
// 	fmt.Println(k, val)
// 	return val, nil
// }

// func (b *basicbucket) SET(k string, v string) error {

// 	if k == "" {
// 		return fmt.Errorf("KeyIsNilError")
// 	}

// 	b.mu.RLock()
// 	_, ok := b.kv[k]
// 	b.mu.RUnlock()
// 	if ok == true {
// 		return fmt.Errorf("KeyAlreadyExistError")
// 	}

// 	b.mu.Lock()
// 	b.kv[k] = v
// 	b.mu.Unlock()

// 	return nil
// }

// func (b *basicbucket) DEL(k string) error {

// 	if k == "" {
// 		return fmt.Errorf("KeyIsNilError")
// 	}

// 	_, e := b.GET(k)
// 	if e == nil {
// 		return fmt.Errorf("KeyAlreadyExistError")
// 	}

// 	b.mu.Lock()
// 	delete(b.kv, k)
// 	b.mu.Unlock()

// 	return nil
// }

// func (b *basicbucket) MODIFY(k string, v string) error {

// 	if k == "" {
// 		return fmt.Errorf("KeyIsNilError")
// 	}

// 	_, e := b.GET(k)

// 	if e != nil {
// 		return fmt.Errorf("KeyNotExistError")
// 	}

// 	b.mu.Lock()
// 	b.kv[k] = v
// 	b.mu.Unlock()

// 	return nil
// }

// func (b *basicbucket) GETALL() error {
// 	b.mu.RLock()
// 	for k, v := range b.kv {
// 		fmt.Println(k, v)
// 	}
// 	b.mu.RUnlock()

// 	return nil
// }

// func main() {
// 	bb := INITBUCKET(true)
// 	bb.SET("1", "a1")
// 	bb.SET("2", "a2")
// 	bb.SET("3", "a3")
// 	bb.GET("1")
// 	bb.GET("2")
// 	bb.GET("3")
// 	bb.GETALL()

// 	bb.MODIFY("1", "b1")

// 	bb.GETALL()

// }
