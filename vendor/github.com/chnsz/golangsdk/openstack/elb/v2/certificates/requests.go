/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package certificates

import (
	"github.com/chnsz/golangsdk"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToCertificateListQuery() (string, error)
}

// ListOpts parameters used to query the certificate.
type ListOpts struct {
	// Specifies the ID of the certificate from which pagination query starts, that is,
	// the ID of the last certificate on the previous page. This parameter must be used with limit.
	Marker string `q:"marker"`
	// Specifies the number of certificates on each page.
	// If this parameter is not set, all certificates are queried by default.
	Limit int `q:"limit"`
	// Specifies the page direction. The value can be true or false, and the default value is false.
	// This parameter must be used with limit.
	PageReverse *bool `q:"page_reverse"`
	// Specifies the certificate ID.
	Id string `q:"id"`
	// Specifies the certificate name.
	// The value contains a maximum of 255 characters.
	Name string `q:"name"`
	// Provides supplementary information about the certificate.
	// The value contains a maximum of 255 characters.
	Description string `q:"description"`
	// Specifies the certificate type. The default value is server.
	// The value range varies depending on the protocol of the backend server group:
	// server: indicates the server certificate.
	// client: indicates the CA certificate.
	Type string `q:"type"`
	// Specifies the domain name associated with the server certificate. The default value is null.
	Domain string `q:"domain"`
	// Specifies the private key of the server certificate.
	PrivateKey string `q:"private_key"`
	// Specifies the public key of the server certificate or CA certificate used to authenticate the client.
	Certificate string `q:"certificate"`
	// Specifies the time when the certificate was created.
	// The UTC time is in YYYY-MM-DD HH:MM:SS format.
	CreateTime string `q:"create_time"`
	// Specifies the time when the certificate was updated.
	// The UTC time is in YYYY-MM-DD HH:MM:SS format.
	UpdateTime string `q:"update_time"`
}

// ToCertificateListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToCertificateListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// List query the certificate list
func List(c *golangsdk.ServiceClient, opts ListOptsBuilder) (ListResult, error) {
	var r ListResult
	var err error
	url := rootURL(c)
	if opts != nil {
		query, err := opts.ToCertificateListQuery()
		if err != nil {
			return r, err
		}
		url += query
	}
	_, r.Err = c.Get(url, &r.Body, nil)
	return r, err
}

// CreateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type CreateOptsBuilder interface {
	ToCertificateCreateMap() (map[string]interface{}, error)
}

// CreateOpts is the common options struct used in this package's Create
// operation.
type CreateOpts struct {
	// Specifies the certificate management status
	AdminStateUP bool `json:"admin_state_up,omitempty"`
	// Specifies the certificate name.
	// The value contains a maximum of 255 characters.
	Name string `json:"name,omitempty"`
	// Provides supplementary information about the certificate.
	// The value contains a maximum of 255 characters.
	Description string `json:"description,omitempty"`
	// Specifies the certificate type. The default value is server.
	// The value range varies depending on the protocol of the backend server group:
	// server: indicates the server certificate.
	// client: indicates the CA certificate.
	Type string `json:"type,omitempty"`
	// Specifies the domain name associated with the server certificate. The default value is null.
	Domain string `json:"domain,omitempty"`
	// Specifies the private key of the server certificate.
	PrivateKey string `json:"private_key,omitempty"`
	// Specifies the public key of the server certificate or CA certificate used to authenticate the client.
	Certificate string `json:"certificate" required:"true"`
	// Specifies the enterprise project id.
	EnterpriseProjectID string `json:"enterprise_project_id,omitempty"`
}

// ToCertificateCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToCertificateCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create is an operation which provisions a new loadbalancer based on the
// configuration defined in the CreateOpts struct. Once the request is
// validated and progress has started on the provisioning process, a
// CreateResult will be returned.
//
// Users with an admin role can create loadbalancers on behalf of other tenants by
// specifying a TenantID attribute different than their own.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToCertificateCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// Get retrieves a particular Loadbalancer based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

// UpdateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Update operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type UpdateOptsBuilder interface {
	ToCertificateUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts is the common options struct used in this package's Update
// operation.
type UpdateOpts struct {
	// Specifies the certificate management status
	AdminStateUP bool `json:"admin_state_up,omitempty"`
	// Specifies the certificate name.
	// The value contains a maximum of 255 characters.
	Name string `json:"name,omitempty"`
	// Provides supplementary information about the certificate.
	// The value contains a maximum of 255 characters.
	Description string `json:"description,omitempty"`
	// Specifies the domain name associated with the server certificate. The default value is null.
	Domain string `json:"domain,omitempty"`
	// Specifies the private key of the server certificate.
	PrivateKey string `json:"private_key,omitempty"`
	// Specifies the public key of the server certificate or CA certificate used to authenticate the client.
	Certificate string `json:"certificate,omitempty"`
}

// ToCertificateUpdateMap casts a UpdateOpts struct to a map.
func (opts UpdateOpts) ToCertificateUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update is an operation which modifies the attributes of the specified Certificate.
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOpts) (r UpdateResult) {
	b, err := opts.ToCertificateUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(resourceURL(c, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete will permanently delete a particular Certificate based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), nil)
	return
}
