package all

import (
	"fmt"

	"github.com/yndd/catalog"
	"github.com/yndd/catalog/vendors/srl"
	"github.com/yndd/catalog/vendors/sros"
	"github.com/yndd/ndd-runtime/pkg/resource"
	statev1alpha1 "github.com/yndd/state/apis/state/v1alpha1"
	targetv1 "github.com/yndd/target/apis/target/v1"
	configsrlv1alpha1 "github.com/yndd/config-srl/apis/srl/v1alpha1"
)

func init() {
	catalog.RegisterEntries(catalog.Default, Entries)
}

var Entries = map[catalog.Key]catalog.Entry{
	{
		Name:    "configure_lldp",
		Version: "latest",
	}: {
		RenderFn:       nil,
		ResourceFn:     nil,
		ResourceListFn: nil,
		MergeFn:        nil,
		GetGvkKeyFn:       GetGvkKey,
	},
	{
		Name:    "state_lldp",
		Version: "latest",
	}: {
		RenderFn: StateLLDP,
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

func GetGvkKey(key catalog.Key, in *catalog.Input) (string, catalog.Key, error) {
	t, err := in.GetTarget()
	if err != nil {
		return "", catalog.Key{}, err
	}

	vendor := t.GetDiscoveryInfo().VendorType

	switch vendor {
	case targetv1.VendorTypeNokiaSRL:
		return configsrlv1alpha1.SrlConfigKindAPIVersion, catalog.Key{Name: key.Name, Version: "latest", Vendor: vendor}, nil
	case targetv1.VendorTypeNokiaSROS:
		return "", catalog.Key{Name: key.Name, Version: "latest", Vendor: vendor}, nil
	default:
		return "", catalog.Key{}, fmt.Errorf("unsupported vendorType: %s", key.Vendor)
	}
}

func StateLLDP(key catalog.Key, in *catalog.Input) (resource.Managed, error) {
	t, err := in.GetTarget()
	if err != nil {
		return nil, err
	}
	key.Vendor = t.GetDiscoveryInfo().VendorType

	switch t.GetDiscoveryInfo().VendorType {
	case targetv1.VendorTypeNokiaSRL:
		return srl.StateLLDP(key, in)
	case targetv1.VendorTypeNokiaSROS:
		return sros.StateLLDP(key, in)
	default:
		return nil, fmt.Errorf("unsupported vendorType: %s", key.Vendor)
	}
}
