package persistentvolumeclaims

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts allows to create a persistent volume claims using given parameters.
type CreateOpts struct {
	// The version of the persistent API, valid value is 'v1'.
	ApiVersion string `json:"apiVersion" required:"true"`
	// Kind is a string value representing the REST resource this object represents.
	// Servers may infer this from the endpoint the c submits requests to.
	Kind string `json:"kind" required:"true"`
	// Standard object's metadata.
	Metadata Metadata `json:"metadata" required:"true"`
	// The desired characteristics of a volume.
	Spec Spec `json:"spec" required:"true"`
}

// Metadata is an object which will be build up standard object metadata.
type Metadata struct {
	// The name of the persistent volume claim.
	Name string `json:"name,omitempty"`
	// The namespace of the persistent volume claim.
	Namespace string `json:"namespace,omitempty"`
	// The labels of the persistent volume claim.
	Labels map[string]string `json:"labels,omitempty"`
	// The annotations of the persistent volume claim.
	Annotations map[string]string `json:"annotations,omitempty"`
}

// Spec defines the detailed description of the cluster object.
type Spec struct {
	// Access mode of the volume. Only the first value in all selected options is valid.
	// ReadWriteOnce: The volume can be mounted as read-write by a single node.
	// NOTE:
	// This function is supported only when the cluster version is v1.13.10 and the storage-driver version is 1.0.19.
	// ReadOnlyMany (default): The volume can be mounted as read-only by many nodes.
	// ReadWriteMany: The volume can be mounted as read-write by many nodes.
	AccessModes []string `json:"accessModes" required:"true"`
	// Storage class name of the PVC.
	StorageClassName string `json:"storageClassName,omitempty"`
	// Resources represents the minimum resources the volume should have.
	Resources ResourceRequest `json:"resources" required:"true"`
	// Name of the PV bound to the PVC.
	VolumeName string `json:"volumeName,omitempty"`
}

// ResourceRequest is an object struct that represents the detailed of the volume.
type ResourceRequest struct {
	// Minimum amount of compute resources required.
	// If requests is omitted for a container, it defaults to Limits if that is explicitly specified,
	// otherwise to an implementation-defined value.
	Requests CapacityReq `json:"requests" required:"true"`
}

// CapacityReq is an object struct that represents the volume capacity.
type CapacityReq struct {
	// Volume size, in GB format: 'xGi'.
	Storage string `json:"storage" required:"true"`
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToPvcCreateMap() (map[string]interface{}, error)
}

// ToPvcCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToPvcCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and uses the namespace name to import a volume into the namespace.
func Create(c *golangsdk.ServiceClient, clusterId, ns string, opts CreateOptsBuilder) (r CreateResult) {
	reqBody, err := opts.ToPvcCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c, clusterId, ns), reqBody, &r.Body, nil)
	return
}

// List is a method to obtain an array of the persistent volume claims for specifies namespace.
func List(c *golangsdk.ServiceClient, clusterId, ns string) pagination.Pager {
	return pagination.NewPager(c, listURL(c, clusterId, ns), func(r pagination.PageResult) pagination.Page {
		return PersistentVolumeClaimPage{pagination.SinglePageBase(r)}
	})
}

// Delete accepts to delete the specifies persistent volume claim form the namespace.
func Delete(c *golangsdk.ServiceClient, clusterId, ns, name string) (r DeleteResult) {
	_, r.Err = c.DeleteWithBodyResp(deleteURL(c, clusterId, ns, name), nil, r.Body, nil)
	return
}
