package builder

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	packagingv1alpha1 "k8s.io/kubectl/pkg/kustomize/apis/packaging/v1alpha1"
)

var PackagingScheme = runtime.NewScheme()
var PackagingCodecs = serializer.NewCodecFactory(PackagingScheme)

func init() {
	packagingv1alpha1.AddToScheme(PackagingScheme)
}

func ParseObject(data []byte) (runtime.Object, *schema.GroupVersionKind, error) {
	return PackagingCodecs.UniversalDecoder(packagingv1alpha1.SchemeGroupVersion).Decode(data, nil, nil)
}
