package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Encryption 对私钥进行加密存储的方式。
type Encryption struct {

	// 取值范围：“kms”或“default”。 - “default”为默认加密方式，适用于没有kms服务的局点。 - “kms”为采用kms服务加密方式。 若局点没有kms服务，请填“default”。
	Type EncryptionType `json:"type"`

	// kms密钥的名称。  - 若“type”为“kms”，则必须填入\"kms_key_name\"或\"kms_key_id\"。
	KmsKeyName *string `json:"kms_key_name,omitempty"`

	// kms密钥的ID。  - 若“type”为“kms”，则必须填入\"kms_key_name\"或\"kms_key_id\"。
	KmsKeyId *string `json:"kms_key_id,omitempty"`
}

func (o Encryption) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Encryption struct{}"
	}

	return strings.Join([]string{"Encryption", string(data)}, " ")
}

type EncryptionType struct {
	value string
}

type EncryptionTypeEnum struct {
	DEFAULT EncryptionType
	KMS     EncryptionType
}

func GetEncryptionTypeEnum() EncryptionTypeEnum {
	return EncryptionTypeEnum{
		DEFAULT: EncryptionType{
			value: "default",
		},
		KMS: EncryptionType{
			value: "kms",
		},
	}
}

func (c EncryptionType) Value() string {
	return c.value
}

func (c EncryptionType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *EncryptionType) UnmarshalJSON(b []byte) error {
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
