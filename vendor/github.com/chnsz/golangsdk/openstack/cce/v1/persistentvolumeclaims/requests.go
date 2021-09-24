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
	// The name of the persistent volume claim, must be unique within a namespace. Cannot be updated.
	Name string `json:"name" required:"true"`
	// Map of string keys and values that can be used to organize and categorize (scope and select) objects.
	// May match selectors of replication controllers and services.
	Labels *Labels `json:"labels,omitempty"`
}

// Labels is an object that can be used to organize and categorize (scope and select) objects.
type Labels struct {
	// Role-based access control (RBAC).
	// If enabled, access to resources in the namespace will be controlled by RBAC policies.
	Region string `json:"failure-domain.beta.kubernetes.io/region,omitempty"`
	// ID of enterprise project.
	AvailabilityZone string `json:"failure-domain.beta.kubernetes.io/zone,omitempty"`
}

// Spec defines the detailed description of the cluster object.
type Spec struct {
	// ID of an existing storage volume.
	// If an SFS, SFS Turbo, or EVS volume is used, set this parameter to the ID of the volume.
	// If an OBS bucket is used, set this parameter to the OBS bucket name.
	VolumeID string `json:"volumeID" required:"true"`
	// Cloud storage class. This parameter is used together with volumeID. That is, volumeID and storageType must be configured at the same time.
	// bs: EVS. For details, see Using EVS Disks as Storage Volumes.
	// nfs: SFS. For details, see Using SFS File Systems as Storage Volumes.
	// obs: OBS. For details, see Using OBS Buckets as Storage Volumes].
	// efs: SFS Turbo. For details, see Using SFS Turbo File Systems as Storage Volumes
	StorageType string `json:"storageType" required:"true"`
	// Access mode of the volume. Only the first value in all selected options is valid.
	// ReadWriteOnce: The volume can be mounted as read-write by a single node.
	// NOTE:
	// This function is supported only when the cluster version is v1.13.10 and the storage-driver version is 1.0.19.
	// ReadOnlyMany (default): The volume can be mounted as read-only by many nodes.
	// ReadWriteMany: The volume can be mounted as read-write by many nodes.
	AccessModes []string `json:"accessModes" required:"true"`
	// Storage class name of the PVC.
	StorageClassName string `json:"storageClassName"`
	// Name of the PV bound to the PVC.
	VolumeName string `json:"volumeName"`
	// PV type specified by the PVC.
	VolumeMode string `json:"volumeMode"`
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
