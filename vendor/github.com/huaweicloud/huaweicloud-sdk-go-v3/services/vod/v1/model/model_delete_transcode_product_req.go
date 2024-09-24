package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type DeleteTranscodeProductReq struct {

	// 媒资Id
	AssetId *string `json:"asset_id,omitempty"`

	// GROUP：模板组级删除。 PRODUCT：产物级删除
	DeleteType *DeleteTranscodeProductReqDeleteType `json:"delete_type,omitempty"`

	// 删除的产物信息
	DeleteInfos *[]ProductGroupInfo `json:"delete_infos,omitempty"`
}

func (o DeleteTranscodeProductReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteTranscodeProductReq struct{}"
	}

	return strings.Join([]string{"DeleteTranscodeProductReq", string(data)}, " ")
}

type DeleteTranscodeProductReqDeleteType struct {
	value string
}

type DeleteTranscodeProductReqDeleteTypeEnum struct {
	GROUP   DeleteTranscodeProductReqDeleteType
	PRODUCT DeleteTranscodeProductReqDeleteType
}

func GetDeleteTranscodeProductReqDeleteTypeEnum() DeleteTranscodeProductReqDeleteTypeEnum {
	return DeleteTranscodeProductReqDeleteTypeEnum{
		GROUP: DeleteTranscodeProductReqDeleteType{
			value: "GROUP",
		},
		PRODUCT: DeleteTranscodeProductReqDeleteType{
			value: "PRODUCT",
		},
	}
}

func (c DeleteTranscodeProductReqDeleteType) Value() string {
	return c.value
}

func (c DeleteTranscodeProductReqDeleteType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *DeleteTranscodeProductReqDeleteType) UnmarshalJSON(b []byte) error {
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
