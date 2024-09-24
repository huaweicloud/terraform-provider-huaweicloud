package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateDeviceAuthorizerResponse Response Object
type CreateDeviceAuthorizerResponse struct {

	// **参数说明**：自定义鉴权ID。
	AuthorizerId *string `json:"authorizer_id,omitempty"`

	// **参数说明**：自定义鉴权器名称，同一租户下的自定义鉴权器名称不能重复。 **取值范围**：长度不超过128，只允许字母、数字、下划线（_）、连接符（-）的组合。
	AuthorizerName *string `json:"authorizer_name,omitempty"`

	// **参数说明**：函数名称。
	FuncName *string `json:"func_name,omitempty"`

	// **参数说明**：函数的URN（Uniform Resource Name），唯一标识函数，即自定义鉴权器对应的处理函数地址。
	FuncUrn *string `json:"func_urn,omitempty"`

	// **参数说明**：是否启动签名校验，启动签名校验后不满足签名要求的鉴权信息将被拒绝，以减少无效的函数调用。推荐用户进行安全的签名校验，默认开启。
	SigningEnable *bool `json:"signing_enable,omitempty"`

	// **参数说明**：签名校验的Key值，开启签名校验时使用。 **取值范围**：长度不超过128，只允许字母、数字、下划线（_）、连接符（-）的组合。
	SigningToken *string `json:"signing_token,omitempty"`

	// **参数说明**：签名校验的公钥，开启签名校验时使用。用于认证设备携带的签名信息是否正确。
	SigningPublicKey *string `json:"signing_public_key,omitempty"`

	// **参数说明**：是否为默认的鉴权方式，默认为false。
	DefaultAuthorizer *bool `json:"default_authorizer,omitempty"`

	// **参数说明**：是否激活该鉴权方式 - ACTIVE：该鉴权为激活状态。 - INACTIVE：该鉴权为停用状态。
	Status *string `json:"status,omitempty"`

	// **参数说明**：是否开启缓存，默认为false，设备为true时，当设备入参（username，clientId，password，以及证书信息，函数urn）不变时，当缓存结果存在时，将直接使用缓存结果，建议在调试时设置为false，生产时设置为true，避免频繁调用函数。
	CacheEnable *bool `json:"cache_enable,omitempty"`

	// 在物联网平台进行自定义鉴权相关操作的时间。格式：yyyyMMdd'T'HHmmss'Z'，如：20151212T121212Z。
	CreateTime *string `json:"create_time,omitempty"`

	// 在物联网平台更新自定义鉴权相关操作的时间。格式：yyyyMMdd'T'HHmmss'Z'，如：20151212T121212Z。
	UpdateTime     *string `json:"update_time,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateDeviceAuthorizerResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateDeviceAuthorizerResponse struct{}"
	}

	return strings.Join([]string{"CreateDeviceAuthorizerResponse", string(data)}, " ")
}
