package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ResourceBody struct {

	// 源站域名或源站IP，源站为IP类型时，仅支持IPv4，如需传入多个源站IP，以多个源站对象传入，除IP其他参数请保持一致，主源站最多支持15个源站IP对象，备源站最多支持15个源站IP对象；源站为域名类型时仅支持1个源站对象。不支持IP源站和域名源站混用。
	Sources []SourceWithPort `json:"sources"`
}

func (o ResourceBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResourceBody struct{}"
	}

	return strings.Join([]string{"ResourceBody", string(data)}, " ")
}
