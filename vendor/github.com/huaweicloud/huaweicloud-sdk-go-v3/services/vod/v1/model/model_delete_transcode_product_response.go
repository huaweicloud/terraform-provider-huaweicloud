package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// DeleteTranscodeProductResponse Response Object
type DeleteTranscodeProductResponse struct {

	// 媒资ID
	AssetId *string `json:"asset_id,omitempty"`

	// SUCCESS：成功 FAIL：失败 PARTIAL_SUCCESS：部分成功
	Status *DeleteTranscodeProductResponseStatus `json:"status,omitempty"`

	// 删除成功的产物
	DeletedProducts *[]ProductGroupInfo `json:"deleted_products,omitempty"`

	// 删除失败的产物
	FailedProducts *[]ProductGroupDelFailInfo `json:"failed_products,omitempty"`
	HttpStatusCode int                        `json:"-"`
}

func (o DeleteTranscodeProductResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteTranscodeProductResponse struct{}"
	}

	return strings.Join([]string{"DeleteTranscodeProductResponse", string(data)}, " ")
}

type DeleteTranscodeProductResponseStatus struct {
	value string
}

type DeleteTranscodeProductResponseStatusEnum struct {
	SUCCESS         DeleteTranscodeProductResponseStatus
	FAIL            DeleteTranscodeProductResponseStatus
	PARTIAL_SUCCESS DeleteTranscodeProductResponseStatus
}

func GetDeleteTranscodeProductResponseStatusEnum() DeleteTranscodeProductResponseStatusEnum {
	return DeleteTranscodeProductResponseStatusEnum{
		SUCCESS: DeleteTranscodeProductResponseStatus{
			value: "SUCCESS",
		},
		FAIL: DeleteTranscodeProductResponseStatus{
			value: "FAIL",
		},
		PARTIAL_SUCCESS: DeleteTranscodeProductResponseStatus{
			value: "PARTIAL_SUCCESS",
		},
	}
}

func (c DeleteTranscodeProductResponseStatus) Value() string {
	return c.value
}

func (c DeleteTranscodeProductResponseStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *DeleteTranscodeProductResponseStatus) UnmarshalJSON(b []byte) error {
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
