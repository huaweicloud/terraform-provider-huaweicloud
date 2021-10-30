package namespaces

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type commonResult struct {
	golangsdk.Result
}

// Namespace is a struct that represents the result of Create and Get methods.
type Namespace struct {
	// Standard object metadata.
	Metadata MetaResp `json:"metadata"`
	// Spec defines the behavior of the Namespace.
	Spec Spec `json:"spec"`
	// Status describes the current status of a Namespace.
	Status Status `json:"status"`
}

// MetaResp is a standard object metadata.
type MetaResp struct {
	// Namespace name.
	Name string `json:"name"`
	// A prefix used by the server to generate a unique name ONLY IF the Name field has not been provided.
	GenerateName string `json:"generateName,omitempty"`
	// URL representing this object.
	SelfLink string `json:"selfLink"`
	UID      string `json:"uid"`
	// String that identifies the server's internal version of this object that can be used by clients to determine when
	// objects have changed.
	ResourceVersion string `json:"resourceVersion"`
	// An unstructured key value map stored with a resource that may be set by external tools to store and retrieve
	// arbitrary metadata.
	Annotations map[string]interface{} `json:"annotations"`
	// Enabled identify whether the resource is available.
	Enable bool `json:"enable"`
	// Map of string keys and values that can be used to organize and categorize (scope and select) objects.
	Labels map[string]interface{} `json:"labels"`
	// An optional prefix used by the server to generate a unique name ONLY IF the Name field has not been provided.
	Genetation        string `json:"generation"`
	CreationTimestamp string `json:"creationTimestamp"`
	DeletionTimestamp string `json:"DeletionTimestamp"`
}

// CreateResult represents a result of the Create method.
type CreateResult struct {
	commonResult
}

// GetResult represents a result of the Get method.
type GetResult struct {
	commonResult
}

// DeleteResult represents a result of the Delete method.
type DeleteResult struct {
	commonResult
}

func (r commonResult) Extract() (*Namespace, error) {
	var s Namespace
	err := r.ExtractInto(&s)
	return &s, err
}

// NamespacePage represents a result of the List method.
type NamespacePage struct {
	pagination.SinglePageBase
}

// ExtractNamespaces is a method which to interpret the namespace pages as a namespace object list.
func ExtractNamespaces(r pagination.Page) ([]Namespace, error) {
	var s []Namespace
	err := r.(NamespacePage).Result.ExtractIntoSlicePtr(&s, "items")
	return s, err
}
