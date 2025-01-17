/*
Copyright 2018 The Kubernetes Authors.

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

package error

import (
	"strings"
	"testing"
)

func TestResourceError_Error(t *testing.T) {
	filepath := "/path/to/kustomize.yaml"
	resourcefilepath := "/path/to/resource/deployment.yaml"
	errorMsg := "file not found"
	me := ResourceError{KustomizationPath: filepath, ResourceFilepath: resourcefilepath, ErrorMsg: errorMsg}
	if !strings.Contains(me.Error(), filepath) {
		t.Errorf("Incorrect ResourceError.Error() message \n")
		t.Errorf("Expected filepath %s, but unfound\n", filepath)
	}
	if !strings.Contains(me.Error(), resourcefilepath) {
		t.Errorf("Incorrect ResourceError.Error() message \n")
		t.Errorf("Expected resourcefilepath %s, but unfound\n", resourcefilepath)
	}
	if !strings.Contains(me.Error(), errorMsg) {
		t.Errorf("Incorrect ResourceError.Error() message \n")
		t.Errorf("Expected errorMsg %s, but unfound\n", errorMsg)
	}
}
