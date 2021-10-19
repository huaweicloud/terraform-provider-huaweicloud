package resources

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
)

// CreateGroupAndUploadOpts is a structure which allows to create a new resource group and upload package file to the
// group using given parameters.
type CreateGroupAndUploadOpts struct {
	// List of OBS object paths. The OBS object path refers to the OBS object URL.
	Paths []string `json:"paths" required:"true"`
	// File type of a package group.
	//   jar: JAR file
	//   pyFile: User Python file
	//   file: User file
	//   modelFile: User AI model file
	// NOTE: If the same group of packages to be uploaded contains different file types, select file as the type of the
	// file to be uploaded.
	Kind string `json:"kind" required:"true"`
	// Name of the group to be created.
	Group string `json:"group" required:"true"`
	// Whether to upload resource packages in asynchronous mode.
	// The default value is false, indicating that the asynchronous mode is not used.
	// You are advised to upload resource packages in asynchronous mode.
	IsAsync bool `json:"is_async,omitempty"`
	// Resource tag.
	Tags []tags.ResourceTag `json:"tags,omitempty"`
}

// CreateGroupAndUpload is a method to create a new resource group and upload package from the specified OBS bucket.
func CreateGroupAndUpload(c *golangsdk.ServiceClient, opts CreateGroupAndUploadOpts) (*Group, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Post(rootURL(c), b, &rst.Body, nil)
	if err == nil {
		var r Group
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

// UploadOpts is a stucture which allows to upload resource package to the specified group using given parameters.
type UploadOpts struct {
	// List of OBS object paths. The OBS object path refers to the OBS object URL.
	Paths []string `json:"paths" required:"true"`
	// Name of a package group.
	Group string `json:"group" required:"true"`
}

// Upload is a method to upload resource package to the specified group.
func Upload(c *golangsdk.ServiceClient, typePath string, opts UploadOpts) (*Group, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Post(resourceURL(c, typePath), b, &rst.Body, nil)
	if err == nil {
		var r Group
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

// ResourceLocatedOpts is a structure which specify the resource package located.
type ResourceLocatedOpts struct {
	// Name of the package group returned when the resource package is uploaded.
	Group string `q:"group" required:"true"`
}

// Get is a method to obtain the resource (packages) from the specified group.
func Get(c *golangsdk.ServiceClient, resourceName string, opts ResourceLocatedOpts) (*Resource, error) {
	url := resourceURL(c, resourceName)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var rst golangsdk.Result
	_, err = c.Get(url, &rst.Body, nil)
	if err == nil {
		var r Resource
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

// ListOpts is a structure which allows to obtain resources groups using given parameters.
type ListOpts struct {
	// Specifies the file type.
	Kind string `q:"kind"`
	// Specifies a label for filtering.
	Tags string `q:"tags"`
}

// List is a method to obtain a list of the groups and resources.
func List(c *golangsdk.ServiceClient, opts ListOpts) (*ListResp, error) {
	url := rootURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var rst golangsdk.Result
	_, err = c.Get(url, &rst.Body, nil)
	if err == nil {
		var r ListResp
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

// UpdateOpts is a structure which allows to update package owner or group owner using given parameters.
type UpdateOpts struct {
	// New username. The name contains 5 to 32 characters, including only digits, letters, underscores (_), and
	// hyphens (-). It cannot start with a digit.
	NewOwner string `json:"new_owner,omitempty"`
	// Group name. The name contains a maximum of 64 characters. Only digits, letters, periods (.), underscores (_),
	// and hyphens (-) are allowed.
	GroupName string `json:"group_name,omitempty"`
	// Package name. The name can contain only digits, letters, underscores (_), exclamation marks (!), hyphens (-),
	// and periods (.), but cannot start with a period.
	// The length (including the file name extension) cannot exceed 128 characters.
	ResourceName string `json:"resource_name,omitempty"`
}

// UpdateOwner is a method to update package owner or group owner.
func UpdateOwner(c *golangsdk.ServiceClient, opts UpdateOpts) (*UpdateResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Put(resourceURL(c, "owner"), b, &rst.Body, nil)
	if err == nil {
		var r UpdateResp
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

// Delete is a method to remove resource package from the specified group.
func Delete(c *golangsdk.ServiceClient, resourceName string, opts ResourceLocatedOpts) *golangsdk.ErrResult {
	var rst golangsdk.ErrResult
	url := resourceURL(c, resourceName)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		rst.Err = err
		return &rst
	}
	url += query.String()

	_, rst.Err = c.Delete(url, nil)
	if err == nil {
		return &rst
	}
	return nil
}
