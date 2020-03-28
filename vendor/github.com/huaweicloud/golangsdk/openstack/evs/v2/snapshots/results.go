package snapshots

import (
	"encoding/json"
	"time"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// Snapshot contains all the information associated with a Cinder Snapshot.
type Snapshot struct {
	// Unique identifier.
	ID string `json:"id"`

	// Date created.
	CreatedAt time.Time `json:"-"`

	// Date updated.
	UpdatedAt time.Time `json:"-"`

	// Display name.
	Name string `json:"name"`

	// Display description.
	Description string `json:"description"`

	// ID of the Volume from which this Snapshot was created.
	VolumeID string `json:"volume_id"`

	// Currect status of the Snapshot.
	Status string `json:"status"`

	// Size of the Snapshot, in GB.
	Size int `json:"size"`

	// User-defined key-value pairs.
	Metadata map[string]string `json:"metadata"`
}

func (r *Snapshot) UnmarshalJSON(b []byte) error {
	type tmp Snapshot
	var s struct {
		tmp
		CreatedAt golangsdk.JSONRFC3339MilliNoZ `json:"created_at"`
		UpdatedAt golangsdk.JSONRFC3339MilliNoZ `json:"updated_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Snapshot(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)
	r.UpdatedAt = time.Time(s.UpdatedAt)

	return err
}

type commonResult struct {
	golangsdk.Result
}

// CreateResult contains the response body and error from a Create request.
type CreateResult struct {
	commonResult
}

// GetResult contains the response body and error from a Get request.
type GetResult struct {
	commonResult
}

// UpdateResult contains the response body and error from an Update request.
type UpdateResult struct {
	commonResult
}

// DeleteResult contains the response body and error from a Delete request.
type DeleteResult struct {
	golangsdk.ErrResult
}

// SnapshotPage is a pagination.Pager that is returned from a call to the List function.
type SnapshotPage struct {
	pagination.SinglePageBase
}

// IsEmpty returns true if a SnapshotPage contains no Snapshots.
func (r SnapshotPage) IsEmpty() (bool, error) {
	volumes, err := ExtractSnapshots(r)
	return len(volumes) == 0, err
}

// ExtractSnapshots extracts and returns Snapshots. It is used while iterating over a snapshots.List call.
func ExtractSnapshots(r pagination.Page) ([]Snapshot, error) {
	var s struct {
		Snapshots []Snapshot `json:"snapshots"`
	}
	err := (r.(SnapshotPage)).ExtractInto(&s)
	return s.Snapshots, err
}

// Extract will get the Snapshot object out of the commonResult object.
func (r commonResult) Extract() (*Snapshot, error) {
	var s struct {
		Snapshot *Snapshot `json:"snapshot"`
	}
	err := r.ExtractInto(&s)
	return s.Snapshot, err
}
