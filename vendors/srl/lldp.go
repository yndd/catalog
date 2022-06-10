package srl

import (
	"errors"

	"github.com/openconfig/ygot/ygot"
	"github.com/yndd/catalog"
	configsrlv1alpha1 "github.com/yndd/config-srl/apis/srl/v1alpha1"
	"github.com/yndd/ndd-runtime/pkg/resource"
	statev1alpha1 "github.com/yndd/state/apis/state/v1alpha1"
	"github.com/yndd/state/pkg/ygotnddpstate"
	targetv1 "github.com/yndd/target/apis/target/v1"
	"github.com/yndd/ygotsrl"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
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
		Name:      "enable_lldp",
		Version:   "latest",
		Vendor:    targetv1.VendorTypeNokiaSRL,
		Platform:  "",
		SwVersion: "",
	}: EnableLLDP,
	{
		Name:      "disable_lldp",
		Version:   "latest",
		Vendor:    targetv1.VendorTypeNokiaSRL,
		Platform:  "",
		SwVersion: "",
	}: DisableLLDP,
	{
		Name:      "state_lldp",
		Version:   "latest",
		Vendor:    targetv1.VendorTypeNokiaSRL,
		Platform:  "",
		SwVersion: "",
	}: StateLLDP,
}

func ConfigureLLDP(in *catalog.Input) (resource.Managed, error) {
	switch data := in.Data.(type) {
	case string: // TODO: use a more elaborate data input type
		switch data {
		case "enable":
			return EnableLLDP(in)
		case "disable":
			return DisableLLDP(in)
		default:
			return nil, errors.New("unexpected data value")
		}
	default:
		return nil, errors.New("unexpected data type")
	}
}

func EnableLLDP(in *catalog.Input) (resource.Managed, error) {
	d := &ygotsrl.Device{
		System: &ygotsrl.SrlNokiaSystem_System{
			Lldp: &ygotsrl.SrlNokiaSystem_System_Lldp{
				AdminState: ygotsrl.SrlNokiaCommon_AdminState_enable,
			},
		},
	}
	b, err := ygot.Marshal7951(d)
	if err != nil {
		return nil, err
	}
	return &configsrlv1alpha1.SrlConfig{
		TypeMeta:   in.TypeMeta,
		ObjectMeta: in.ObjectMeta,
		Spec: configsrlv1alpha1.ConfigSpec{
			Properties: runtime.RawExtension{Raw: b},
		},
	}, nil
}

func DisableLLDP(in *catalog.Input) (resource.Managed, error) {
	d := &ygotsrl.Device{
		System: &ygotsrl.SrlNokiaSystem_System{
			Lldp: &ygotsrl.SrlNokiaSystem_System_Lldp{
				AdminState: ygotsrl.SrlNokiaCommon_AdminState_disable,
			},
		},
	}
	b, err := ygot.Marshal7951(d)
	if err != nil {
		return nil, err
	}
	return &configsrlv1alpha1.SrlConfig{
		TypeMeta:   in.TypeMeta,
		ObjectMeta: in.ObjectMeta,
		Spec: configsrlv1alpha1.ConfigSpec{
			Properties: runtime.RawExtension{Raw: b},
		},
	}, nil
}

func StateLLDP(in *catalog.Input) (resource.Managed, error) {
	paths := []string{
		// "system/lldp/chassis-id",
		// "system/lldp/chassis-id-type",
		"system/lldp/interface[name=*]/neighbor[name=*]",
	}
	d := &ygotnddpstate.YnddState_StateEntry{
		Name: pointer.String("lldp_state"),
		Path: paths,
	}
	b, err := ygot.Marshal7951(d)
	if err != nil {
		return nil, err
	}
	return &statev1alpha1.State{
		TypeMeta:   in.TypeMeta,
		ObjectMeta: in.ObjectMeta,
		Spec: statev1alpha1.StateSpec{
			Properties: runtime.RawExtension{
				Raw: b,
			},
		},
	}, nil
}
