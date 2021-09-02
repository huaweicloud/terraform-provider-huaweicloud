package certificates

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type Certificate struct {
	// Certificate ID
	Id string `json:"id"`
	// Certificate Name
	Name string `json:"name"`
	// the time when the certificate expires in unix timestamp
	ExpireTime int `json:"expireTime"`
	// the time when the certificate is uploaded in unix timestamp
	TimeStamp int `json:"timestamp"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a certificate.
func (r commonResult) Extract() (*Certificate, error) {
	var response Certificate
	err := r.ExtractInto(&response)
	return &response, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Certificate.
type CreateResult struct {
	commonResult
}

// UpdateResult represents the result of a update operation. Call its Extract
// method to interpret it as a Certificate.
type UpdateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Certificate.
type GetResult struct {
	commonResult
}

// CertificateDetail reprersents the result of list certificates.
type CertificateDetail struct {
	// Certificate ID
	Id string `json:"id"`
	// Certificate Name
	Name string `json:"name"`
	// Certificate Id
	CertificateId string `json:"certificateid"`
	// Certificate name
	CertificateName string `json:"certificatename"`
	// the key of the certificate
	Timestamp int64 `json:"timestamp"`
	// the privite key of the certificate
	BindHosts []BindHost `json:"bind_host,omitempty"`
	// the time when the certificate expires in unix timestamp
	ExpireTime int `json:"expire_time"`
	// the expire status of the certificate
	// 0-not expired, 1-expired, 2-expired soon
	ExpStatus int `json:"exp_status"`
}

type BindHost struct {
	Id       string `json:"id"`
	Hostname string `json:"hostname"`
	WafType  string `json:"waf_type"`
	Mode     string `json:"mode"`
}

type CertificatePage struct {
	pagination.SinglePageBase
}

func ExtractCertificates(r pagination.Page) ([]CertificateDetail, error) {
	var s struct {
		Certificates []CertificateDetail `json:"items"`
	}
	err := (r.(CertificatePage)).ExtractInto(&s)
	return s.Certificates, err
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
