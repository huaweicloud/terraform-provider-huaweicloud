package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type VulInfoCveList struct {

	// CVE ID
	CveId *string `json:"cve_id,omitempty"`

	// CVSS分值
	Cvss *float32 `json:"cvss,omitempty"`
}

func (o VulInfoCveList) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VulInfoCveList struct{}"
	}

	return strings.Join([]string{"VulInfoCveList", string(data)}, " ")
}
