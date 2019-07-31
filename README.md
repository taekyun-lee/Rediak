# Rediak 
- Distribute(work in progress) in-memory key-value store written in Go
- (wanna) fast NoSQL DB uses RESP (REdis Serialization Protocol).
- multi-thread key-value store using goroutine
- Sync.map 
- Making your own custom RESP-based command available

## Overview

-(Distributed) in-memory key-value store 
- Compatible with RESP protocol(some commands only, see below)


## Install
```go
go build -ldflags "-w" *.go
```

## Usage
### Redis compatible use like redis,
- In some commands only though... this project is my summer projects and currently working on


### Current implemented command

- String and keys
    - [x] get 
    - [x] set 
    - [x] del 
    - [x] exists 
    - [x] incr 
    - [x] incrby
    - [x] expire

- Hash
    - [x] hget
    - [x] hset
    - [x] hdel
    - [x] hexists
    
- list
    - [x] lpush
    - [x] lpop
    - [x] lindex
    - [x] llen
    - [x] lrange
    
- custom and etc
    - [x] ping (response :pong)
    - [x] gc (execute GC )
    - [x] exit
    

## Configuration

```bash
./rediak -evict-interval second (some options...)
```

 - Strong-lock
        use mutex to all modification command,  
  - evict-interval int
        Default interval of eviction, 0 means no active eviction
  - num-core int
        number of cores using this instances (default 4)
  - printinfo-interval int
        Default time interval of print interval (default 1)
  - rediak-addr string
        the address of rediak server (default "127.0.0.1") -> currently not used
  - rediak-port int
        the port of rediak server (default 6380)
  - restore-snapshot string
        if restore needed, write absolute path of files( like /path/of/folder/file.rdb ) 
  - snapshot-interval int
        Default time interval of take snapshot, 0 means no active snapshot
  - snapshot-modify-interval int
        Default time interval of take snapshot, 0 means no active eviction (default 1)
  - storage-dir string
        Default persistent storage location /path/to/snapshotfolder (default "./")


## References 

*for basic structure*
- https://github.com/alash3al/redix 
- Apache License 2.0

*for RESP*
- https://github.com/tidwall/redcon
- MIT License

*consistent hashing for cluster*
- https://github.com/buraksezer/consistent
- MIT License

## Contributing
 - Sorry.
 - I want to develop this project on my own this time.
 - I know this project have lots of dirty code.
 
## Issue
- Please make new issue if you have to.



## License
[GNU General Public License v3.0](https://github.com/taekyun-lee/Rediak/blob/master/LICENSE)
> [For references's license, See NOTICE](https://github.com/taekyun-lee/Rediak/blob/master/NOTICE)
> and some files...
