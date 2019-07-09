package KangDB

import (
	"hash/adler32"
	"hash/crc32"
)

type ShardingAlg interface {
	GetLocationbyKey(key string) uint32
	GetRealValue(key string) string
}

type CRC32Shard struct {
	n uint32
}

func (s CRC32Shard) GetLocationbyKey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key)) % s.n
}

func (s CRC32Shard) GetRealValue(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}



type Adler32Shard struct {
	n uint32
}

func (s Adler32Shard) GetLocationbyKey(key string) uint32 {
	return adler32.Checksum([]byte(key)) % s.n
}
func (s Adler32Shard) GetRealValue(key string) uint32 {
	return adler32.Checksum([]byte(key))
}
