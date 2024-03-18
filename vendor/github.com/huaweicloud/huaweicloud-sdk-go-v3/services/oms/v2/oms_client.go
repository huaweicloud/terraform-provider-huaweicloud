package v2

import (
	httpclient "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/oms/v2/model"
)

type OmsClient struct {
	HcClient *httpclient.HcHttpClient
}

func NewOmsClient(hcClient *httpclient.HcHttpClient) *OmsClient {
	return &OmsClient{HcClient: hcClient}
}

func OmsClientBuilder() *httpclient.HcHttpClientBuilder {
	builder := httpclient.NewHcHttpClientBuilder()
	return builder
}

// BatchUpdateTasks 批量更新任务
//
// 批量更新迁移任务，可指定单个迁移任务组下所有的迁移任务或通过迁移任务ID来执行。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *OmsClient) BatchUpdateTasks(request *model.BatchUpdateTasksRequest) (*model.BatchUpdateTasksResponse, error) {
	requestDef := GenReqDefForBatchUpdateTasks()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.BatchUpdateTasksResponse), nil
	}
}

// BatchUpdateTasksInvoker 批量更新任务
func (c *OmsClient) BatchUpdateTasksInvoker(request *model.BatchUpdateTasksRequest) *BatchUpdateTasksInvoker {
	requestDef := GenReqDefForBatchUpdateTasks()
	return &BatchUpdateTasksInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CheckPrefix 检查前缀是否在源端桶中存在
//
// 检查前缀是否在源端桶中存在
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *OmsClient) CheckPrefix(request *model.CheckPrefixRequest) (*model.CheckPrefixResponse, error) {
	requestDef := GenReqDefForCheckPrefix()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CheckPrefixResponse), nil
	}
}

// CheckPrefixInvoker 检查前缀是否在源端桶中存在
func (c *OmsClient) CheckPrefixInvoker(request *model.CheckPrefixRequest) *CheckPrefixInvoker {
	requestDef := GenReqDefForCheckPrefix()
	return &CheckPrefixInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateSyncEvents 创建同步事件
//
// 源端有对象需要进行同步时，调用该接口创建一个同步事件，系统将根据同步事件中包含的对象名称进行同步(目前只支持华北-北京四、华东-上海一地区)。
//
// Please refer to HUAWEI cloud API Explorer for details.
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

// CreateSyncTask 创建同步任务
//
// 创建同步任务，创建成功后，任务会被自动启动，不需要额外调用启动任务命令。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *OmsClient) CreateSyncTask(request *model.CreateSyncTaskRequest) (*model.CreateSyncTaskResponse, error) {
	requestDef := GenReqDefForCreateSyncTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateSyncTaskResponse), nil
	}
}

// CreateSyncTaskInvoker 创建同步任务
func (c *OmsClient) CreateSyncTaskInvoker(request *model.CreateSyncTaskRequest) *CreateSyncTaskInvoker {
	requestDef := GenReqDefForCreateSyncTask()
	return &CreateSyncTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateTask 创建迁移任务
//
// 创建迁移任务，创建成功后，任务会被自动启动，不需要额外调用启动任务命令。
//
// Please refer to HUAWEI cloud API Explorer for details.
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

// CreateTaskGroup 创建迁移任务组
//
// 创建迁移任务组，创建成功后，迁移任务组会自动创建迁移任务，不需要额外调用启动任务命令。
//
// Please refer to HUAWEI cloud API Explorer for details.
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

// DeleteSyncTask 删除同步任务
//
// 调用该接口删除同步任务。
// 正在同步的任务不允许删除，如果删除会返回失败；若要删除，请先行暂停任务。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *OmsClient) DeleteSyncTask(request *model.DeleteSyncTaskRequest) (*model.DeleteSyncTaskResponse, error) {
	requestDef := GenReqDefForDeleteSyncTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteSyncTaskResponse), nil
	}
}

