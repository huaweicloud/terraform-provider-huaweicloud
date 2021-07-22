/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package certificates

import "github.com/huaweicloud/golangsdk"

type ListResult struct {
	golangsdk.Result
}

// CertificateList the struct returned by the query list.
type CertificateList struct {
	// Lists the certificates. For details, see Table 4.
	Certificates []Certificate `json:"certificates"`
	// Specifies the number of certificates.
	InstanceNum int `json:"instance_num"`
}

// Certificate the certificate detail returned by the query list.
type Certificate struct {
	Id           string `json:"id"`
	TenantId     string `json:"tenant_id"`
	AdminStateUp bool   `json:"admin_state_up"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Type         string `json:"type"`
	Domain       string `json:"domain"`
	PrivateKey   string `json:"private_key"`
	Certificate  string `json:"certificate"`
	ExpireTime   string `json:"expire_time"`
	CreateTime   string `json:"create_time"`
	UpdateTime   string `json:"update_time"`
}

// Extract extract the `Result.Body` to CertificateList
func (r ListResult) Extract() (*CertificateList, error) {
	var s CertificateList
	err := r.ExtractInto(&s)
	return &s, err
}
