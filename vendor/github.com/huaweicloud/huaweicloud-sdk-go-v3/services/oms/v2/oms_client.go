package v2

import (
	http_client "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/oms/v2/model"
)

type OmsClient struct {
	HcClient *http_client.HcHttpClient
}

func NewOmsClient(hcClient *http_client.HcHttpClient) *OmsClient {
	return &OmsClient{HcClient: hcClient}
}

func OmsClientBuilder() *http_client.HcHttpClientBuilder {
	builder := http_client.NewHcHttpClientBuilder()
	return builder
}

// CreateSyncEvents 创建同步事件
//
// 源端有对象需要进行同步时，调用该接口创建一个同步事件，系统将根据同步事件中包含的对象名称进行同步(目前只支持华北-北京四、华东-上海一地区)。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *OmsClient) CreateSyncEvents(request *model.CreateSyncEventsRequest) (*model.CreateSyncEventsResponse, error) {
	requestDef := GenReqDefForCreateSyncEvents()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateSyncEventsResponse), nil
	}
}

// CreateSyncEventsInvoker 创建同步事件
func (c *OmsClient) CreateSyncEventsInvoker(request *model.CreateSyncEventsRequest) *CreateSyncEventsInvoker {
	requestDef := GenReqDefForCreateSyncEvents()
	return &CreateSyncEventsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateTask 创建迁移任务
//
// 创建迁移任务，创建成功后，任务会被自动启动，不需要额外调用启动任务命令。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *OmsClient) CreateTask(request *model.CreateTaskRequest) (*model.CreateTaskResponse, error) {
	requestDef := GenReqDefForCreateTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateTaskResponse), nil
	}
}

// CreateTaskInvoker 创建迁移任务
func (c *OmsClient) CreateTaskInvoker(request *model.CreateTaskRequest) *CreateTaskInvoker {
	requestDef := GenReqDefForCreateTask()
	return &CreateTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteTask 删除迁移任务
//
// 调用该接口删除迁移任务。
// 正在运行的任务不允许删除，如果删除会返回失败；若要删除，请先行暂停任务。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *OmsClient) DeleteTask(request *model.DeleteTaskRequest) (*model.DeleteTaskResponse, error) {
	requestDef := GenReqDefForDeleteTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteTaskResponse), nil
	}
}

// DeleteTaskInvoker 删除迁移任务
func (c *OmsClient) DeleteTaskInvoker(request *model.DeleteTaskRequest) *DeleteTaskInvoker {
	requestDef := GenReqDefForDeleteTask()
	return &DeleteTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListTasks 查询迁移任务列表
//
// 查询用户账户下的所有任务信息。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *OmsClient) ListTasks(request *model.ListTasksRequest) (*model.ListTasksResponse, error) {
	requestDef := GenReqDefForListTasks()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListTasksResponse), nil
	}
}

// ListTasksInvoker 查询迁移任务列表
func (c *OmsClient) ListTasksInvoker(request *model.ListTasksRequest) *ListTasksInvoker {
	requestDef := GenReqDefForListTasks()
	return &ListTasksInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowTask 查询指定ID的任务详情
//
// 查询指定ID的任务详情。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *OmsClient) ShowTask(request *model.ShowTaskRequest) (*model.ShowTaskResponse, error) {
	requestDef := GenReqDefForShowTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowTaskResponse), nil
	}
}

// ShowTaskInvoker 查询指定ID的任务详情
func (c *OmsClient) ShowTaskInvoker(request *model.ShowTaskRequest) *ShowTaskInvoker {
	requestDef := GenReqDefForShowTask()
	return &ShowTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StartTask 启动迁移任务
//
// 迁移任务暂停或失败后，调用该接口以启动任务。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *OmsClient) StartTask(request *model.StartTaskRequest) (*model.StartTaskResponse, error) {
	requestDef := GenReqDefForStartTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartTaskResponse), nil
	}
}

// StartTaskInvoker 启动迁移任务
func (c *OmsClient) StartTaskInvoker(request *model.StartTaskRequest) *StartTaskInvoker {
	requestDef := GenReqDefForStartTask()
	return &StartTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StopTask 暂停迁移任务
//
// 当迁移任务处于迁移中时，调用该接口停止任务。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *OmsClient) StopTask(request *model.StopTaskRequest) (*model.StopTaskResponse, error) {
	requestDef := GenReqDefForStopTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StopTaskResponse), nil
	}
}

