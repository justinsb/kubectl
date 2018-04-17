/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Kit is a specification for a Kit resource
type Kit struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KitSpec   `json:"spec"`
	Status KitStatus `json:"status"`
}

// KitSpec is the spec for a Kit resource
type KitSpec struct {
	Bases     []KitBase                   `json:"bases"`
	Patchsets []KitPatchset               `json:"patchsets"`
	Objects   []unstructured.Unstructured `json:"objects"`
}

// KitBase describes a base package, which can be seen as inheritance or inclusion
type KitBase struct {
	Source   string `json:"source"`
	Optional bool   `json:"optional,omitempty"`
}

// KitPatchset specifies a patchset to be included into a Kit
type KitPatchset struct {
	// Patch is an inline patch
	Patch *unstructured.Unstructured `json:"patch,omitempty"`
	// Source is a reference to an external patch
	Source string `json:"source,omitempty"`
}

// KitStatus is the status for a Kit resource
type KitStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KitList is a list of Kit resources
type KitList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Kit `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Patchset describes a list of patches that should be applied to a package
type Patchset struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PatchsetSpec   `json:"spec"`
	Status PatchsetStatus `json:"status"`
}

// PatchsetSpec is the spec for a Patchset resource
type PatchsetSpec struct {
	Patches []PatchSpec `json:"patches"`
}

// PatchSpec is an individual patch, part of a Patchset
type PatchSpec struct {
	Patch *unstructured.Unstructured `json:"patch,omitempty"`
}

// PatchsetStatus is the status for a Patchset resource
type PatchsetStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PatchsetList is a list of Patchset resources
type PatchsetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Patchset `json:"items"`
}
