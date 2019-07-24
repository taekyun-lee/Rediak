package main


// Default Command set
// GET/SET/DELETE/REFRESH/CONTAINS for all datatype/interface
// USE WITH CAUTION

// GET key -> value
// SET key value (ttl,optional integer in second)
// DELETE key
// REFRESH key
// CONTAINS key -> boolean of existence

// OPS============================================
// THIS IS OPTION COMMAND FOR EVERY DATA TYPE (Delete,refresh, contain for "entity")
// REFRESH key
// CONTAINS key -> boolean of existence
// OPS============================================

// Primitive value
// string (or byte array, numbers)
// GET key -> value
// SET key value (ttl,optional integer in second)
// DELETE key
// OPS

// bytearray []byte in value
// BGET key -> value
// BSET key value (ttl,optional integer in second)
// OPS

// numbers in value
// NGET key -> value
// NSET key value (ttl,optional integer in second)
// NADD key (value to add, numbers)
// NSUB key (value to subtract, numbers)
// OPS


// GET/SET/DELETE/REFRESH/CONTAINS -> string or []byte
// (cmd of string) + ADD,SUBTRACT for numbers

// List Command set
// List : []string slice of strings

// LGET key -> value
// LSET key value (ttl,optional integer in second) value v1 v2 ...
// OPS
// LGETIDX key idx -> value (of specific string by index)
// LSETIDX key idx value (specific string by index)
// LDELIDX key idx -> delete string by index


// Hashmap command set
// map[string][]byte key = string(strict string), value: byte array

// HGET key -> value
// HSET key value (ttl,optional integer in second)
// OPS
// HGETAUX key auxkey -> value (of specific string by index)
// HSETAUX key auxkey value (specific string by index)
// HDELAUX key auxkey -> delete string by index