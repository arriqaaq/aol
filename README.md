# aol

A simple immutable append only log in Go. aol is a rewrite of [wal](https://github.com/tidwall/wal) to allow only appends to a log.

## Features

- High durability
- Fast operations
- Batch writes

## Getting Started

### Installing

To start using `aol`, install Go and run `go get`:

```sh
$ go get -u github.com/arriqaaq/aol
```

This will retrieve the library.

### Example

```go
// open a new log file
log, err := Open("mylog", nil)

// write some entries
err = log.Write([]byte("first entry"))
err = log.Write([]byte("second entry"))
err = log.Write([]byte("third entry"))

// read an entry
data, err := log.Read(1,0)
println(string(data))  // output: first entry

// close the log
err = log.Close()
```

Batch writes:

```go
// write three entries as a batch
batch := new(Batch)
batch.Write(1, []byte("first entry"))
batch.Write(2, []byte("second entry"))
batch.Write(3, []byte("third entry"))

err = log.WriteBatch(batch)
```

### Inspiration

[wal](https://github.com/tidwall/wal) by [@tidwall](http://twitter.com/tidwall)
