package destruct_test

import (
	"fmt"
	"github.com/conclave/destruct"
	"runtime"
	"sync"
	"testing"
	"time"
	"unsafe"
)

var count int = 0
var mutex sync.Mutex

type DummyDestructable struct {
	_ uintptr // should not be an empty struct
}

func (t *DummyDestructable) Close() error {
	fmt.Printf("\t\t\t\tdestructing 0x%X\n", t.ptr())
	mutex.Lock()
	count--
	mutex.Unlock()
	return nil
}

func (t *DummyDestructable) ptr() uintptr {
	return uintptr(unsafe.Pointer(t))
}

func construct() {
	t := &DummyDestructable{}
	mutex.Lock()
	count++
	mutex.Unlock()
	destruct.Register(t)
	fmt.Printf("0x%X constructed\n", t.ptr())
}

func TestDestruct(t *testing.T) {
	for i := 0; i < 10; i++ {
		construct()
		time.Sleep(time.Second)
		runtime.GC()
	}
	time.Sleep(time.Second)
	fmt.Printf("\ndone.\nobjects left: %d\n\n", count)
	if count < 0 {
		t.Fatal("over purged.")
	} else if count > 0 {
		t.Error("partial purged.")
	}
}
