package srl

import (
	"github.com/yndd/catalog"
	configsrlv1alpha "github.com/yndd/config-srl/apis/srl/v1alpha1"
	"github.com/yndd/ndd-runtime/pkg/resource"
	statev1alpha "github.com/yndd/state/apis/state/v1alpha1"
	targetv1 "github.com/yndd/target/apis/target/v1"
)

func init() {
	catalog.RegisterFns(catalog.Default, Fns)
}

var Fns = map[catalog.FnKey]catalog.Fn{
	{
		Name:      "configure_lldp",
		Version:   "latest",
		Vendor:    targetv1.VendorTypeNokiaSRL,
		Platform:  "",
		SwVersion: "",
	}: ConfigureLLDP,
	{
		Name:      "state_lldp",
		Version:   "latest",
		Vendor:    targetv1.VendorTypeNokiaSRL,
		Platform:  "",
		SwVersion: "",
	}: StateLLDP,
}

func ConfigureLLDP(in *catalog.Input) (resource.Managed, error) {
	return &configsrlv1alpha.SrlConfig{}, nil
}

func StateLLDP(in *catalog.Input) (resource.Managed, error) {
	return &statev1alpha.State{}, nil
}
