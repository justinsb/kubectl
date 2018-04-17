package builder

import (
	"fmt"

	"k8s.io/kubectl/pkg/loader"

	"k8s.io/kubectl/pkg/kustomize/transformers"

	"k8s.io/kubectl/pkg/kustomize/apis/packaging/v1alpha1"
	"k8s.io/kubectl/pkg/kustomize/resource"
)

type Builder struct {
	resources []*resource.Resource
	patches   []*resource.Resource
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) loadResources(kit *v1alpha1.Kit, l loader.Loader) error {
	for i := range kit.Spec.Bases {
		contents, err := l.Load(kit.Spec.Bases[i].Source + ".yaml")
		if err != nil {
			return err
		}
		obj, gvk, err := ParseObject(contents)
		if err != nil {
			return fmt.Errorf("error parsing %s: %v", kit.Spec.Patchsets[i].Source, err)
		}
		if gvk == nil {
			return fmt.Errorf("no GroupVersionKind returned from parse")
		}
		if gvk.Kind != "Kit" {
			return fmt.Errorf("unexpected object %s, expected Kit", gvk.Kind)
		}
		baseKit, ok := obj.(*v1alpha1.Kit)
		if !ok {
			return fmt.Errorf("unexpected type for Kit: %T", obj)
		}
		if err := b.loadResources(baseKit, l); err != nil {
			return err
		}
	}

	for i := range kit.Spec.Objects {
		b.resources = append(b.resources, &resource.Resource{Data: &kit.Spec.Objects[i]})
	}

	for i := range kit.Spec.Patchsets {
		if kit.Spec.Patchsets[i].Patch != nil {
			b.patches = append(b.patches, &resource.Resource{Data: kit.Spec.Patchsets[i].Patch})
		}

		if kit.Spec.Patchsets[i].Source != "" {
			contents, err := l.Load(kit.Spec.Patchsets[i].Source + ".yaml")
			if err != nil {
				return err
			}
			obj, gvk, err := ParseObject(contents)
			if err != nil {
				return fmt.Errorf("error parsing %s: %v", kit.Spec.Patchsets[i].Source, err)
			}
			if gvk == nil {
				return fmt.Errorf("no GroupVersionKind returned from parse")
			}
			if gvk.Kind != "Patchset" {
				return fmt.Errorf("unexpected object %s, expected Patchset", gvk.Kind)
			}
			patchset, ok := obj.(*v1alpha1.Patchset)
			if !ok {
				return fmt.Errorf("unexpected type for Patchset: %T", obj)
			}
			for i := range patchset.Spec.Patches {
				b.patches = append(b.patches, &resource.Resource{Data: patchset.Spec.Patches[i].Patch})
			}
		}
	}

	return nil
}

func (b *Builder) ExpandKit(kit *v1alpha1.Kit, l loader.Loader) (resource.ResourceCollection, error) {
	if err := b.loadResources(kit, l); err != nil {
		return nil, err
	}

	rc := resource.ResourceCollection{}
	for _, res := range b.resources {
		gvkn := res.GVKN()
		if _, found := rc[gvkn]; found {
			return nil, fmt.Errorf("GroupVersionKindName: %#v already exists in the map", gvkn)
		}
		rc[gvkn] = res
	}

	ts := []transformers.Transformer{}
	ot, err := transformers.NewOverlayTransformer(b.patches)
	if err != nil {
		return nil, err
	}
	ts = append(ts, ot)

	// npt, err := transformers.NewDefaultingNamePrefixTransformer(string(a.kustomization.NamePrefix))
	// if err != nil {
	// 	return nil, err
	// }
	// ts = append(ts, npt)

	// lt, err := transformers.NewDefaultingLabelsMapTransformer(a.kustomization.LabelsToAdd)
	// if err != nil {
	// 	return nil, err
	// }
	// ts = append(ts, lt)

	// at, err := transformers.NewDefaultingAnnotationsMapTransformer(a.kustomization.AnnotationsToAdd)
	// if err != nil {
	// 	return nil, err
	// }
	// ts = append(ts, at)

	t := transformers.NewMultiTransformer(ts)

	err = t.Transform(rc)
	if err != nil {
		return nil, err
	}

	return rc, nil
}
