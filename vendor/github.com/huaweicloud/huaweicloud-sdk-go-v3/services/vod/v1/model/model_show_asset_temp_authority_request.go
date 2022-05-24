package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type ShowAssetTempAuthorityRequest struct {

	// 分段上传时调用OBS接口的HTTP方法，具体操作需要的HTTP方法请参考OBS的接口文档。  - 初始化上传任务：POST - 上传段：PUT - 合并段：POST - 取消段：DELETE - 列举已上传段：GET
	HttpVerb string `json:"http_verb"`

	// 桶名。  调用[创建媒资：上传方式](https://support.huaweicloud.com/api-vod/vod_04_0196.html)接口中返回的响应体中的target字段获得的bucket值。
	Bucket string `json:"bucket"`

	// 对象名。  调用[创建媒资：上传方式](https://support.huaweicloud.com/api-vod/vod_04_0196.html)接口中返回的响应体中的target字段获得的object值。
	ObjectKey string `json:"object_key"`

	// 文件类型对应的content-type，如MP4对应video/mp4。
	ContentType *string `json:"content_type,omitempty"`

	// 上传段时每段的MD5。
	ContentMd5 *string `json:"content_md5,omitempty"`

	// 每一个上传任务的id，是OBS进行初始段后OBS返回的。
	UploadId *string `json:"upload_id,omitempty"`

	// 上传段时每一段的id。  取值范围：[1,10000]。
	PartNumber *int32 `json:"part_number,omitempty"`
}

func (o ShowAssetTempAuthorityRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAssetTempAuthorityRequest struct{}"
	}

	return strings.Join([]string{"ShowAssetTempAuthorityRequest", string(data)}, " ")
}