// DeleteSyncTaskInvoker 删除同步任务
func (c *OmsClient) DeleteSyncTaskInvoker(request *model.DeleteSyncTaskRequest) *DeleteSyncTaskInvoker {
	requestDef := GenReqDefForDeleteSyncTask()
	return &DeleteSyncTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteTask 删除迁移任务
//
// 调用该接口删除迁移任务。
// 正在运行的任务不允许删除，如果删除会返回失败；若要删除，请先行暂停任务。
//
// Please refer to HUAWEI cloud API Explorer for details.
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

// DeleteTaskGroup 删除指定id的迁移任务组
//
// 删除指定的迁移任务组.
// 创建任务中、监控中、暂停中状态的任务不允许删除，如果删除会返回失败；若要删除，请先行暂停任务。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *OmsClient) DeleteTaskGroup(request *model.DeleteTaskGroupRequest) (*model.DeleteTaskGroupResponse, error) {
	requestDef := GenReqDefForDeleteTaskGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteTaskGroupResponse), nil
	}
}

// DeleteTaskGroupInvoker 删除指定id的迁移任务组
func (c *OmsClient) DeleteTaskGroupInvoker(request *model.DeleteTaskGroupRequest) *DeleteTaskGroupInvoker {
	requestDef := GenReqDefForDeleteTaskGroup()
	return &DeleteTaskGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListSyncTaskStatistic 查询指定ID的同步任务统计数据
//
// 查询指定ID同步任务的接收同步请求对象数、同步成功对象数、同步失败对象数、同步跳过对象数、同步成功对象容量统计数据(目前只支持华北-北京四、华东-上海一地区)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *OmsClient) ListSyncTaskStatistic(request *model.ListSyncTaskStatisticRequest) (*model.ListSyncTaskStatisticResponse, error) {
	requestDef := GenReqDefForListSyncTaskStatistic()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSyncTaskStatisticResponse), nil
	}
}

// ListSyncTaskStatisticInvoker 查询指定ID的同步任务统计数据
func (c *OmsClient) ListSyncTaskStatisticInvoker(request *model.ListSyncTaskStatisticRequest) *ListSyncTaskStatisticInvoker {
	requestDef := GenReqDefForListSyncTaskStatistic()
	return &ListSyncTaskStatisticInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListSyncTasks 查询同步任务列表
//
// 查询用户名下所有同步任务信息
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *OmsClient) ListSyncTasks(request *model.ListSyncTasksRequest) (*model.ListSyncTasksResponse, error) {
	requestDef := GenReqDefForListSyncTasks()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListSyncTasksResponse), nil
	}
}

// ListSyncTasksInvoker 查询同步任务列表
func (c *OmsClient) ListSyncTasksInvoker(request *model.ListSyncTasksRequest) *ListSyncTasksInvoker {
	requestDef := GenReqDefForListSyncTasks()
	return &ListSyncTasksInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListTaskGroup 查询迁移任务组列表
//
// 查询用户账户下的任务组信息
//
// Please refer to HUAWEI cloud API Explorer for details.
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

// ListTasks 查询迁移任务列表
//
// 查询用户账户下的所有任务信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
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

// RetryTaskGroup 对已经失败的指定id迁移任务组进行重启
//
// 当迁移任务组处于迁移失败状态时，调用该接口重启指定id的迁移任务组。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *OmsClient) RetryTaskGroup(request *model.RetryTaskGroupRequest) (*model.RetryTaskGroupResponse, error) {
	requestDef := GenReqDefForRetryTaskGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.RetryTaskGroupResponse), nil
	}
}

// RetryTaskGroupInvoker 对已经失败的指定id迁移任务组进行重启
func (c *OmsClient) RetryTaskGroupInvoker(request *model.RetryTaskGroupRequest) *RetryTaskGroupInvoker {
	requestDef := GenReqDefForRetryTaskGroup()
	return &RetryTaskGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowBucketList 查询桶列表
//
// 查询桶列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *OmsClient) ShowBucketList(request *model.ShowBucketListRequest) (*model.ShowBucketListResponse, error) {
	requestDef := GenReqDefForShowBucketList()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowBucketListResponse), nil
	}
}

// ShowBucketListInvoker 查询桶列表
func (c *OmsClient) ShowBucketListInvoker(request *model.ShowBucketListRequest) *ShowBucketListInvoker {
	requestDef := GenReqDefForShowBucketList()
	return &ShowBucketListInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowBucketObjects 查询桶对象列表
//
// 查询桶对象列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *OmsClient) ShowBucketObjects(request *model.ShowBucketObjectsRequest) (*model.ShowBucketObjectsResponse, error) {
	requestDef := GenReqDefForShowBucketObjects()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowBucketObjectsResponse), nil
	}
}

