package certificates

import (
	"github.com/chnsz/golangsdk"
)

/*
Part 1:
The response of the import operation.
*/
// The struct defines the information about the imported certificate.
type CertificateImportInfo struct {
	CertificateId string `json:"certificate_id"`
}

type ImportResult struct {
	golangsdk.Result
}

func (r ImportResult) Extract() (*CertificateImportInfo, error) {
	var s CertificateImportInfo
	err := r.ExtractInto(&s)
	return &s, err
}

/*
Part 2:
The response from query escrowed certificate.
*/
// This struct defines the authentication information of the domain.
// API that will be used to Obtain Certificate Information
type Authentification struct {
	RecordName  string `json:"record_name"`
	RecordType  string `json:"record_type"`
	RecordValue string `json:"record_value"`
	Domain      string `json:"domain"`
}

// The struct defines the detail information about the escrow certificate information.
type CertificateEscrowInfo struct {
	Id                  string             `json:"id"`
	Status              string             `json:"status"`
	OrderId             string             `json:"order_id"`
	Name                string             `json:"name"`
	CertificateType     string             `json:"type"`
	Brand               string             `json:"brand"`
	PushSupport         string             `json:"push_support"`
	RevokeReason        string             `json:"revoke_reason"`
	SignatureAlgrithm   string             `json:"signature_algrithm"`
	IssueTime           string             `json:"issue_time"`
	NotBefore           string             `json:"not_before"`
	NotAfter            string             `json:"not_after"`
	ValidityPeriod      int                `json:"validity_period,omitempty"`
	ValidationMethod    string             `json:"validation_method"`
	DomainType          string             `json:"domain_type"`
	MultiDomainType     string             `json:"multi_domain_type"`
	Domain              string             `json:"domain"`
	Sans                string             `json:"sans"`
	DomainCount         int                `json:"domain_count,omitempty"`
	WildcardCount       int                `json:"wildcard_count,omitempty"`
	Fingerprint         string             `json:"fingerprint"`
	EnterpriseProjectID string             `json:"enterprise_project_id"`
	Authentifications   []Authentification `json:"authentification,omitempty"`
}

type GetResult struct {
	golangsdk.Result
}

func (r GetResult) Extract() (*CertificateEscrowInfo, error) {
	var s CertificateEscrowInfo
	err := r.ExtractInto(&s)
	return &s, err
}

/*
-- Part 3:
-- The response from exported certificate.
*/
// The struct defines the detail information about the imported certificate.
type CertificateDetail struct {
	Certificate      string `json:"certificate" required:"true"`
	CertificateChain string `json:"certificate_chain" required:"true"`
	PrivateKey       string `json:"private_key" required:"true"`
}

func (r ExportResult) Extract() (*CertificateDetail, error) {
	var s CertificateDetail
	err := r.ExtractInto(&s)
	return &s, err
}

type ExportResult struct {
	golangsdk.Result
}

/*
Part 4:
The response from pushing certificate.
*/
type PushResult struct {
	golangsdk.ErrResult
}

/*
Part 5:
The response from deleting certificate.
*/
type DeleteResult struct {
	golangsdk.ErrResult
}
