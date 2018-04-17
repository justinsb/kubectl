package builder

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"k8s.io/kubectl/pkg/kustomize/util/fs"
	"k8s.io/kubectl/pkg/loader"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/kubectl/pkg/kustomize/apis/packaging/v1alpha1"
	"k8s.io/kubectl/pkg/kustomize/controller/builder/testhelpers"
)

func TestReadKit(t *testing.T) {
	ts := testhelpers.NewSuite(t)
	data := ts.MustRead("testdata/coredns.io/coredns/1.1.1/kit.yaml")
	obj, gvk := ts.MustParse(data)
	ts.MustDeepEqual("Kit", gvk.Kind)
	t.Logf("Kit: %v", ts.MustSerialize(obj))
	kit := obj.(*v1alpha1.Kit)
	ts.MustDeepEqual("coredns-base", kit.GetName())
	ts.MustDeepEqual(6, len(kit.Spec.Objects))
}

func TestReadPatchset(t *testing.T) {
	ts := testhelpers.NewSuite(t)
	data := ts.MustRead("testdata/coredns.io/coredns/1.1.1/system/coredns-kops.yaml")
	obj, gvk := ts.MustParse(data)
	ts.MustDeepEqual("Patchset", gvk.Kind)
	t.Logf("Patchset: %v", ts.MustSerialize(obj))
	ps := obj.(*v1alpha1.Patchset)
	ts.MustDeepEqual(1, len(ps.Spec.Patches))
}

func TestBuild_NoPatches(t *testing.T) {
	runBuildTest(t, "nopatches")
}

func TestBuild_InlinePatch(t *testing.T) {
	runBuildTest(t, "inlinepatch")
}

func TestBuild_ExternalPatch(t *testing.T) {
	runBuildTest(t, "externalpatch")
}

func TestBuild_Base(t *testing.T) {
	runBuildTest(t, "base")
}

func runBuildTest(t *testing.T, tn string) {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("error getting working directory: %v", err)
	}

	basedir := filepath.Join(currentDir, "testdata", "dummy", tn)

	ts := testhelpers.NewSuite(t)
	data := ts.MustRead(filepath.Join(basedir, "kit.yaml"))
	obj, gvk := ts.MustParse(data)
	ts.MustDeepEqual("Kit", gvk.Kind)
	t.Logf("Kit: %v", ts.MustSerialize(obj))
	kit := obj.(*v1alpha1.Kit)

	b := NewBuilder()

	vfsBase := fs.MakeRealFS()
	vfs := loader.NewFileLoader(vfsBase)

	l := loader.Init([]loader.SchemeLoader{vfs})
	subdir, err := l.New(basedir)
	if err != nil {
		t.Fatalf("error building vfs: %v", err)
	}

	expanded, err := b.ExpandKit(kit, subdir)
	if err != nil {
		t.Fatalf("error expanding Kit: %v", err)
	}
	t.Logf("expanded: %v", expanded)

	buffer := &bytes.Buffer{}

	for _, r := range expanded {
		o := r.Data

		a := ts.MustSerializeCore(o)
		buffer.WriteString(a)
	}

	expectedYaml := ts.MustRead("testdata/dummy/" + tn + "/_expected.yaml")
	expectedYaml = strings.TrimSpace(expectedYaml)

	actualYaml := strings.TrimSpace(buffer.String())

	ts.MustDeepEqual(expectedYaml, actualYaml)
}

func parseExpected(ts *testhelpers.Suite, p string) map[string]*unstructured.Unstructured {
	data := ts.MustRead(p)

	decoder := k8syaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(data)), 1024)

	resources := make(map[string]*unstructured.Unstructured)
	for {
		o := &unstructured.Unstructured{}

		err := decoder.Decode(o)
		if err != nil {
			if err == io.EOF {
				break
			}

			ts.T.Fatalf("error parsing %s: %v", p, err)
		}

		k := o.GetKind() + ":" + o.GetNamespace() + "/" + o.GetName()

		resources[k] = o
	}
	return resources
}

// type testLoader struct {
// 	data map[string][]byte
// }

// var _ loader.SchemeLoader = &testLoader{}

// // Does this location correspond to this scheme.
// func (l *testLoader) IsScheme(root string, location string) bool {
// 	return true
// }

// // Combines the root and path into a full location string.
// func (l *testLoader) FullLocation(root string, path string) (string, error) {
// 	if root == "" {
// 		return path, nil
// 	}
// 	return "", fmt.Errorf("testLoader::FullLocation(%s, %s) not implemented", root, path)
// }

// // Load bytes at scheme-specific location or an error.
// func (l *testLoader) Load(location string) ([]byte, error) {
// 	b := l.data[location]
// 	if b != nil {
// 		return b, nil
// 	}
// 	return nil, fmt.Errorf("testLoader did not have a file for %q", location)
// }

// func (l *testLoader) Put(path string, data []byte) {
// 	if l.data == nil {
// 		l.data = make(map[string][]byte)
// 	}
// 	l.data[path] = data
// }
