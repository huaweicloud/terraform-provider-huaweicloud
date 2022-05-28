package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowAssetTempAuthorityResponse struct {

	// 带授权签名字符串的URL。具体调用示例请参见[示例2：媒资分段上传（20M以上）](https://support.huaweicloud.com/api-vod/vod_04_0216.html)。  示例：https://{obs_domain}/{bucket}?AWSAccessKeyId={AccessKeyID}&Expires={ExpiresValue}&Signature={Signature}
	SignStr        *string `json:"sign_str,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowAssetTempAuthorityResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAssetTempAuthorityResponse struct{}"
	}

	return strings.Join([]string{"ShowAssetTempAuthorityResponse", string(data)}, " ")
}
