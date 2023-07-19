package certificates

import "github.com/chnsz/golangsdk"

type CertOpts struct {
	// The certificate name.
	Name string `json:"name" required:"true"`
	// The certificate content.
	Content string `json:"cert_content" required:"true"`
	// The private key of the certificate.
	PrivateKey string `json:"private_key" required:"true"`
	// The certificate type. The valid values are as follows:
	// + instance
	// + global
	Type string `json:"type,omitempty"`
	// The dedicated instance ID to which the certificate belongs.
	// If the certificate type is global and the instance ID is omitted, the value 'common' will be used by default.
	InstanceId string `json:"instance_id,omitempty"`
	// The trusted root certificate (CA).
	TrustedRootCA string `json:"trusted_root_ca,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method used to create a new certificate using given parameters.
func Create(client *golangsdk.ServiceClient, opts CertOpts) (*Certificate, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Certificate
	_, err = client.Post(rootURL(client), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Get is a method to obtain an existing certificate by its ID.
func Get(client *golangsdk.ServiceClient, certificateId string) (*Certificate, error) {
	var r Certificate
	_, err := client.Get(resourceURL(client, certificateId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Update is a method used to update the configuration of an existing certificate using given parameters.
func Update(client *golangsdk.ServiceClient, certificateId string, opts CertOpts) (*Certificate, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Certificate
	_, err = client.Put(resourceURL(client, certificateId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Delete is a method to delete an existing certificate using its ID.
func Delete(client *golangsdk.ServiceClient, certificateId string) error {
	_, err := client.Delete(resourceURL(client, certificateId), nil)
	return err
}
