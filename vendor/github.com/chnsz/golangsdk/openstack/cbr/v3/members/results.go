package members

import "github.com/chnsz/golangsdk/pagination"

// Member is the structure that represents the backup shared member detail.
type Member struct {
	// The status of the backup share.
	Status string `json:"status"`
	// The creation time of the backup sharing member.
	CreatedAt string `json:"created_at"`
	// The latest update time of the backup sharing member.
	UpdatedAt string `json:"updated_at"`
	// The backup ID.
	BackupId string `json:"backup_id"`
	// The latest update time of the backup shared member.
	ImageId string `json:"image_id"`
	// The ID of the project with which the backup is shared.
	DestProjectId string `json:"dest_project_id"`
	// The ID of the vault where the shared backup is stored.
	VaultId string `json:"vault_id"`
	// The ID of the backup shared member record.
	ID string `json:"id"`
}

// MemberPage is a single page maximum result representing a query by offset page.
type MemberPage struct {
	pagination.OffsetPageBase
}

// IsEmpty checks whether a MemberPage struct is empty.
func (b MemberPage) IsEmpty() (bool, error) {
	arr, err := ExtractMembers(b)
	return len(arr) == 0, err
}

// ExtractMembers is a method to extract the list of sharing members.
func ExtractMembers(r pagination.Page) ([]Member, error) {
	var s []Member
	err := r.(MemberPage).Result.ExtractIntoSlicePtr(&s, "members")
	return s, err
}
