package srl

import (
	"fmt"

	"github.com/openconfig/ygot/ygot"
	configsrlv1alpha1 "github.com/yndd/config-srl/apis/srl/v1alpha1"
	"github.com/yndd/ndd-runtime/pkg/resource"
	"github.com/yndd/ygotsrl"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func configSRLMergeFn(crs ...resource.Managed) (resource.Managed, error) {
	if len(crs) == 0 {
		return nil, nil
	}
	dc := &ygotsrl.Device{}
	var objMeta *metav1.ObjectMeta
	var err error

	for _, cr := range crs {
		srlConfig, ok := cr.(*configsrlv1alpha1.SrlConfig)
		if !ok {
			return nil, fmt.Errorf("unexpected resource type")
		}
		err = ygotsrl.Unmarshal(srlConfig.Spec.Properties.Raw, dc)
		if err != nil {
			return nil, err
		}
		if objMeta == nil {
			objMeta = &srlConfig.ObjectMeta
		}
	}

	b, err := ygot.Marshal7951(dc)
	if err != nil {
		return nil, err
	}
	return &configsrlv1alpha1.SrlConfig{
		ObjectMeta: *objMeta,
		Spec: configsrlv1alpha1.ConfigSpec{
			Properties: runtime.RawExtension{
				Raw: b,
			},
		},
	}, nil
}
