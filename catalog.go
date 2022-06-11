package catalog

import (
	"errors"
	"sync"

	"github.com/yndd/ndd-runtime/pkg/resource"
	targetv1 "github.com/yndd/target/apis/target/v1"
	topologyv1alpha1 "github.com/yndd/topology/apis/topo/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Catalog interface {
	GetFn(fnKey FnKey) (Fn, error)
	RegisterFn(fnKey FnKey, fn Fn)
	List() []FnKey
}

type Fn func(in *Input) (resource.Managed, error)

type Input struct {
	ObjectMeta metav1.ObjectMeta
	//
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

func (c *catalog) List() []FnKey {
	c.m.RLock()
	defer c.m.RLock()
	r := make([]FnKey, 0, len(c.fns))
	for k := range c.fns {
		r = append(r, k)
	}
	return r
}

func RegisterFns(c Catalog, fns map[FnKey]Fn) {
	for k, v := range fns {
		c.RegisterFn(k, v)
	}
}

func GetFn(c Catalog, name, version string, in *Input) (resource.Managed, error) {
	key := FnKey{
		Name:    name,
		Version: version,
	}
	switch obj := in.Object.(type) {
	case targetv1.Target:
		key.Vendor = obj.Spec.Properties.VendorType
		key.Platform = obj.Spec.DiscoveryInfo.Platform
		key.SwVersion = obj.Spec.DiscoveryInfo.SwVersion
	case topologyv1alpha1.Node:
		key.Vendor = obj.Spec.Properties.VendorType
		key.Platform = obj.Spec.Properties.Platform
		key.SwVersion = obj.Spec.Properties.ExpectedSWVersion
	default:
		return nil, errors.New("unexpected obj type")
	}
	fn, err := c.GetFn(key)
	if err != nil {
		return nil, err
	}
	return fn(in)
}