// ShowBucketObjectsInvoker 查询桶对象列表
func (c *OmsClient) ShowBucketObjectsInvoker(request *model.ShowBucketObjectsRequest) *ShowBucketObjectsInvoker {
	requestDef := GenReqDefForShowBucketObjects()
	return &ShowBucketObjectsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowBucketRegion 查询桶对应的region
//
// 查询桶对应的region
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *OmsClient) ShowBucketRegion(request *model.ShowBucketRegionRequest) (*model.ShowBucketRegionResponse, error) {
	requestDef := GenReqDefForShowBucketRegion()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowBucketRegionResponse), nil
	}
}

// ShowBucketRegionInvoker 查询桶对应的region
func (c *OmsClient) ShowBucketRegionInvoker(request *model.ShowBucketRegionRequest) *ShowBucketRegionInvoker {
	requestDef := GenReqDefForShowBucketRegion()
	return &ShowBucketRegionInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowCdnInfo 查桶对应的CDN信息
//
// 查桶对应的CDN信息
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *OmsClient) ShowCdnInfo(request *model.ShowCdnInfoRequest) (*model.ShowCdnInfoResponse, error) {
	requestDef := GenReqDefForShowCdnInfo()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowCdnInfoResponse), nil
	}
}

// ShowCdnInfoInvoker 查桶对应的CDN信息
func (c *OmsClient) ShowCdnInfoInvoker(request *model.ShowCdnInfoRequest) *ShowCdnInfoInvoker {
	requestDef := GenReqDefForShowCdnInfo()
	return &ShowCdnInfoInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowCloudType 查询所有支持的云厂商
//
// 查询所有支持的云厂商
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *OmsClient) ShowCloudType(request *model.ShowCloudTypeRequest) (*model.ShowCloudTypeResponse, error) {
	requestDef := GenReqDefForShowCloudType()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowCloudTypeResponse), nil
	}
}

// ShowCloudTypeInvoker 查询所有支持的云厂商
func (c *OmsClient) ShowCloudTypeInvoker(request *model.ShowCloudTypeRequest) *ShowCloudTypeInvoker {
	requestDef := GenReqDefForShowCloudType()
	return &ShowCloudTypeInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowRegionInfo 查询云厂商支持的reigon
//
// 查询云厂商支持的reigon
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *OmsClient) ShowRegionInfo(request *model.ShowRegionInfoRequest) (*model.ShowRegionInfoResponse, error) {
	requestDef := GenReqDefForShowRegionInfo()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowRegionInfoResponse), nil
	}
}

// ShowRegionInfoInvoker 查询云厂商支持的reigon
func (c *OmsClient) ShowRegionInfoInvoker(request *model.ShowRegionInfoRequest) *ShowRegionInfoInvoker {
	requestDef := GenReqDefForShowRegionInfo()
	return &ShowRegionInfoInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowSyncTask 查询指定ID的同步任务详情
//
// 查询指定ID的同步任务详情。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *OmsClient) ShowSyncTask(request *model.ShowSyncTaskRequest) (*model.ShowSyncTaskResponse, error) {
	requestDef := GenReqDefForShowSyncTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowSyncTaskResponse), nil
	}
}

// ShowSyncTaskInvoker 查询指定ID的同步任务详情
func (c *OmsClient) ShowSyncTaskInvoker(request *model.ShowSyncTaskRequest) *ShowSyncTaskInvoker {
	requestDef := GenReqDefForShowSyncTask()
	return &ShowSyncTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowTask 查询指定ID的任务详情
//
// 查询指定ID的任务详情。
//
// Please refer to HUAWEI cloud API Explorer for details.
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

// ShowTaskGroup 获取指定id的taskgroup信息
//
// 获取指定id的taskgroup信息
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *OmsClient) ShowTaskGroup(request *model.ShowTaskGroupRequest) (*model.ShowTaskGroupResponse, error) {
	requestDef := GenReqDefForShowTaskGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowTaskGroupResponse), nil
	}
}

