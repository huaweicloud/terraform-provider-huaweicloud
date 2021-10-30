package namespaces

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts allows to create a namespace using given parameters.
type CreateOpts struct {
	// Kind is a string value representing the REST resource this object represents.
	// Servers may infer this from the endpoint the client submits requests to. Cannot be updated.
	// The value of this parameter is Namespace.
	Kind string `json:"kind" required:"true"`
	// ApiVersion defines the versioned schema of this representation of an object.
	// Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values.
	// The value of this parameter is v1.
	ApiVersion string `json:"apiVersion" required:"true"`
	// Standard object metadata.
	Metadata Metadata `json:"metadata" required:"true"`
	// Spec defines the behavior of the Namespace.
	Spec *Spec `json:"spec,omitempty"`
	// Status describes the current status of a Namespace.
	Status *Status `json:"status,omitempty"`
}

// Metadata is an object which will be build up standard object metadata.
type Metadata struct {
	// Namespace name. Is required when creating resources, although some resources may allow a client to request the
	// generation of an appropriate name automatically.
	// Name is primarily intended for creation idempotence and configuration definition. Cannot be updated.
	// The name is consists of 1 to 63 characters and must be a regular expression [a-z0-9]([-a-z0-9]*[a-z0-9])?.
	Name string `json:"name,omitempty"`
	// An initializer is a controller which enforces some system invariant at object creation time.
	Initializers *Initializers `json:"initializers,omitempty"`
	// Enabled identify whether the resource is available.
	Enabled bool `json:"enable,omitempty"`
	// An optional prefix used by the server to generate a unique name ONLY IF the Name field has not been provided.
	// The name is consists of 1 to 253 characters and must bee a regular expression [a-z0-9]([-a-z0-9]*[a-z0-9])?.
	GenerateName string `json:"generateName,omitempty"`
	// Map of string keys and values that can be used to organize and categorize (scope and select) objects.
	// May match selectors of replication controllers and services.
	Labels map[string]interface{} `json:"labels,omitempty"`
	// An unstructured key value map stored with a resource that may be set by external tools to store and retrieve
	// arbitrary metadata. They are not queryable and should be preserved when modifying objects.
	Annotations map[string]interface{} `json:"annotations,omitempty"`
	// List of objects depended by this object.
	OwnerReferences []OwnerReference `json:"ownerReferences,omitempty"`
	// Finalizers is an opaque list of values that must be empty to permanently remove object from storage.
	Finalizers []string `json:"finalizers,omitempty"`
}

// Initializers is a controller which enforces some system invariant at namespace creation time.
type Initializers struct {
	// List of initializers that must execute in order before this object is visible.
	Pendings []Pending `json:"pending,omitempty"`
}

// Pending is an object of initializers that must execute in order before this object is visible.
type Pending struct {
	// Name of the process that is responsible for initializing this object.
	Name string `json:"name,omitempty"`
}

// OwnerReference is a list of objects depended by this object.
type OwnerReference struct {
	// API version of the referent.
	ApiVersion string `json:"apiVersion" required:"true"`
	// Kind of the referent.
	Kind string `json:"kind" required:"true"`
	// Name of the referent.
	Name string `json:"name" required:"true"`
	// If true, AND if the owner has the "foregroundDeletion" finalizer, then the owner cannot be deleted from the
	// key-value store until this reference is removed. Defaults to false.
	BlockOwnerDeletion bool `json:"blockOwnerDeletion,omitempty"`
}

// Spec defines the behavior of the Namespace.
type Spec struct {
	// Finalizers is an opaque list of values that must be empty to permanently remove object from storage.
	Finalizers []string `json:"finalizers,omitempty"`
}

// Status describes the current status of a Namespace.
type Status struct {
	Phase string `json:"phase,omitempty"`
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToNamespaceCreateMap() (map[string]interface{}, error)
}

// ToNamespaceCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToNamespaceCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and uses the values to create a new namespace.
func Create(client *golangsdk.ServiceClient, clusterID string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToNamespaceCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client, clusterID), b, &r.Body, nil)
	return
}

// Get is a method to obtain the specified CCI namespace by name.
func Get(client *golangsdk.ServiceClient, clusterID, name string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, clusterID, name), &r.Body, nil)
	return
}

// ListOpts allows to filter list data using given parameters.
type ListOpts struct {
	Pretty               string `q:"pretty"`
	Continue             string `q:"continue"`
	LabelSelector        string `q:"labelSelector"`
	FieldSelector        string `q:"fieldSelector"`
	Limit                string `q:"limit"`
	IncludeUninitialized string `q:"includeUninitialized"`
	Watch                string `q:"watch"`
	ResourceVersion      string `q:"resourceVersion"`
	Timeout              string `q:"timeoutSeconds"`
}

// ToNamespaceListQuery is a method which to build a request query by the ListOpts.
func (opts ListOpts) ToNamespaceListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// ListOptsBuilder is an interface which to support request query build of the namespace search.
type ListOptsBuilder interface {
	ToNamespaceListQuery() (string, error)
}

// List is a method to obtain the specified namespaces according to the ListOpts.
func List(client *golangsdk.ServiceClient, clusterID string, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(client, clusterID)
	if opts != nil {
		query, err := opts.ToNamespaceListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return NamespacePage{pagination.SinglePageBase(r)}
	})
}

// Delete is a method to delete an existing namespace.
func Delete(client *golangsdk.ServiceClient, clusterID, name string) (r DeleteResult) {
	reqOpt := &golangsdk.RequestOpts{
		JSONBody: map[string]string{
			"kind":       "DeleteOptions",
			"apiVersion": "v1",
		},
	}
	_, r.Err = client.Delete(resourceURL(client, clusterID, name), reqOpt)
	return
}
