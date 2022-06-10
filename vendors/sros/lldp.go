package sros

import (
	"github.com/yndd/catalog"
	"github.com/yndd/ndd-runtime/pkg/resource"
	statev1alpha1 "github.com/yndd/state/apis/state/v1alpha1"
	targetv1 "github.com/yndd/target/apis/target/v1"
)

func init() {
	catalog.RegisterFns(catalog.Default, Fns)
}

var Fns = map[catalog.FnKey]catalog.Fn{
	{
		Name:      "configure_lldp",
		Version:   "latest",
		Vendor:    targetv1.VendorTypeNokiaSROS,
		Platform:  "",
		SwVersion: "",
	}: ConfigureLLDP,
	{
		Name:      "state_lldp",
		Version:   "latest",
		Vendor:    targetv1.VendorTypeNokiaSROS,
		Platform:  "",
		SwVersion: "",
	}: StateLLDP,
}

func ConfigureLLDP(in *catalog.Input) (resource.Managed, error) {
	return nil, nil
}

func StateLLDP(in *catalog.Input) (resource.Managed, error) {
	return &statev1alpha1.State{}, nil
}
