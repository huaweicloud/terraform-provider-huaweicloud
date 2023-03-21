package imagecopy

import "github.com/chnsz/golangsdk"

// WithinRegionCopyOptsBuilder allows extensions to add parameters to the Create request.
type WithinRegionCopyOptsBuilder interface {
	// ToWithinRegionCopyMap Returns value that can be passed to json.Marshal
	ToWithinRegionCopyMap() (map[string]interface{}, error)
}

// WithinRegionCopyOpts represents options used to create an image.
type WithinRegionCopyOpts struct {
	// the name of the copy image
	Name string `json:"name" required:"true"`
	// description of the copy image
	Description string `json:"description,omitempty"`
	// the master key used for encrypting an image.
	CmkId string `json:"cmk_id,omitempty"`
	// Enterprise project ID
	EnterpriseProjectID string `json:"enterprise_project_id,omitempty"`
}

func (opts WithinRegionCopyOpts) ToWithinRegionCopyMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// WithinRegionCopy implements create image request.
func WithinRegionCopy(client *golangsdk.ServiceClient, imageId string, opts WithinRegionCopyOptsBuilder) (r JobResult) {
	b, err := opts.ToWithinRegionCopyMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(withinRegionCopyURL(client, imageId), b, &r.Body, nil)
	return
}

// CrossRegionCopyOptsBuilder allows extensions to add parameters to the Create request.
type CrossRegionCopyOptsBuilder interface {
	// ToCrossRegionCopyMap Returns value that can be passed to json.Marshal
	ToCrossRegionCopyMap() (map[string]interface{}, error)
}

// CrossRegionCopyOpts represents options used to create an image.
type CrossRegionCopyOpts struct {
	// the name of the copy image
	Name string `json:"name" required:"true"`
	// description of the copy image
	Description string `json:"description,omitempty"`
	// the target region name.
	TargetRegion string `json:"region" required:"true"`
	// the name of the project in the destination region.
	TargetProjectName string `json:"project_name" required:"true"`
	// the agency name.
	AgencyName string `json:"agency_name" required:"true"`
	// the ID of the vault
	VaultId string `json:"vault_id,omitempty"`
}

func (opts CrossRegionCopyOpts) ToCrossRegionCopyMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// CrossRegionCopy implements create image request.
func CrossRegionCopy(client *golangsdk.ServiceClient, imageId string, opts CrossRegionCopyOptsBuilder) (r JobResult) {
	b, err := opts.ToCrossRegionCopyMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(crossRegionCopyURL(client, imageId), b, &r.Body, nil)
	return
}
