package persistentvolumeclaims

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type PersistentVolumeClaim struct {
	Kind       string   `json:"kind"`
	ApiVersion string   `json:"apiVersion"`
	Metadata   MetaResp `json:"metadata"`
	Spec       SpecResp `json:"spec"`
	Status     Status   `json:"status"`
}

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

// The extra response of Persistent Volume are Capacity, FlexVolume, ClaimRef and PersistentVolumeReclaimPolicy.
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
	// The capacity of the storage.
	Capacity ResourceName `json:"capacity"`
	// PersistentVolumeClaim.
	FlexVolume FlexVolume `json:"flexVolume"`
	// Part of a bi-directional binding between persistentVolume and persistentVolumeClaim.
	ClaimRef ClaimRef `json:"claimRef"`
	// Specifies what happens to a persistent volume when released form its claim.
	PersistentVolumeReclaimPolicy string `json:"persistentVolumeReclaimPolicy"`
}

type FlexVolume struct {
	Driver  string  `json:"driver"`
	FsType  string  `json:"fsType"`
	Options Options `json:"options"`
}

type Options struct {
	// The type of the file system.
	FsType string `json:"fsType"`
	// ID of the volume.
	VolumeID string `json:"volumeID"`
	// The Shared path of the SFS and the SFS Turbo.
	DeviceMountPath string `json:"deviceMountPath"`
}

type ClaimRef struct {
	// Kind of the referent.
	Kind string `json:"kind"`
	// Namespace of the referent.
	Namespace string `json:"namespace"`
	// Name of the referent.
	Name string `json:"name"`
	// UID of the referent.
	UID string `json:"uid"`
	// API version of the referent.
	AapiVersion string `json:"apiVersion"`
	// Specifies resource version to which this reference is made, If any.
	ResourceVersion string `json:"resourceVersion"`
}

type Status struct {
	// Phase represents the current phase of persistentVolumeClaim.
	//     pending: used for PersistentVolumeClaims that are not yet bound.
	//     Bound: used for PersistentVolumeClaims that are bound.
	//     Lost: used for PersistentVolumeClaims that lost their underlying.
	Phase string `json:"phase"`
}

type commonResult struct {
	golangsdk.Result
}

type CreateResult struct {
	commonResult
}

func (r commonResult) Extract() (*PersistentVolumeClaim, error) {
	var s PersistentVolumeClaim
	err := r.ExtractInto(&s)
	return &s, err
}

type ListResp struct {
	PersistentVolumeClaim PersistentVolumeClaim `json:"persistentVolumeClaim"`
	PersistentVolume      PersistentVolumeClaim `json:"persistentVolume"`
}

type PersistentVolumeClaimPage struct {
	pagination.SinglePageBase
}

func ExtractPersistentVolumeClaims(r pagination.Page) ([]ListResp, error) {
	var s []ListResp
	err := r.(PersistentVolumeClaimPage).Result.ExtractIntoSlicePtr(&s, "")
	return s, err
}

type DeleteResult struct {
	commonResult
}

func (r DeleteResult) Extract() ([]PersistentVolumeClaim, error) {
	var s []PersistentVolumeClaim
	err := r.ExtractInto(&s)
	return s, err
}