// StopTaskInvoker 暂停迁移任务
func (c *OmsClient) StopTaskInvoker(request *model.StopTaskRequest) *StopTaskInvoker {
	requestDef := GenReqDefForStopTask()
	return &StopTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateBandwidthPolicy 更新任务带宽策略
//
// 当迁移任务未执行完成时，修改迁移任务的流量控制策略。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *OmsClient) UpdateBandwidthPolicy(request *model.UpdateBandwidthPolicyRequest) (*model.UpdateBandwidthPolicyResponse, error) {
	requestDef := GenReqDefForUpdateBandwidthPolicy()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateBandwidthPolicyResponse), nil
	}
}

// UpdateBandwidthPolicyInvoker 更新任务带宽策略
func (c *OmsClient) UpdateBandwidthPolicyInvoker(request *model.UpdateBandwidthPolicyRequest) *UpdateBandwidthPolicyInvoker {
	requestDef := GenReqDefForUpdateBandwidthPolicy()
	return &UpdateBandwidthPolicyInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateTaskGroup 创建迁移任务组
//
// 创建迁移任务组，创建成功后，迁移任务组会自动创建迁移任务，不需要额外调用启动任务命令（目前只支持华南-广州用户友好环境、西南-贵阳一、亚太-香港和亚太-新加坡地区）。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *OmsClient) CreateTaskGroup(request *model.CreateTaskGroupRequest) (*model.CreateTaskGroupResponse, error) {
	requestDef := GenReqDefForCreateTaskGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateTaskGroupResponse), nil
	}
}

// CreateTaskGroupInvoker 创建迁移任务组
func (c *OmsClient) CreateTaskGroupInvoker(request *model.CreateTaskGroupRequest) *CreateTaskGroupInvoker {
	requestDef := GenReqDefForCreateTaskGroup()
	return &CreateTaskGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteTaskGroup 删除指定ID的迁移任务组
//
// 删除指定的迁移任务组.（目前只支持华南-广州用户友好环境、西南-贵阳一、亚太-香港和亚太-新加坡地区）
// 创建任务中、监控中、暂停中状态的任务不允许删除，如果删除会返回失败；若要删除，请先行暂停任务。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *OmsClient) DeleteTaskGroup(request *model.DeleteTaskGroupRequest) (*model.DeleteTaskGroupResponse, error) {
	requestDef := GenReqDefForDeleteTaskGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteTaskGroupResponse), nil
	}
}

// DeleteTaskGroupInvoker 删除指定ID的迁移任务组
func (c *OmsClient) DeleteTaskGroupInvoker(request *model.DeleteTaskGroupRequest) *DeleteTaskGroupInvoker {
	requestDef := GenReqDefForDeleteTaskGroup()
	return &DeleteTaskGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListTaskGroup 查询迁移任务组列表
//
// 查询用户账户下的任务组信息（目前只支持华南-广州用户友好环境、西南-贵阳一、亚太-香港和亚太-新加坡地区）。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *OmsClient) ListTaskGroup(request *model.ListTaskGroupRequest) (*model.ListTaskGroupResponse, error) {
	requestDef := GenReqDefForListTaskGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListTaskGroupResponse), nil
	}
}

// ListTaskGroupInvoker 查询迁移任务组列表
func (c *OmsClient) ListTaskGroupInvoker(request *model.ListTaskGroupRequest) *ListTaskGroupInvoker {
	requestDef := GenReqDefForListTaskGroup()
	return &ListTaskGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// RetryTaskGroup 对已经失败的指定ID迁移任务组进行重启
//
// 当迁移任务组处于迁移失败状态时，调用该接口重启指定ID的迁移任务组（目前只支持华南-广州用户友好环境、西南-贵阳一、亚太-香港和亚太-新加坡地区）。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *OmsClient) RetryTaskGroup(request *model.RetryTaskGroupRequest) (*model.RetryTaskGroupResponse, error) {
	requestDef := GenReqDefForRetryTaskGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RetryTaskGroupResponse), nil
	}
}

