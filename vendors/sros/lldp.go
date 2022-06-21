package sros

import (
	"github.com/yndd/catalog"
	"github.com/yndd/ndd-runtime/pkg/resource"
	statev1alpha1 "github.com/yndd/state/apis/state/v1alpha1"
	targetv1 "github.com/yndd/target/apis/target/v1"
)

func init() {
	catalog.RegisterEntries(catalog.Default, Entries)
}

var Entries = map[catalog.Key]catalog.Entry{
	{
		Name:      "configure_lldp",
		Version:   "latest",
		Vendor:    targetv1.VendorTypeNokiaSROS,
		Platform:  "",
		SwVersion: "",
	}: {
		RenderRn: ConfigureLLDP,
		ResourceFn: func() resource.Managed {
			return nil
		},
		ResourceListFn: func() resource.ManagedList {
			return nil
		},
		MergeFn: func(crs ...resource.Managed) (resource.Managed, error) {
			return nil, nil
		},
	},
	{
		Name:      "state_lldp",
		Version:   "latest",
		Vendor:    targetv1.VendorTypeNokiaSROS,
		Platform:  "",
		SwVersion: "",
	}: {
		RenderRn: StateLLDP,
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
	return nil, nil
}

func StateLLDP(key catalog.Key, in *catalog.Input) (resource.Managed, error) {
	return &statev1alpha1.State{}, nil
}
