package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type AssetInfo struct {

	// 媒资ID。
	AssetId *string `json:"asset_id,omitempty"`

	// 媒资状态。  取值如下： - UNCREATED：未创建（媒资ID不存在 ） - DELETED：已删除 - CANCELLED：上传取消 - SERVER_ERROR：上传失败（点播服务端故障） - UPLOAD_FAILED：上传失败（向OBS上传失败） - CREATING：创建中 - PUBLISHED：已发布 - WAITING_TRANSCODE：待发布（转码排队中） - TRANSCODING：待发布（转码中） - TRANSCODE_FAILED：待发布（转码失败） - TRANSCODE_SUCCEED：待发布（转码成功） - CREATED：待发布（未转码） - NO_ASSET：该媒资不存在 - DELETING：正在删除 - DELETE_FAILED：删除失败 - OBS_CREATING：OBS转存方式创建中 - OBS_CREATE_FAILED： OBS转存失败 - OBS_CREATE_SUCCESS： OBS转存成功
	Status *string `json:"status,omitempty"`

	// 媒资子状态或描述信息。 - 对于媒资异常场景，描述具体的异常原因。 - 对于正常场景，描述媒资的处理信息。
	Description *string `json:"description,omitempty"`

	BaseInfo *BaseInfo `json:"base_info,omitempty"`

	// 转码文件的播放信息。 - HLS或DASH：此数组的成员个数为n+1，n为转码输出路数。 - MP4：此数组的成员个数为n，n为转码输出路数。
	PlayInfoArray *[]PlayInfo `json:"play_info_array,omitempty"`
}

func (o AssetInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AssetInfo struct{}"
	}

	return strings.Join([]string{"AssetInfo", string(data)}, " ")
}
