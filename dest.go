package destruct

import (
	"io"
	"runtime"
)

type destructable interface {
	io.Closer
}

func finalizer(d destructable) {
	d.Close()
}

func Register(d destructable) {
	runtime.SetFinalizer(d, finalizer)
}
