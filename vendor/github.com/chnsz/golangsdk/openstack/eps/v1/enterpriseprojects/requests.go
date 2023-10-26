package enterpriseprojects

import (
	"github.com/chnsz/golangsdk"
)

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

type ListOpts struct {
	Name   string `q:"name"`
	ID     string `q:"id"`
	Status int    `q:"status"`
}

func (opts ListOpts) ToEnterpriseProjectListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

type ListOptsBuilder interface {
	ToEnterpriseProjectListQuery() (string, error)
}

func List(c *golangsdk.ServiceClient, opts ListOptsBuilder) (r ListResult) {
	url := rootURL(c)
	if opts != nil {
		query, err := opts.ToEnterpriseProjectListQuery()
		if err != nil {
			r.Err = err
		}
		url += query
	}

	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// CreateOpts allows to create a enterprise project using given parameters.
type CreateOpts struct {
	// A name can contain 1 to 64 characters.
	// Only letters, digits, underscores (_), and hyphens (-) are allowed.
	// The name must be unique in the domain and cannot include any form of
	// the word "default" ("deFaulT", for instance).
	Name string `json:"name" required:"true"`
	// A description can contain a maximum of 512 characters.
	Description string `json:"description"`
	// Specifies the enterprise project type.
	// The options are as follows:
	// poc: indicates a test project.
	// prod: indicates a commercial project.
	Type string `json:"type,omitempty"`
}

// Create accepts a CreateOpts struct and uses the values to create a new enterprise project.
func Create(client *golangsdk.ServiceClient, opts CreateOpts) (r CreatResult) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(rootURL(client), b, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return
}

// Get is a method to obtain the specified enterprise project by id.
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, id), &r.Body, nil)
	return
}

// Update accepts a CreateOpts struct and uses the values to Update a enterprise project.
func Update(client *golangsdk.ServiceClient, opts CreateOpts, id string) (r UpdateResult) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(resourceURL(client, id), b, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return
}

type ActionOpts struct {
	// enable: Enable an enterprise project.
	// disable: Disable an enterprise project.
	Action string `json:"action" required:"true"`
}

// Update accepts a ActionOpts struct and uses the values to enable or diaable a enterprise project.
func Action(client *golangsdk.ServiceClient, opts ActionOpts, id string) (r ActionResult) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(actionURL(client, id), b, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
		OkCodes:     []int{204},
	})
	return
}

type MigrateResourceOpts struct {
	ResourceId string `json:"resource_id" required:"true"`

	ResourceType string `json:"resource_type" required:"true"`
	// this filed is required when resource_type is bucket
	RegionId string `json:"region_id,omitempty"`

	// this filed is required when resource_type is region level
	ProjectId string `json:"project_id,omitempty"`

	// only support for EVS、EIP
	Associated *bool `json:"associated,omitempty"`
}

func Migrate(client *golangsdk.ServiceClient, opts MigrateResourceOpts, id string) (r MigrateResult) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(migrateURL(client, id), b, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
		OkCodes:     []int{204},
	})
	return
}

type ListResourcesOpts struct {
	// Target enterprise project ID.
	EnterpriseProjectId string `json:"-" required:"true"`

	ResourceTypes []string `json:"resource_types" required:"true"`

	Projects []string `json:"projects,omitempty"`

	Offset int32 `json:"offset,omitempty"`

	Limit int32 `json:"limit,omitempty"`

	Matches []Match `json:"matches,omitempty"`
}

type Match struct {
	Key string `json:"key" required:"true"`

	Value string `json:"value" required:"true"`
}

// ListAssociatedResources is a method that used to query associated resources for specified enterprise project using
// given parameters.
func ListAssociatedResources(client *golangsdk.ServiceClient, opts ListResourcesOpts) (*FilterResult, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r FilterResult
	_, err = client.Post(resourceFilterURL(client, opts.EnterpriseProjectId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &r, err
}
