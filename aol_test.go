package aol

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
)

func TestReadWrite(t *testing.T) {
	logPath := "tmp/" + strings.Join(strings.Split(t.Name(), "/")[1:], "/")
	l, err := Open(logPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	defer os.RemoveAll("tmp/")

	for i := 1; i <= 10; i++ {
		key := fmt.Sprintf("key_%d", i)
		err = l.Write([]byte(key))
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}
	}

	segs := l.Segments()
	for i := 1; i <= segs; i++ {
		for j := 0; j < 10; j++ {
			key := fmt.Sprintf("key_%d", j+1)
			data, err := l.Read(uint64(i), uint64(j))
			if err != nil {
				t.Fatalf("expected %v, got %v", nil, err)
			}
			if string(data) != key {
				t.Fatalf("expected %s, got %s", key, data)
			}
		}
	}

	// Read -- random access
	for _, i := range rand.Perm(10) {
		index := uint64(i)
		key := fmt.Sprintf("key_%d", index+1)
		data, err := l.Read(1, index)
		if err != nil {
			if err == ErrEOF {
				continue
			}
			t.Fatal(err)
		}
		if key != string(data) {
			t.Fatalf("expected %v, got %v", key, string(data))
		}
	}
}

func TestReadWrite_Close(t *testing.T) {
	logPath := "tmp/" + strings.Join(strings.Split(t.Name(), "/")[1:], "/")
	l, err := Open(logPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	for i := 1; i <= 10; i++ {
		// Write - append next item
		key := fmt.Sprintf("key_%d", i)
		err = l.Write([]byte(key))
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}
	}
	if err := l.Close(); err != nil {
		t.Fatal(err)
	}

	// Reopen file and read data
	p, err := Open(logPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	defer os.RemoveAll("tmp/")

	segs := l.Segments()
	for i := 1; i <= segs; i++ {
		for j := 0; j < 10; j++ {
			key := fmt.Sprintf("key_%d", j+1)
			data, err := p.Read(uint64(i), uint64(j))
			if err != nil {
				t.Fatalf("expected %v, got %v", nil, err)
			}
			if string(data) != key {
				t.Fatalf("expected %s, got %s", key, data)
			}
		}
	}
}

func TestCycle_AOL(t *testing.T) {
	logPath := "tmp/" + strings.Join(strings.Split(t.Name(), "/")[1:], "/")
	l, err := Open(logPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	defer os.RemoveAll("tmp/")
	for i := 1; i <= 100; i++ {
		key := fmt.Sprintf("key_%d", i)
		l.Write([]byte(key))
	}
	l.cycle()
	for i := 101; i <= 200; i++ {
		key := fmt.Sprintf("key_%d", i)
		l.Write([]byte(key))
	}

	segs := l.Segments()
	var lastKey string
	for i := 1; i <= segs; i++ {
		j := 0
		for {
			data, err := l.Read(uint64(i), uint64(j))
			// fmt.Println("---->", i, l.Segments(), string(data), err)
			if err != nil {
				if err == ErrEOF {
					break
				}
				t.Fatalf("expected %v, got %v", nil, err)
			}
			lastKey = string(data)
			j++
		}
	}

	if lastKey != "key_200" {
		t.Fatalf("expected %v, got %v", "key_200", lastKey)

	}
}
