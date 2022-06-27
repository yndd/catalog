package srl

import (
	"errors"
	"strings"

	"github.com/openconfig/ygot/ygot"
	"github.com/yndd/catalog"
	configsrlv1alpha1 "github.com/yndd/config-srl/apis/srl/v1alpha1"
	nddv1 "github.com/yndd/ndd-runtime/apis/common/v1"
	"github.com/yndd/ndd-runtime/pkg/resource"
	statev1alpha1 "github.com/yndd/state/apis/state/v1alpha1"
	"github.com/yndd/state/pkg/ygotnddpstate"
	targetv1 "github.com/yndd/target/apis/target/v1"
	"github.com/yndd/ygotsrl"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
)

func init() {
	catalog.RegisterEntries(catalog.Default, Entries)
}

var Entries = map[catalog.Key]catalog.Entry{
	{
		Name:      "configure_lldp",
		Version:   "latest",
		Vendor:    targetv1.VendorTypeNokiaSRL,
		Platform:  "",
		SwVersion: "",
	}: {
		RenderFn: ConfigureLLDP,
		ResourceFn: func() resource.Managed {
			return &configsrlv1alpha1.SrlConfig{}
		},
		ResourceListFn: func() resource.ManagedList {
			return &configsrlv1alpha1.SrlConfigList{}
		},
		MergeFn: configSRLMergeFn,
	},
	{
		Name:      "enable_lldp",
		Version:   "latest",
		Vendor:    targetv1.VendorTypeNokiaSRL,
		Platform:  "",
		SwVersion: "",
	}: {
		RenderFn: EnableLLDP,
		ResourceFn: func() resource.Managed {
			return &configsrlv1alpha1.SrlConfig{}
		},
		ResourceListFn: func() resource.ManagedList {
			return &configsrlv1alpha1.SrlConfigList{}
		},
		MergeFn: configSRLMergeFn,
	},
	{
		Name:      "disable_lldp",
		Version:   "latest",
		Vendor:    targetv1.VendorTypeNokiaSRL,
		Platform:  "",
		SwVersion: "",
	}: {
		RenderFn: DisableLLDP,
		ResourceFn: func() resource.Managed {
			return &configsrlv1alpha1.SrlConfig{}
		},
		ResourceListFn: func() resource.ManagedList {
			return &configsrlv1alpha1.SrlConfigList{}
		},
		MergeFn: configSRLMergeFn,
	},
	{
		Name:      "state_lldp",
		Version:   "latest",
		Vendor:    targetv1.VendorTypeNokiaSRL,
		Platform:  "",
		SwVersion: "",
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

func ConfigureLLDP(key catalog.Key, in *catalog.Input) (resource.Managed, error) {
	switch data := in.Data.(type) {
	case string: // TODO: use a more elaborate data input type
		switch data {
		case "enable":
			return EnableLLDP(key, in)
		case "disable":
			return DisableLLDP(key, in)
		default:
			return nil, errors.New("unexpected data value")
		}
	default:
		return nil, errors.New("unexpected data type")
	}
}

func EnableLLDP(key catalog.Key, in *catalog.Input) (resource.Managed, error) {
	t, err := in.GetTarget()
	if err != nil {
		return nil, err
	}
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
	in.ObjectMeta.Name = strings.Join([]string{in.ObjectMeta.Name, t.GetName()}, ".")
	return &configsrlv1alpha1.SrlConfig{
		ObjectMeta: in.ObjectMeta,
		Spec: configsrlv1alpha1.ConfigSpec{
			ResourceSpec: nddv1.ResourceSpec{
				TargetReference: &nddv1.Reference{
					Name: t.GetName(),
				},
			},
			Properties: runtime.RawExtension{Raw: b},
		},
	}, nil
}

func DisableLLDP(key catalog.Key, in *catalog.Input) (resource.Managed, error) {
	t, err := in.GetTarget()
	if err != nil {
		return nil, err
	}
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
	in.ObjectMeta.Name = strings.Join([]string{in.ObjectMeta.Name, t.GetName()}, ".")
	return &configsrlv1alpha1.SrlConfig{
		ObjectMeta: in.ObjectMeta,
		Spec: configsrlv1alpha1.ConfigSpec{
			ResourceSpec: nddv1.ResourceSpec{
				TargetReference: &nddv1.Reference{
					Name: t.GetName(),
				},
			},
			Properties: runtime.RawExtension{Raw: b},
		},
	}, nil
}

func StateLLDP(key catalog.Key, in *catalog.Input) (resource.Managed, error) {
	t, err := in.GetTarget()
	if err != nil {
		return nil, err
	}
	paths := []string{
		// "system/lldp/chassis-id",
		// "system/lldp/chassis-id-type",
		"/system/lldp/interface[name=*]/neighbor[id=*]",
	}
	d := &ygotnddpstate.YnddState_StateEntry{
		Name: pointer.String("lldp_state"),
		Path: paths,
	}
	b, err := ygot.Marshal7951(d)
	if err != nil {
		return nil, err
	}

	if in.ObjectMeta.Annotations == nil {
		in.ObjectMeta.Annotations = make(map[string]string)
	}
	in.ObjectMeta.Annotations["state.yndd.io/paths"] = strings.Join(paths, ",")
	in.ObjectMeta.Name = strings.Join([]string{in.ObjectMeta.Name, t.GetName()}, ".")

	return &statev1alpha1.State{
		ObjectMeta: in.ObjectMeta,
		Spec: statev1alpha1.StateSpec{
			ResourceSpec: nddv1.ResourceSpec{
				Lifecycle: nddv1.Lifecycle{},
				TargetReference: &nddv1.Reference{
					Name: t.GetName(),
				},
			},
			Properties: runtime.RawExtension{Raw: b},
		},
	}, nil
}
