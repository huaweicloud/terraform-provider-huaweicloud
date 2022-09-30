package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type LoadCustomThesaurusReq struct {

	// 词库文件存放的OBS桶（桶类型必须为标准存储或者低频存储，不支持归档存储）。
	BucketName string `json:"bucketName"`

	// 主词库文件对象，必须为UTF-8无BOM编码的文本文件，一行一个分词，文件大小最大支持100M。 mainObject, stopObject, synonymObject三个参数至少要填写一个。  >一次只能加载一个主词库，不支持同时加载多个主词库。
	MainObject string `json:"mainObject"`

	// 停词词库文件对象，必须为UTF-8无BOM编码的文本文件，一行一个分词，文件大小最大支持20M。  mainObject, stopObject, synonymObject三个参数至少要填写一个。
	StopObject string `json:"stopObject"`

	// 同义词词库文件，必须为UTF-8无BOM编码的文本文件，一行一组分词，文件大小最大支持20M。  mainObject, stopObject, synonymObject三个参数至少要填写一个。
	SynonymObject string `json:"synonymObject"`
}

func (o LoadCustomThesaurusReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "LoadCustomThesaurusReq struct{}"
	}

	return strings.Join([]string{"LoadCustomThesaurusReq", string(data)}, " ")
}
