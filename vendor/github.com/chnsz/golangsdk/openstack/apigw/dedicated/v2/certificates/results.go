package certificates

// Certificate is the structure that represents the APIG SSL certificate details.
type Certificate struct {
	// The certificate ID.
	ID string `json:"id"`
	// The certificate name.
	Name string `json:"name"`
	// The certificate type.
	Type string `json:"type"`
	// The instance ID to which the certificate belongs.
	// For global certificates, this value is not empty and "common" will be returned.
	InstanceId string `json:"instance_id"`
	// The project ID.
	ProjectId string `json:"project_id"`
	// The domain name.
	CommonName string `json:"common_name"`
	// The SAN domain list.
	SANs []string `json:"san"`
	// The expiration time.
	NotAfter string `json:"not_after"`
	// What signature algorithm the certificate uses.
	SignatureAlgorithm string `json:"signature_algorithm"`
	// The creation time of the certificate.
	CreatedAt string `json:"create_time"`
	// The update time of the certificate.
	UpdatedAt string `json:"update_time"`
	// Whether a trusted root certificate (CA) exists.
	// The value is true if trusted_root_ca exists in the bound certificate.
	// Defaults to false.
	IsHasTrustedRootCA bool `json:"is_has_trusted_root_ca"`
	// The certificate version
	Version int `json:"version"`
	// The Company or organization list to which the certificate belongs.
	Organizations []string `json:"organization"`
	// The department list.
	OrganizationalUnits []string `json:"organizational_unit"`
	// The city name.
	Locality []string `json:"locality"`
	// The state or province.
	State []string `json:"state"`
	// The country or region.
	Country []string `json:"country"`
	// The effective time
	NotBefore string `json:"not_before"`
	// The serial number.
	SerialNumber string `json:"serial_number"`
	// Issuer.
	Issuer []string `json:"issuer"`
}
