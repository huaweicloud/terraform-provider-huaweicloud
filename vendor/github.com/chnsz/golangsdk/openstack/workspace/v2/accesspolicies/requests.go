package accesspolicies

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpts is the method that used to create an access policy.
type CreateOpts struct {
	// Access policy detail.
	Policy AccessPolicy `json:"policy" required:"true"`
	// List of policy objects.
	PolicyObjectsList []AccessPolicyObjectInfo `json:"policy_objects_list,omitempty"`
}

// AccessPolicy is the structure that represents the basic configuration of the access policy.
type AccessPolicy struct {
	// Policy name.
	// + PRIVATE_ACCESS: Private line access
	PolicyName string `json:"policy_name" required:"true"`
	// Blacklist type.
	// + INTERNET: Internet.
	BlacklistType string `json:"blacklist_type" required:"true"`
}

// AccessPolicyObjectInfo is the structure that represents the object list configuration.
type AccessPolicyObjectInfo struct {
	// Object ID of the blacklist.
	ObjectId string `json:"object_id" required:"true"`
	// Object type of the blacklist.
	// + USER
	// + USERGROUP
	ObjectType string `json:"object_type" required:"true"`
	// Obejct name of the blacklist.
	ObjectName string `json:"object_name,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method to create an access policy using given parameters.
func Create(c *golangsdk.ServiceClient, opts CreateOpts) error {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return err
	}

	_, err = c.Post(rootURL(c), b, nil, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}

// List is a method to query all access policies under a specified region.
func List(c *golangsdk.ServiceClient) ([]AccessPolicyDetailInfo, error) {
	var r policiesResp
	_, err := c.Get(rootURL(c), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return r.Policies, err
}

// ListObjectsOpts is the structure that used to query object list.
type ListObjectsOpts struct {
	// Policy ID.
	PolicyId string `json:"-" required:"true"`
	// Number of records to be queried.
	// Value range: 0–2000.
	// Default value: 10.
	Limit int `q:"limit"`
	// The offset number.
	// Value range: 0–1999.
	// Default value: 0
	Offset int `q:"offset"`
}

// ListObjects is the method that used to query object list under a specified policy using given parameters.
func ListObjects(c *golangsdk.ServiceClient, opts ListObjectsOpts) ([]AccessPolicyObject, error) {
	url := resourceURL(c, opts.PolicyId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pager := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := AccessPolicyPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	})
	pager.Headers = requestOpts.MoreHeaders
	pages, err := pager.AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractAccessPolicies(pages)
}

// UpdateOpts is the method that used to update an existing access policy.
type UpdateOpts struct {
	// The associated policy ID to which the objects belong.
	PolicyId string `json:"-" required:"true"`
	// List of policy objects.
	PolicyObjectsList []AccessPolicyObjectInfo `json:"policy_objects_list,omitempty"`
}

// UpdateObjects is a method to modify the policy objects using given parameters.
func UpdateObjects(c *golangsdk.ServiceClient, opts UpdateOpts) error {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return err
	}

	_, err = c.Put(resourceURL(c, opts.PolicyId), b, nil, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}

// DeleteOpts is the method that used to delete an existing access policy.
type DeleteOpts struct {
	// List of policy ID.
	PolicyIdList []string `json:"policy_id_list,omitempty"`
}

// Delete is a method to delete an existing access policy using given parameters.
func Delete(c *golangsdk.ServiceClient, opts DeleteOpts) error {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return err
	}

	_, err = c.DeleteWithBody(rootURL(c), b, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}
