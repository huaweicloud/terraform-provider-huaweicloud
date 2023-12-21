package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type CreateUrlAuthchainReq struct {

	// 播放域名或推流域名
	Domain string `json:"domain"`

	// 域名类型
	DomainType CreateUrlAuthchainReqDomainType `json:"domain_type"`

	// 流名称，与推流或播放地址中的StreamName一致。
	Stream string `json:"stream"`

	// 应用名称，与推流或播放地址中的AppName一致。
	App string `json:"app"`

	// 鉴权方式C必选。 检查级别。LiveID由AppName和StreamName组成,即\"<app_name>/<stream_name>\"。 包含如下取值： - 3：只检查LiveID是否匹配，不检查鉴权URL是否过期。 - 5：检查LiveID是否匹配，Timestamp是否超时。
	CheckLevel *int32 `json:"check_level,omitempty"`

	// 用户定义的有效访问时间起始点；例如：2006-01-02T15:04:05Z07:00 不传或为空表示当前时间
	StartTime *string `json:"start_time,omitempty"`
}

func (o CreateUrlAuthchainReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateUrlAuthchainReq struct{}"
	}

	return strings.Join([]string{"CreateUrlAuthchainReq", string(data)}, " ")
}

type CreateUrlAuthchainReqDomainType struct {
	value string
}

type CreateUrlAuthchainReqDomainTypeEnum struct {
	PULL CreateUrlAuthchainReqDomainType
	PUSH CreateUrlAuthchainReqDomainType
}

func GetCreateUrlAuthchainReqDomainTypeEnum() CreateUrlAuthchainReqDomainTypeEnum {
	return CreateUrlAuthchainReqDomainTypeEnum{
		PULL: CreateUrlAuthchainReqDomainType{
			value: "pull",
		},
		PUSH: CreateUrlAuthchainReqDomainType{
			value: "push",
		},
	}
}

func (c CreateUrlAuthchainReqDomainType) Value() string {
	return c.value
}

func (c CreateUrlAuthchainReqDomainType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateUrlAuthchainReqDomainType) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}
