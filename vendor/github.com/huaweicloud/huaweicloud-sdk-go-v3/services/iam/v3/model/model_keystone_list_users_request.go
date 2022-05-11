package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneListUsersRequest struct {

	// IAM用户所属账号ID，获取方式请参见：[获取账号ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	DomainId *string `json:"domain_id,omitempty"`

	// 是否启IAM用户，true为启用，false为停用，默认为true。
	Enabled *bool `json:"enabled,omitempty"`

	// IAM用户名。
	Name *string `json:"name,omitempty"`

	// 密码过期时间，格式为：password_expires_at={operator}:{timestamp}。timestamp格式为：YYYY-MM-DDTHH:mm:ssZ。示例：  ``` password_expires_at=lt:2016-12-08T22:02:00Z ``` > - operator取值范围：lt，lte，gt，gte，eq，neq。 > - lt：过期时间小于timestamp。 > - lte：过期时间小于等于timestamp。 > - gt：过期时间大于timestamp。 > - gte：过期时间大于等于timestamp。 > - eq：过期时间等于timestamp。 > - neq：过期时间不等于timestamp。
	PasswordExpiresAt *string `json:"password_expires_at,omitempty"`
}

func (o KeystoneListUsersRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneListUsersRequest struct{}"
	}

	return strings.Join([]string{"KeystoneListUsersRequest", string(data)}, " ")
}
