package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// EventProcessResponseInfo 进程信息
type EventProcessResponseInfo struct {

	// 进程名称
	ProcessName *string `json:"process_name,omitempty"`

	// 进程文件路径
	ProcessPath *string `json:"process_path,omitempty"`

	// 进程id
	ProcessPid *int32 `json:"process_pid,omitempty"`

	// 进程用户id
	ProcessUid *int32 `json:"process_uid,omitempty"`

	// 运行进程的用户名
	ProcessUsername *string `json:"process_username,omitempty"`

	// 进程文件命令行
	ProcessCmdline *string `json:"process_cmdline,omitempty"`

	// 进程文件名
	ProcessFilename *string `json:"process_filename,omitempty"`

	// 进程启动时间
	ProcessStartTime *int64 `json:"process_start_time,omitempty"`

	// 进程组ID
	ProcessGid *int32 `json:"process_gid,omitempty"`

	// 进程有效组ID
	ProcessEgid *int32 `json:"process_egid,omitempty"`

	// 进程有效用户ID
	ProcessEuid *int32 `json:"process_euid,omitempty"`

	// 祖父进程文件路径
	AncestorProcessPath *string `json:"ancestor_process_path,omitempty"`

	// 祖父进程id
	AncestorProcessPid *int32 `json:"ancestor_process_pid,omitempty"`

	// 祖父进程文件命令行
	AncestorProcessCmdline *string `json:"ancestor_process_cmdline,omitempty"`

	// 父进程名称
	ParentProcessName *string `json:"parent_process_name,omitempty"`

	// 父进程文件路径
	ParentProcessPath *string `json:"parent_process_path,omitempty"`

	// 父进程id
	ParentProcessPid *int32 `json:"parent_process_pid,omitempty"`

	// 父进程用户id
	ParentProcessUid *int32 `json:"parent_process_uid,omitempty"`

	// 父进程文件命令行
	ParentProcessCmdline *string `json:"parent_process_cmdline,omitempty"`

	// 父进程文件名
	ParentProcessFilename *string `json:"parent_process_filename,omitempty"`

	// 父进程启动时间
	ParentProcessStartTime *int64 `json:"parent_process_start_time,omitempty"`

	// 父进程组ID
	ParentProcessGid *int32 `json:"parent_process_gid,omitempty"`

	// 父进程有效组ID
	ParentProcessEgid *int32 `json:"parent_process_egid,omitempty"`

	// 父进程有效用户ID
	ParentProcessEuid *int32 `json:"parent_process_euid,omitempty"`

	// 子进程名称
	ChildProcessName *string `json:"child_process_name,omitempty"`

	// 子进程文件路径
	ChildProcessPath *string `json:"child_process_path,omitempty"`

	// 子进程id
	ChildProcessPid *int32 `json:"child_process_pid,omitempty"`

	// 子进程用户id
	ChildProcessUid *int32 `json:"child_process_uid,omitempty"`

	// 子进程文件命令行
	ChildProcessCmdline *string `json:"child_process_cmdline,omitempty"`

	// 子进程文件名
	ChildProcessFilename *string `json:"child_process_filename,omitempty"`

	// 子进程启动时间
	ChildProcessStartTime *int64 `json:"child_process_start_time,omitempty"`

	// 子进程组ID
	ChildProcessGid *int32 `json:"child_process_gid,omitempty"`

	// 子进程有效组ID
	ChildProcessEgid *int32 `json:"child_process_egid,omitempty"`

	// 子进程有效用户ID
	ChildProcessEuid *int32 `json:"child_process_euid,omitempty"`

	// 虚拟化命令
	VirtCmd *string `json:"virt_cmd,omitempty"`

	// 虚拟化进程名称
	VirtProcessName *string `json:"virt_process_name,omitempty"`

	// 逃逸方式
	EscapeMode *string `json:"escape_mode,omitempty"`

	// 逃逸后后执行的命令
	EscapeCmd *string `json:"escape_cmd,omitempty"`

	// 进程启动文件hash
	ProcessHash *string `json:"process_hash,omitempty"`

	// 进程文件hash
	ProcessFileHash *string `json:"process_file_hash,omitempty"`

	// 父进程文件hash
	ParentProcessFileHash *string `json:"parent_process_file_hash,omitempty"`

	// 是否阻断成功，1阻断成功 0阻断失败
	Block *int32 `json:"block,omitempty"`
}

func (o EventProcessResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EventProcessResponseInfo struct{}"
	}

	return strings.Join([]string{"EventProcessResponseInfo", string(data)}, " ")
}
