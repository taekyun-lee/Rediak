package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/bsm/redeo"
	"github.com/bsm/redeo/resp"
	"github.com/buraksezer/consistent"
	"hash/crc64"
	"io"
	"log"
	"net"



)

type RedisMember struct{
	uuid string
	addr string

	numfailed int


}
var consistentring *consistent.Consistent


var (
	crc64table = crc64.MakeTable(crc64.ECMA)
)
var localAddr *string = flag.String("l", "127.0.0.1:9999", "local address")

type hasher struct{}

func (h hasher) Sum64(data []byte) uint64 {
	// you should use a proper hash function for uniformity.
	return crc64.Checksum(data,crc64table)
}


func init(){

	cfg := consistent.Config{
		PartitionCount:    7,
		ReplicationFactor: 20,
		Load:              1.25,
		Hasher:            hasher{},
	}
	consistentring := consistent.New(nil, cfg)

	node1 := Newmem("127.0.0.1:6380")
	consistentring.Add(node1)

	node2 := Newmem("127.0.0.1:6381")
	consistentring.Add(node2)
	//


}

func  Newmem(addr string) RedisMember{
	uuid , err := newUUID()
	if err !=nil{
		log.Fatal(err)

	}
	return RedisMember{
		uuid:      uuid,
		addr:      addr,
		numfailed: 0,
	}
}

func (m RedisMember)String() string{
	return m.addr
}


// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}


func main() {
	srv := redeo.NewServer(nil)

	// Define handlers
	srv.HandleFunc("ping", func(w resp.ResponseWriter, _ *resp.Command) {
		w.AppendInlineString("PONG")
	})
	srv.HandleFunc("info", func(w resp.ResponseWriter, _ *resp.Command) {
		w.AppendBulkString(srv.Info().String())
	})

	// More handlers; demo usage of redeo.WrapperFunc
	srv.Handle("echo", redeo.WrapperFunc(func(c *resp.Command) interface{} {
		if c.ArgN() != 1 {
			return redeo.ErrWrongNumberOfArgs(c.Name)
		}
		return c.Arg(0)
	}))
	srv.HandleFunc("get", func(w resp.ResponseWriter, c *resp.Command) {

		fmt.Println(c.Args)
		w.AppendInlineString("ok")
	})
	// Open a new listener
	lis, err := net.Listen("tcp", ":9736")
	if err != nil {
		panic(err)
	}
	defer lis.Close()

	// Start serving (blocking)
	srv.Serve(lis)
}