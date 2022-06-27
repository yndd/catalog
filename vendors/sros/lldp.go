package sros

import (
	"strings"

	"github.com/openconfig/ygot/ygot"
	"github.com/yndd/catalog"
	nddv1 "github.com/yndd/ndd-runtime/apis/common/v1"
	"github.com/yndd/ndd-runtime/pkg/resource"
	statev1alpha1 "github.com/yndd/state/apis/state/v1alpha1"
	"github.com/yndd/state/pkg/ygotnddpstate"
	targetv1 "github.com/yndd/target/apis/target/v1"
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
		Vendor:    targetv1.VendorTypeNokiaSROS,
		Platform:  "",
		SwVersion: "",
	}: {
		RenderFn: ConfigureLLDP,
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
	return nil, nil
}

func StateLLDP(key catalog.Key, in *catalog.Input) (resource.Managed, error) {
	t, err := in.GetTarget()
	if err != nil {
		return nil, err
	}
	paths := []string{
		// "system/lldp/chassis-id",
		// "system/lldp/chassis-id-type",
		"/system/lldp/interface[name=*]/neighbor[id=*]", // TODO update the field
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
