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

	for i := 1; i <= 10; i++ {
		key := fmt.Sprintf("key_%d", i)
		data, err := l.Read(uint64(i))
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}
		if string(data) != key {
			t.Fatalf("expected %s, got %s", key, data)
		}
	}

	// Read -- random access
	for _, i := range rand.Perm(10) {
		index := uint64(i + 1)
		key := fmt.Sprintf("key_%d", index)
		data, err := l.Read(index)
		if err != nil {
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

	for i := 1; i <= 10; i++ {
		key := fmt.Sprintf("key_%d", i)
		data, err := p.Read(uint64(i))
		if err != nil {
			t.Fatalf("expected %v, got %v", nil, err)
		}
		if string(data) != key {
			t.Fatalf("expected %s, got %s", key, data)
		}
	}
}