// ShowTaskGroupInvoker 获取指定id的taskgroup信息
func (c *OmsClient) ShowTaskGroupInvoker(request *model.ShowTaskGroupRequest) *ShowTaskGroupInvoker {
	requestDef := GenReqDefForShowTaskGroup()
	return &ShowTaskGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StartSyncTask 启动同步任务
//
// 同步任务停止后，调用该接口以启动同步任务(目前只支持华北-北京四、华东-上海一地区)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *OmsClient) StartSyncTask(request *model.StartSyncTaskRequest) (*model.StartSyncTaskResponse, error) {
	requestDef := GenReqDefForStartSyncTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartSyncTaskResponse), nil
	}
}

// StartSyncTaskInvoker 启动同步任务
func (c *OmsClient) StartSyncTaskInvoker(request *model.StartSyncTaskRequest) *StartSyncTaskInvoker {
	requestDef := GenReqDefForStartSyncTask()
	return &StartSyncTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StartTask 启动迁移任务
//
// 迁移任务暂停或失败后，调用该接口以启动任务。
//
// Please refer to HUAWEI cloud API Explorer for details.
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

// StartTaskGroup 恢复指定id的迁移任务组
//
// 当迁移任务组处于暂停状态时，调用该接口启动指定id的迁移任务组。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *OmsClient) StartTaskGroup(request *model.StartTaskGroupRequest) (*model.StartTaskGroupResponse, error) {
	requestDef := GenReqDefForStartTaskGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StartTaskGroupResponse), nil
	}
}

// StartTaskGroupInvoker 恢复指定id的迁移任务组
func (c *OmsClient) StartTaskGroupInvoker(request *model.StartTaskGroupRequest) *StartTaskGroupInvoker {
	requestDef := GenReqDefForStartTaskGroup()
	return &StartTaskGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StopSyncTask 暂停同步任务
//
// 当同步任务处于同步中时，调用该接口停止任务(目前只支持华北-北京四、华东-上海一地区)。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *OmsClient) StopSyncTask(request *model.StopSyncTaskRequest) (*model.StopSyncTaskResponse, error) {
	requestDef := GenReqDefForStopSyncTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StopSyncTaskResponse), nil
	}
}

// StopSyncTaskInvoker 暂停同步任务
func (c *OmsClient) StopSyncTaskInvoker(request *model.StopSyncTaskRequest) *StopSyncTaskInvoker {
	requestDef := GenReqDefForStopSyncTask()
	return &StopSyncTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// StopTask 暂停迁移任务
//
// 当迁移任务处于迁移中时，调用该接口停止任务。
//
// Please refer to HUAWEI cloud API Explorer for details.
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

// StopTaskGroup 暂停指定id的迁移任务组
//
// 当迁移任务组处于创建任务中或监控中时，调用该接口暂停指定迁移任务组。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *OmsClient) StopTaskGroup(request *model.StopTaskGroupRequest) (*model.StopTaskGroupResponse, error) {
	requestDef := GenReqDefForStopTaskGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.StopTaskGroupResponse), nil
	}
}

// StopTaskGroupInvoker 暂停指定id的迁移任务组
func (c *OmsClient) StopTaskGroupInvoker(request *model.StopTaskGroupRequest) *StopTaskGroupInvoker {
	requestDef := GenReqDefForStopTaskGroup()
	return &StopTaskGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateBandwidthPolicy 更新任务带宽策略
//
// 当迁移任务未执行完成时，修改迁移任务的流量控制策略。
//
// Please refer to HUAWEI cloud API Explorer for details.
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

// UpdateTaskGroup 更新指定id的迁移任务组的流控策略
//
// 当迁移任务组未执行完成时，修改迁移任务组的流量控制策略。
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *OmsClient) UpdateTaskGroup(request *model.UpdateTaskGroupRequest) (*model.UpdateTaskGroupResponse, error) {
	requestDef := GenReqDefForUpdateTaskGroup()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateTaskGroupResponse), nil
	}
}

// UpdateTaskGroupInvoker 更新指定id的迁移任务组的流控策略
func (c *OmsClient) UpdateTaskGroupInvoker(request *model.UpdateTaskGroupRequest) *UpdateTaskGroupInvoker {
	requestDef := GenReqDefForUpdateTaskGroup()
	return &UpdateTaskGroupInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListApiVersions 查询API版本信息列表
//
// 查询对象存储迁移服务的API版本信息。
//
// Please refer to HUAWEI cloud API Explorer for details.
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
// Please refer to HUAWEI cloud API Explorer for details.
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
