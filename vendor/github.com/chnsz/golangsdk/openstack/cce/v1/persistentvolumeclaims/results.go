package persistentvolumeclaims

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// PersistentVolumeClaim is a struct that represents the result of Create, Get and List methods.
type PersistentVolumeClaim struct {
	// The version of the persistent API.
	ApiVersion string `json:"apiVersion"`
	// The REST resource of object represents.
	Kind string `json:"kind"`
	// Standard object's metadata.
	Metadata MetaResp `json:"metadata"`
	// The desired characteristics of a volume.
	Spec SpecResp `json:"spec"`
	// PVC status.
	Status Status `json:"status"`
}

// MetaResp is an object struct that represents the persistent volume claim metadata.
type MetaResp struct {
	// The name of the Persistent Volume Claim.
	Name string `json:"name"`
	// The namespace where the Persistent Volume Claim is located.
	Namespace string `json:"namespace"`
	// An unstructured key value map stored with a resource that may be set by external tools to store and retrieve
	// arbitrary metadata.
	Annotations map[string]string `json:"annotations"`
	// ID of the Persistent Volume Claim in UUID format.
	UID string `json:"uid"`
	// String that identifies the server's internal version of this object that can be used by clients to determine
	// when objects have changed.
	ResourceVersion string `json:"resourceVersion"`
	// A timestamp representing the server time when this object was created.
	CreationTimestamp string `json:"creationTimestamp"`
	// SelfLink is a URL representing this object.
	SelfLink string `json:"selfLink"`
	// Map of string keys and values that can be used to organize and categorize (scope and select) objects.
	Labels map[string]string `json:"labels"`
	// Each finalizer string of array is an identifier for the responsible component that will remove the entry form
	// the list.
	Finalizers []string `json:"finalizers"`
	// Enable identify whether the resource is available.
	Enable bool `json:"enable"`
}

// SpecResp is an object struct that represents the detailed description.
type SpecResp struct {
	// The name of the volume.
	VolumeName string `json:"volumeName"`
	// AccessModes contains the actual access modes the volume backing the PVC has.
	AccessModes []string `json:"accessModes"`
	// Resources represents the minimum resources the volume should have.
	Resources ResourceRequirement `json:"resources"`
	// Name of the storage class required by the claim.
	StorageClassName string `json:"storageClassName"`
	// Mode of the volume.
	VolumeMode string `json:"volumeMode"`
}

// ResourceRequirement is an object struct that represents the detailed of the volume.
type ResourceRequirement struct {
	// Minimum amount of compute resources required.
	// If requests is omitted for a container, it defaults to Limits if that is explicitly specified,
	// otherwise to an implementation-defined value.
	Requests Capacity `json:"requests"`
}

// Capacity is an object struct that represents the volume capacity.
type Capacity struct {
	// Volume size, in GB format: 'xGi'.
	Storage string `json:"storage"`
}

// Status is an object struct that represents the volume capacity.
type Status struct {
	// AccessModes contains the actual access modes the volume backing the PVC has.
	Capacity Capacity `json:"capacity"`
	// AccessModes contains the actual access modes the volume backing the PVC has.
	AccessModes []string `json:"accessModes"`
	// Phase represents the current phase of persistentVolumeClaim.
	//   pending: used for PersistentVolumeClaims that are not yet bound.
	//   Bound: used for PersistentVolumeClaims that are bound.
	//   Lost: used for PersistentVolumeClaims that lost their underlying.
	Phase string `json:"phase"`
}

type commonResult struct {
	golangsdk.Result
}

// CreateResult represents a result of the Create method.
type CreateResult struct {
	commonResult
}

func (r commonResult) Extract() (*PersistentVolumeClaim, error) {
	var s PersistentVolumeClaim
	err := r.ExtractInto(&s)
	return &s, err
}

// PersistentVolumeClaimPage represents the result of a List method.
type PersistentVolumeClaimPage struct {
	pagination.SinglePageBase
}

// ExtractPersistentVolumeClaims is a method to interpret the PersistentVolumeClaimPage as a PVC list.
func ExtractPersistentVolumeClaims(r pagination.Page) ([]PersistentVolumeClaim, error) {
	var s []PersistentVolumeClaim
	err := r.(PersistentVolumeClaimPage).Result.ExtractIntoSlicePtr(&s, "items")
	return s, err
}

// DeleteResult represents a result of the Create method.
type DeleteResult struct {
	commonResult
}

// Extract is a method which to extract the response to a PVC object.
func (r DeleteResult) Extract() ([]PersistentVolumeClaim, error) {
	var s []PersistentVolumeClaim
	err := r.ExtractInto(&s)
	return s, err
}
