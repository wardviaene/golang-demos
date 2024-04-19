package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type StaticPageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []StaticPage `json:"items"`
}

type StaticPage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec StaticPageSpec `json:"spec"`
}

type StaticPageSpec struct {
	Contents string `json:"contents"`
	Image    string `json:"image"`
	Replicas int    `json:"replicas"`
}
