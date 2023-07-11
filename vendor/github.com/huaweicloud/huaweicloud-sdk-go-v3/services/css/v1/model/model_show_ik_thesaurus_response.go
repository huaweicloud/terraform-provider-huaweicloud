package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowIkThesaurusResponse Response Object
type ShowIkThesaurusResponse struct {

	// 加载状态。  - Loaded表示加载成功。 - Loading表示正在加载中。 - Failed表示加载失败。
	Status *string `json:"status,omitempty"`

	// 存放词库文件的OBS桶。
	Bucket *string `json:"bucket,omitempty"`

	// 主词库文件对象。
	MainObj *string `json:"mainObj,omitempty"`

	// 停词词库文件对象。
	StopObj *string `json:"stopObj,omitempty"`

	// 同义词词库文件对象。
	SynonymObj *string `json:"synonymObj,omitempty"`

	// 词库最近更新时间。
	UpdateTime *string `json:"updateTime,omitempty"`

	// 更新详情。
	UpdateDetails *string `json:"updateDetails,omitempty"`

	// 指定配置自定义词库的集群ID。
	ClusterId *string `json:"clusterId,omitempty"`

	// 操作状态。
	OperateStatus *string `json:"operateStatus,omitempty"`

	// 词库的ID。
	Id             *string `json:"id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowIkThesaurusResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowIkThesaurusResponse struct{}"
	}

	return strings.Join([]string{"ShowIkThesaurusResponse", string(data)}, " ")
}
