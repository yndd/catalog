package catalog

import (
	"errors"
	"fmt"
	"sync"

	"github.com/yndd/ndd-runtime/pkg/resource"
	targetv1 "github.com/yndd/target/apis/target/v1"
	topologyv1alpha1 "github.com/yndd/topology/apis/topo/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Catalog interface {
	Get(key Key) (Entry, error)
	Register(key Key, e Entry)
	Merge(c ...Catalog)
	List() []Key
}

type GetGvkKeyFn func(key Key, in *Input) (string, Key, error)

type Fn func(key Key, in *Input) (resource.Managed, error)

type MergeFn func(crs ...resource.Managed) (resource.Managed, error)

type Entry struct {
	RenderFn       Fn
	ResourceFn     func() resource.Managed
	ResourceListFn func() resource.ManagedList
	MergeFn        MergeFn
	GetGvkKeyFn    GetGvkKeyFn
}

type Input struct {
	ParentInput *Input
	ObjectMeta  metav1.ObjectMeta // metadata for the returned CR
	//
	Meta   map[string]interface{}
	Object interface{} // target or node or a referenced object
	Data   interface{}
}

type Key struct {
	Name    string
	Version string
	//
	Vendor    targetv1.VendorType
	Platform  string
	SwVersion string
}

type catalog struct {
	m       *sync.RWMutex
	entries map[Key]Entry
}

func New(scs ...Catalog) Catalog {
	c := &catalog{
		m:       &sync.RWMutex{},
		entries: map[Key]Entry{},
	}
	c.Merge(scs...)
	return c
}

var Default = &catalog{
	m:       &sync.RWMutex{},
	entries: make(map[Key]Entry),
}

func (c *catalog) Get(key Key) (Entry, error) {
	c.m.RLock()
	defer c.m.RUnlock()
	if f, ok := c.entries[key]; ok {
		return f, nil
	}
	return Entry{}, fmt.Errorf("key: %v, not found", key)
}

func (c *catalog) Register(key Key, e Entry) {
	c.m.Lock()
	defer c.m.Unlock()
	c.entries[key] = e
}

func (c *catalog) Merge(scs ...Catalog) {
	for _, sc := range scs {
		for _, k := range sc.List() {
			e, err := sc.Get(k)
			if err != nil {
				continue
			}
			c.Register(k, e)
		}
	}
}

func (c *catalog) List() []Key {
	c.m.RLock()
	defer c.m.RLock()
	r := make([]Key, 0, len(c.entries))
	for k := range c.entries {
		r = append(r, k)
	}
	return r
}

func RegisterEntries(c Catalog, fns map[Key]Entry) {
	for k, v := range fns {
		c.Register(k, v)
	}
}

func (in *Input) GetTarget() (*targetv1.Target, error) {
	switch obj := in.Object.(type) {
	case *targetv1.Target:
		return obj, nil
	default:
		return nil, errors.New("unexpected obj type")
	}
}

func (in *Input) GetNode() (*topologyv1alpha1.Node, error) {
	switch obj := in.Object.(type) {
	case *topologyv1alpha1.Node:
		return obj, nil
	default:
		return nil, errors.New("unexpected obj type")
	}
}
