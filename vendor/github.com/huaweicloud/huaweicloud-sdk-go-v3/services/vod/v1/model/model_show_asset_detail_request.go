package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type ShowAssetDetailRequest struct {

	// 媒资ID。
	AssetId string `json:"asset_id"`

	// 查询的信息类型。 - 为空时表示查询所有信息。 - 不为空时支持同时查询一个或者多个类型的信息，取值如下： - - base_info：媒资基本信息。 - - transcode_info：转码结果信息。 - - thumbnail_info：截图结果信息。 - - review_info：审核结果信息。
	Categories *[]ShowAssetDetailRequestCategories `json:"categories,omitempty"`
}

func (o ShowAssetDetailRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAssetDetailRequest struct{}"
	}

	return strings.Join([]string{"ShowAssetDetailRequest", string(data)}, " ")
}

type ShowAssetDetailRequestCategories struct {
	value string
}

type ShowAssetDetailRequestCategoriesEnum struct {
	BASE_INFO      ShowAssetDetailRequestCategories
	TRANSCODE_INFO ShowAssetDetailRequestCategories
	THUMBNAIL_INFO ShowAssetDetailRequestCategories
	REVIEW_INFO    ShowAssetDetailRequestCategories
}

func GetShowAssetDetailRequestCategoriesEnum() ShowAssetDetailRequestCategoriesEnum {
	return ShowAssetDetailRequestCategoriesEnum{
		BASE_INFO: ShowAssetDetailRequestCategories{
			value: "base_info",
		},
		TRANSCODE_INFO: ShowAssetDetailRequestCategories{
			value: "transcode_info",
		},
		THUMBNAIL_INFO: ShowAssetDetailRequestCategories{
			value: "thumbnail_info",
		},
		REVIEW_INFO: ShowAssetDetailRequestCategories{
			value: "review_info",
		},
	}
}

func (c ShowAssetDetailRequestCategories) Value() string {
	return c.value
}

func (c ShowAssetDetailRequestCategories) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowAssetDetailRequestCategories) UnmarshalJSON(b []byte) error {
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
