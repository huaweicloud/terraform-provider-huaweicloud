package persistentvolumeclaims

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type CreateOpts struct {
	// The version of the persistent API, valid value is 'v1'.
	ApiVersion string `json:"apiVersion" required:"true"`
	// Kind is a string value representing the REST resource this object represents.
	// Servers may infer this from the endpoint the client submits requests to.
	Kind string `json:"kind" required:"true"`
	// Standard object's metadata.
	Metadata Metadata `json:"metadata" required:"true"`
	// The desired characteristics of a volume.
	Spec Spec `json:"spec" required:"true"`
}

type Metadata struct {
	// The name of the persistent volume claim, must be unique within a namespace. Cannot be updated.
	Name string `json:"name" required:"true"`
	// Namespace defines the space within each name must be unique.
	// An empty namespace is equivalent to the 'default' namespace, but 'default' is the canonical representation.
	Namespace string `json:"namespace,omitempty"`
	// An unstructured key value map stored with a resource that may be set by external tools to store and retrieve
	// arbitrary metadata.
	Annotations *Annotations `json:"annotations,omitempty"`
}

type Annotations struct {
	// The type of the file system.
	// The valid values are ext4 (EVS disk), obs (OBS bucket) and nfs (SFS or SFS Turbo).
	FsType string `json:"fsType" required:"true"`
	// ID of the volume
	VolumeID string `json:"volumeID" required:"true"`
	// The Shared path of the SFS and the SFS Turbo.
	DeviceMountPath string `json:"deviceMountPath,omitempty"`
}

type Spec struct {
	// Resources represents the minimum resources the volume should have.
	Resources ResourceRequirement `json:"resources" required:"true"`
	// Name of the storage class required by the claim.
	// The following fields are supported:
	//     EVS: sas, ssd and sata
	//     SFS: nfs-rw
	//     SFS Turbo: efs-performance and efs-standard
	//     OBS: obs
	StorageClassName string `json:"storageClassName" required:"true"`
	// AccessModes contains the actual access modes the volume backing the PVC has.
	//     ReadWriteOnce: can be mount read/write mode to exactly 1 host.
	//     ReadOnlyMany: can be mount in read-only mode to many hosts.
	//     ReadWriteMany: can be mount in read/write mode to many hosts.
	AccessModes []string `json:"accessModes,omitempty"`
}

type ResourceRequirement struct {
	// Minimum amount of compute resources required.
	// If requests is omitted for a container, it defaults to Limits if that is explicitly specified,
	// otherwise to an implementation-defined value.
	Requests *ResourceName `json:"requests,omitempty"`
}

type ResourceName struct {
	// Volume size, in GB format: 'xGi'.
	Storage string `json:"storage,omitempty"`
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToPVCCreateMap() (map[string]interface{}, error)
}

// ToPVCCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToPVCCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and uses the namespace name to import a volume into the namespace.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder, ns string) (r CreateResult) {
	reqBody, err := opts.ToPVCCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client, ns), reqBody, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

type ListOpts struct {
	// Type of the storage, valid values are bs, obs, nfs and efs.
	StorageType string `q:"storage_type"`
}

type ListOptsBuilder interface {
	ToPVCListQuery() (string, error)
}

func (opts ListOpts) ToPVCListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// List is a method to obtain an array of the persistent volume claims for specifies namespace.
func List(client *golangsdk.ServiceClient, opts ListOptsBuilder, ns string) pagination.Pager {
	url := rootURL(client, ns)
	if opts != nil {
		query, err := opts.ToPVCListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return PersistentVolumeClaimPage{pagination.SinglePageBase(r)}
	})
}

// Delete accepts to delete the specifies persistent volume claim form the namespace.
func Delete(client *golangsdk.ServiceClient, ns, name string) (r DeleteResult) {
	_, r.Err = client.DeleteWithBodyResp(resourceURL(client, ns, name), nil, r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
