package testhelpers

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	packagingv1alpha1 "k8s.io/kubectl/pkg/kustomize/apis/packaging/v1alpha1"
)

var PackagingScheme = runtime.NewScheme()
var PackagingCodecs = serializer.NewCodecFactory(PackagingScheme)
var CoreScheme = runtime.NewScheme()
var CoreCodecs = serializer.NewCodecFactory(CoreScheme)

func init() {
	packagingv1alpha1.AddToScheme(PackagingScheme)
	v1.AddToScheme(CoreScheme)
}

type Suite struct {
	T *testing.T
}

func NewSuite(t *testing.T) *Suite {
	return &Suite{T: t}
}

func (s *Suite) MustRead(p string) string {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		s.T.Fatalf("error reading file %q: %v", p, err)
	}
	return string(b)
}

func (s *Suite) MustParse(data string) (runtime.Object, *schema.GroupVersionKind) {
	obj, gvk, err := PackagingCodecs.UniversalDecoder(packagingv1alpha1.SchemeGroupVersion).Decode([]byte(data), nil, nil)
	if err != nil {
		s.T.Fatalf("error parsing data %q: %v", data, err)
	}
	return obj, gvk
}

func (s *Suite) MustParseCore(data string) (runtime.Object, *schema.GroupVersionKind) {
	obj, gvk, err := CoreCodecs.UniversalDecoder(v1.SchemeGroupVersion).Decode([]byte(data), nil, nil)
	if err != nil {
		s.T.Fatalf("error parsing data %q: %v", data, err)
	}
	return obj, gvk
}

func (s *Suite) MustSerialize(obj runtime.Object) string {
	mediaType := "application/yaml"
	gv := packagingv1alpha1.SchemeGroupVersion

	e, ok := runtime.SerializerInfoForMediaType(PackagingCodecs.SupportedMediaTypes(), mediaType)
	if !ok {
		s.T.Fatalf("no %s serializer registered", mediaType)
	}
	encoder := PackagingCodecs.EncoderForVersion(e.Serializer, gv)

	var w bytes.Buffer
	err := encoder.Encode(obj, &w)
	if err != nil {
		s.T.Fatalf("error encoding %T: %v", obj, err)
	}
	return w.String()
}

func (s *Suite) MustSerializeCore(obj runtime.Object) string {
	mediaType := "application/yaml"
	gv := v1.SchemeGroupVersion

	e, ok := runtime.SerializerInfoForMediaType(CoreCodecs.SupportedMediaTypes(), mediaType)
	if !ok {
		s.T.Fatalf("no %s serializer registered", mediaType)
	}
	encoder := CoreCodecs.EncoderForVersion(e.Serializer, gv)

	var w bytes.Buffer
	err := encoder.Encode(obj, &w)
	if err != nil {
		s.T.Fatalf("error encoding %T: %v", obj, err)
	}
	return w.String()
}

func (s *Suite) MustDeepEqual(e, a interface{}) {
	if !reflect.DeepEqual(e, a) {
		s.T.Fatalf("unexpected value.  expected=%v, actual=%v", e, a)
	}
}
