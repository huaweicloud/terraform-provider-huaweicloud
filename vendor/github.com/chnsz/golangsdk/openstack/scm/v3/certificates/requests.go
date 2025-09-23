package certificates

import (
	"strings"

	"github.com/chnsz/golangsdk"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// ImportOptsBuilder is the interface options structs have to satisfy in order
// to be used in the Import operation in this package.
type ImportOptsBuilder interface {
	ToCertificateImportMap() (map[string]interface{}, error)
}

// ImportOpts is the struct be used in the Import operation
type ImportOpts struct {
	Name                string `json:"name,omitempty" required:"true"`
	Certificate         string `json:"certificate" required:"true"`
	PrivateKey          string `json:"private_key" required:"true"`
	CertificateChain    string `json:"certificate_chain,omitempty"`
	EnterpriseProjectID string `json:"enterprise_project_id,omitempty"`
	EncCertificate      string `json:"enc_certificate,omitempty"`
	EncPrivateKey       string `json:"enc_private_key,omitempty"`
}

// ToCertificateImportMap casts a CreateOpts struct to a map.
func (opts ImportOpts) ToCertificateImportMap() (map[string]interface{}, error) {
	// Remove the blank content of both ends of the authentication value.
	// Otherwise, the service will go wrong.
	opts.PrivateKey = strings.Trim(opts.PrivateKey, "\r\n")
	opts.Certificate = strings.Trim(opts.Certificate, "\r\n")
	opts.CertificateChain = strings.Trim(opts.CertificateChain, "\r\n")

	return golangsdk.BuildRequestBody(opts, "")
}

// Import the certification into Huawei Cloud
func Import(c *golangsdk.ServiceClient, opts ImportOptsBuilder) (r ImportResult) {
	b, err := opts.ToCertificateImportMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(importURL(c), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200, 201},
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return
}

// PushOptsBuilder is the interface options structs have to satisfy in order
// to be used in the Push operation in this package.
type PushOptsBuilder interface {
	ToCertificatePushMap() (map[string]interface{}, error)
}

// PushOpts is the struct be used in the Import operation
type PushOpts struct {
	TargetProject string `json:"target_project"`
	TargetService string `json:"target_service" required:"true"`
}

// ToCertificatePushMap casts a PushOpts struct to a map.
func (opts PushOpts) ToCertificatePushMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Push the certification of imported to services
func Push(c *golangsdk.ServiceClient, id string, opts PushOptsBuilder) (r PushResult) {
	b, err := opts.ToCertificatePushMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = c.Post(pushURL(c, id), b, nil, &golangsdk.RequestOpts{
		OkCodes:     []int{204},
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return
}

// Obtain information about the imported certificate by ID.
// Contain no certificate key or private key.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return
}

// Get the certification key、private key and certification chain by id.
func Export(c *golangsdk.ServiceClient, id string) (r ExportResult) {
	body := map[string]interface{}{}
	_, r.Err = c.Post(exportURL(c, id), body, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200, 201},
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return
}

// Delete the imported certificate based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), nil)
	return
}
