package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// EventFileResponseInfo 文件信息
type EventFileResponseInfo struct {

	// 文件路径
	FilePath *string `json:"file_path,omitempty"`

	// 文件别名
	FileAlias *string `json:"file_alias,omitempty"`

	// 文件大小
	FileSize *int32 `json:"file_size,omitempty"`

	// 文件最后一次修改时间
	FileMtime *int64 `json:"file_mtime,omitempty"`

	// 文件最后一次访问时间
	FileAtime *int64 `json:"file_atime,omitempty"`

	// 文件最后一次状态改变时间
	FileCtime *int64 `json:"file_ctime,omitempty"`

	// 文件hash,当前为sha256
	FileHash *string `json:"file_hash,omitempty"`

	// 文件md5
	FileMd5 *string `json:"file_md5,omitempty"`

	// 文件sha256
	FileSha256 *string `json:"file_sha256,omitempty"`

	// 文件类型
	FileType *string `json:"file_type,omitempty"`

	// 文件内容
	FileContent *string `json:"file_content,omitempty"`

	// 文件属性
	FileAttr *string `json:"file_attr,omitempty"`

	// 文件操作类型
	FileOperation *int32 `json:"file_operation,omitempty"`

	// 文件动作
	FileAction *string `json:"file_action,omitempty"`

	// 变更前后的属性
	FileChangeAttr *string `json:"file_change_attr,omitempty"`

	// 新文件路径
	FileNewPath *string `json:"file_new_path,omitempty"`

	// 文件描述
	FileDesc *string `json:"file_desc,omitempty"`

	// 文件关键字
	FileKeyWord *string `json:"file_key_word,omitempty"`

	// 是否目录
	IsDir *bool `json:"is_dir,omitempty"`

	// 文件句柄信息
	FdInfo *string `json:"fd_info,omitempty"`

	// 文件句柄数量
	FdCount *int32 `json:"fd_count,omitempty"`
}

func (o EventFileResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EventFileResponseInfo struct{}"
	}

	return strings.Join([]string{"EventFileResponseInfo", string(data)}, " ")
}
