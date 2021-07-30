/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package certificates

import (
	"github.com/huaweicloud/golangsdk"
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
