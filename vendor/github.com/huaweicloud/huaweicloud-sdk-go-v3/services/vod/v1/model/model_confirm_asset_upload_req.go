package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type ConfirmAssetUploadReq struct {

	// 媒资ID。
	AssetId string `json:"asset_id"`

	// 上传状态。  取值如下： - CREATED：创建成功。 - FAILED：创建失败。 - CANCELLED：创建取消。
	Status ConfirmAssetUploadReqStatus `json:"status"`
}

func (o ConfirmAssetUploadReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ConfirmAssetUploadReq struct{}"
	}

	return strings.Join([]string{"ConfirmAssetUploadReq", string(data)}, " ")
}

type ConfirmAssetUploadReqStatus struct {
	value string
}

type ConfirmAssetUploadReqStatusEnum struct {
	CREATED   ConfirmAssetUploadReqStatus
	FAILED    ConfirmAssetUploadReqStatus
	CANCELLED ConfirmAssetUploadReqStatus
}

func GetConfirmAssetUploadReqStatusEnum() ConfirmAssetUploadReqStatusEnum {
	return ConfirmAssetUploadReqStatusEnum{
		CREATED: ConfirmAssetUploadReqStatus{
			value: "CREATED",
		},
		FAILED: ConfirmAssetUploadReqStatus{
			value: "FAILED",
		},
		CANCELLED: ConfirmAssetUploadReqStatus{
			value: "CANCELLED",
		},
	}
}

func (c ConfirmAssetUploadReqStatus) Value() string {
	return c.value
}

func (c ConfirmAssetUploadReqStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ConfirmAssetUploadReqStatus) UnmarshalJSON(b []byte) error {
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