// RetryTaskGroupInvoker 对已经失败的指定ID迁移任务组进行重启
func (c *OmsClient) RetryTaskGroupInvoker(request *model.RetryTaskGroupRequest) *RetryTaskGroupInvoker {
	requestDef := GenReqDefForRetryTaskGroup()
	return &RetryTaskGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowTaskGroup 获取指定ID的taskgroup信息
//
// 获取指定ID的taskgroup信息（目前只支持华南-广州用户友好环境、西南-贵阳一、亚太-香港和亚太-新加坡地区）。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *OmsClient) ShowTaskGroup(request *model.ShowTaskGroupRequest) (*model.ShowTaskGroupResponse, error) {
	requestDef := GenReqDefForShowTaskGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowTaskGroupResponse), nil
	}
}

// ShowTaskGroupInvoker 获取指定ID的taskgroup信息
func (c *OmsClient) ShowTaskGroupInvoker(request *model.ShowTaskGroupRequest) *ShowTaskGroupInvoker {
	requestDef := GenReqDefForShowTaskGroup()
	return &ShowTaskGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StartTaskGroup 恢复指定ID的迁移任务组
//
// 当迁移任务组处于暂停状态时，调用该接口启动指定ID的迁移任务组（目前只支持华南-广州用户友好环境、西南-贵阳一、亚太-香港和亚太-新加坡地区）。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *OmsClient) StartTaskGroup(request *model.StartTaskGroupRequest) (*model.StartTaskGroupResponse, error) {
	requestDef := GenReqDefForStartTaskGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartTaskGroupResponse), nil
	}
}

// StartTaskGroupInvoker 恢复指定ID的迁移任务组
func (c *OmsClient) StartTaskGroupInvoker(request *model.StartTaskGroupRequest) *StartTaskGroupInvoker {
	requestDef := GenReqDefForStartTaskGroup()
	return &StartTaskGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StopTaskGroup 暂停指定ID的迁移任务组
//
// 当迁移任务组处于创建任务中或监控中时，调用该接口暂停指定迁移任务组（目前只支持华南-广州用户友好环境、西南-贵阳一、亚太-香港和亚太-新加坡地区）。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *OmsClient) StopTaskGroup(request *model.StopTaskGroupRequest) (*model.StopTaskGroupResponse, error) {
	requestDef := GenReqDefForStopTaskGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StopTaskGroupResponse), nil
	}
}

// StopTaskGroupInvoker 暂停指定ID的迁移任务组
func (c *OmsClient) StopTaskGroupInvoker(request *model.StopTaskGroupRequest) *StopTaskGroupInvoker {
	requestDef := GenReqDefForStopTaskGroup()
	return &StopTaskGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateTaskGroup 更新指定ID的迁移任务组的流控策略
//
// 当迁移任务组未执行完成时，修改迁移任务组的流量控制策略（目前只支持华南-广州用户友好环境、西南-贵阳一、亚太-香港和亚太-新加坡地区）。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *OmsClient) UpdateTaskGroup(request *model.UpdateTaskGroupRequest) (*model.UpdateTaskGroupResponse, error) {
	requestDef := GenReqDefForUpdateTaskGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateTaskGroupResponse), nil
	}
}

// UpdateTaskGroupInvoker 更新指定ID的迁移任务组的流控策略
func (c *OmsClient) UpdateTaskGroupInvoker(request *model.UpdateTaskGroupRequest) *UpdateTaskGroupInvoker {
	requestDef := GenReqDefForUpdateTaskGroup()
	return &UpdateTaskGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListApiVersions 查询API版本信息列表
//
// 查询对象存储迁移服务的API版本信息。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *OmsClient) ListApiVersions(request *model.ListApiVersionsRequest) (*model.ListApiVersionsResponse, error) {
	requestDef := GenReqDefForListApiVersions()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListApiVersionsResponse), nil
	}
}

// ListApiVersionsInvoker 查询API版本信息列表
func (c *OmsClient) ListApiVersionsInvoker(request *model.ListApiVersionsRequest) *ListApiVersionsInvoker {
	requestDef := GenReqDefForListApiVersions()
	return &ListApiVersionsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowApiInfo 查询指定API版本信息
//
// 查询对象存储迁移服务指定API版本信息。
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *OmsClient) ShowApiInfo(request *model.ShowApiInfoRequest) (*model.ShowApiInfoResponse, error) {
	requestDef := GenReqDefForShowApiInfo()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowApiInfoResponse), nil
	}
}

// ShowApiInfoInvoker 查询指定API版本信息
func (c *OmsClient) ShowApiInfoInvoker(request *model.ShowApiInfoRequest) *ShowApiInfoInvoker {
	requestDef := GenReqDefForShowApiInfo()
	return &ShowApiInfoInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}
