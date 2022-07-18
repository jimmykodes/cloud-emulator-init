package atomicerr

import (
	"sync"

	"go.uber.org/multierr"
)

type Error struct {
	errors []error
	mutex  sync.Mutex
}

func (a *Error) Err() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return multierr.Combine(a.errors...)
}

func (a *Error) Append(err error) {
	if err == nil {
		return
	}
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.errors = append(a.errors, err)
}
