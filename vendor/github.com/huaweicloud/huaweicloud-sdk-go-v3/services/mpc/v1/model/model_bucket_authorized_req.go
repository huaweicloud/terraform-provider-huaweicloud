package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type BucketAuthorizedReq struct {

	// 桶名
	Bucket string `json:"bucket"`

	// 操作标记，取值[0,1]，0表示取消授权，1表示授权
	Operation BucketAuthorizedReqOperation `json:"operation"`
}

func (o BucketAuthorizedReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BucketAuthorizedReq struct{}"
	}

	return strings.Join([]string{"BucketAuthorizedReq", string(data)}, " ")
}

type BucketAuthorizedReqOperation struct {
	value string
}

type BucketAuthorizedReqOperationEnum struct {
	E_0 BucketAuthorizedReqOperation
	E_1 BucketAuthorizedReqOperation
}

func GetBucketAuthorizedReqOperationEnum() BucketAuthorizedReqOperationEnum {
	return BucketAuthorizedReqOperationEnum{
		E_0: BucketAuthorizedReqOperation{
			value: "0",
		},
		E_1: BucketAuthorizedReqOperation{
			value: "1",
		},
	}
}

func (c BucketAuthorizedReqOperation) Value() string {
	return c.value
}

func (c BucketAuthorizedReqOperation) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *BucketAuthorizedReqOperation) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter != nil {
		val, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
		if err == nil {
			c.value = val.(string)
			return nil
		}
		return err
	} else {
		return errors.New("convert enum data to string error")
	}
}
