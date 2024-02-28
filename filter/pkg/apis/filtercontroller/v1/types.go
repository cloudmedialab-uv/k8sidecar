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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MyResource is a specification for a MyResource resource
type Filter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec FilterSpec `json:"spec"`
}

// FilterSpec is the spec for a Filter resource
type FilterSpec struct {
	Sidecars []Sidecar       `json:"sidecars"`
	Volumes  []corev1.Volume `json:"volumes"`
}

// Object represents a single object in the array
type Sidecar struct {
	Image       string               `json:"image"`
	Name        string               `json:"name,omitempty"`
	Priority    int8                 `json:"priority,omitempty"`
	Env         []corev1.EnvVar      `json:"env,omitempty"`
	VolumeMount []corev1.VolumeMount `json:"volumeMount,omitempty"`
}
