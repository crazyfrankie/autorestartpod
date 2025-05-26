/*
Copyright 2025.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AutoRestartPodSpec defines the desired state of AutoRestartPod.
type AutoRestartPodSpec struct {
	Schedule string               `json:"schedule"`           // 定义Cron表达式 (例如 "0 3 * * *" 或 "30 */5 * * * *")
	Selector metav1.LabelSelector `json:"selector"`           // 定义用于选择要重启的Pod的标签选择器
	TimeZone string               `json:"timeZone,omitempty"` // 可选：时区 (例如 "Asia/Shanghai")
}

// AutoRestartPodStatus defines the observed state of AutoRestartPod.
type AutoRestartPodStatus struct {
	LastRestartTime *metav1.Time `json:"lastRestartTime,omitempty"` // Record the last reboot time
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// AutoRestartPod is the Schema for the autorestartpods API.
type AutoRestartPod struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AutoRestartPodSpec   `json:"spec,omitempty"`
	Status AutoRestartPodStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AutoRestartPodList contains a list of AutoRestartPod.
type AutoRestartPodList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AutoRestartPod `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AutoRestartPod{}, &AutoRestartPodList{})
}
