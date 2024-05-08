package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Encryption 加密信息
type Encryption struct {

	// 密钥缓存时间。如果密钥不变，默认缓存七天
	KeyRotationIntervalSeconds *int32 `json:"key_rotation_interval_seconds,omitempty"`

	// 加密方式
	EncryptionMethod *EncryptionEncryptionMethod `json:"encryption_method,omitempty"`

	// 取值如下： - content：一个频道对应一个密钥 - profile：一个码率对应一个密钥  默认值：content
	Level *EncryptionLevel `json:"level,omitempty"`

	// 客户生成的DRM内容ID
	DrmContentId string `json:"drm_content_id"`

	// system_id枚举值
	SystemIds []EncryptionSystemIds `json:"system_ids"`

	// 增加到请求消息体header中的鉴权信息
	AuthInfo string `json:"auth_info"`

	// 获取密钥的DRM地址
	KmUrl string `json:"km_url"`
}

func (o Encryption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Encryption struct{}"
	}

	return strings.Join([]string{"Encryption", string(data)}, " ")
}

type EncryptionEncryptionMethod struct {
	value string
}

type EncryptionEncryptionMethodEnum struct {
	SAMPLE_AES EncryptionEncryptionMethod
	AES_128    EncryptionEncryptionMethod
}

func GetEncryptionEncryptionMethodEnum() EncryptionEncryptionMethodEnum {
	return EncryptionEncryptionMethodEnum{
		SAMPLE_AES: EncryptionEncryptionMethod{
			value: "SAMPLE-AES",
		},
		AES_128: EncryptionEncryptionMethod{
			value: "AES-128",
		},
	}
}

func (c EncryptionEncryptionMethod) Value() string {
	return c.value
}

func (c EncryptionEncryptionMethod) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *EncryptionEncryptionMethod) UnmarshalJSON(b []byte) error {
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

type EncryptionLevel struct {
	value string
}

type EncryptionLevelEnum struct {
	CONTENT EncryptionLevel
	PROFILE EncryptionLevel
}

func GetEncryptionLevelEnum() EncryptionLevelEnum {
	return EncryptionLevelEnum{
		CONTENT: EncryptionLevel{
			value: "content",
		},
		PROFILE: EncryptionLevel{
			value: "profile",
		},
	}
}

func (c EncryptionLevel) Value() string {
	return c.value
}

func (c EncryptionLevel) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *EncryptionLevel) UnmarshalJSON(b []byte) error {
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

type EncryptionSystemIds struct {
	value string
}

type EncryptionSystemIdsEnum struct {
	WIDEVINE   EncryptionSystemIds
	PLAY_READY EncryptionSystemIds
	FAIR_PLAY  EncryptionSystemIds
}

func GetEncryptionSystemIdsEnum() EncryptionSystemIdsEnum {
	return EncryptionSystemIdsEnum{
		WIDEVINE: EncryptionSystemIds{
			value: "Widevine",
		},
		PLAY_READY: EncryptionSystemIds{
			value: "PlayReady",
		},
		FAIR_PLAY: EncryptionSystemIds{
			value: "FairPlay",
		},
	}
}

func (c EncryptionSystemIds) Value() string {
	return c.value
}

func (c EncryptionSystemIds) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *EncryptionSystemIds) UnmarshalJSON(b []byte) error {
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
