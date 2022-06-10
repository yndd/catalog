package catalog

import (
	"errors"
	"sync"

	"github.com/yndd/ndd-runtime/pkg/resource"
	targetv1 "github.com/yndd/target/apis/target/v1"
)

type Catalog interface {
	GetFn(fnKey FnKey) (Fn, error)
	RegisterFn(fnKey FnKey, fn Fn)
}

type Fn func(in *Input) (resource.Managed, error)

type Input struct {
	Meta   map[string]interface{}
	Object interface{}
	Data   interface{}
}

type FnKey struct {
	Name    string
	Version string
	//
	Vendor    targetv1.VendorType
	Platform  string
	SwVersion string
}

type catalog struct {
	m   *sync.RWMutex
	fns map[FnKey]Fn
}

var Default = &catalog{
	m:   &sync.RWMutex{},
	fns: map[FnKey]Fn{},
}

func (c *catalog) GetFn(key FnKey) (Fn, error) {
	c.m.RLock()
	defer c.m.RUnlock()
	if f, ok := c.fns[key]; ok {
		return f, nil
	}
	return nil, errors.New("not found")
}

func (c *catalog) RegisterFn(key FnKey, fn Fn) {
	c.m.Lock()
	defer c.m.Unlock()
	c.fns[key] = fn
}

func RegisterFns(c Catalog, fns map[FnKey]Fn) {
	for k, v := range fns {
		c.RegisterFn(k, v)
	}
}
