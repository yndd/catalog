package all

import (
	"fmt"

	"github.com/yndd/catalog"
	"github.com/yndd/catalog/vendors/srl"
	"github.com/yndd/catalog/vendors/sros"
	"github.com/yndd/ndd-runtime/pkg/resource"
	targetv1 "github.com/yndd/target/apis/target/v1"
	statev1alpha1 "github.com/yndd/state/apis/state/v1alpha1"
)

func init() {
	catalog.RegisterEntries(catalog.Default, Entries)
}

var Entries = map[catalog.Key]catalog.Entry{
	{
		Name:    "configure_lldp",
		Version: "latest",
	}: {
		RenderRn:       ConfigureLLDP,
		ResourceFn:     nil,
		ResourceListFn: nil,
		MergeFn:        nil,
	},
	{
		Name:    "state_lldp",
		Version: "latest",
	}: {
		RenderRn:       StateLLDP,
		ResourceFn: func() resource.Managed {
			return &statev1alpha1.State{}
		},
		ResourceListFn: func() resource.ManagedList {
			return &statev1alpha1.StateList{}
		},
		MergeFn: func(crs ...resource.Managed) (resource.Managed, error) {
			return nil, nil
		},
	},
}

func ConfigureLLDP(key catalog.Key, in *catalog.Input) (resource.Managed, error) {
	t, err := in.GetTarget()
	if err != nil {
		return nil, err
	}
	key.Vendor = t.GetDiscoveryInfo().VendorType

	switch t.GetDiscoveryInfo().VendorType {
	case targetv1.VendorTypeNokiaSRL:
		return srl.ConfigureLLDP(key, in)
	case targetv1.VendorTypeNokiaSROS:
		return sros.ConfigureLLDP(key, in)
	default:
		return nil, fmt.Errorf("unsupported vendorType: %s", key.Vendor)
	}
}

func StateLLDP(key catalog.Key, in *catalog.Input) (resource.Managed, error) {
	t, err := in.GetTarget()
	if err != nil {
		return nil, err
	}
	switch t.GetDiscoveryInfo().VendorType {
	case targetv1.VendorTypeNokiaSRL:
		return srl.StateLLDP(key, in)
	case targetv1.VendorTypeNokiaSROS:
		return sros.StateLLDP(key, in)
	default:
		return nil, fmt.Errorf("unsupported vendorType: %s", key.Vendor)
	}
}
