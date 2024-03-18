package v1

import (
	httpclient "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cpts/v1/model"
)

type CptsClient struct {
	HcClient *httpclient.HcHttpClient
}

func NewCptsClient(hcClient *httpclient.HcHttpClient) *CptsClient {
	return &CptsClient{HcClient: hcClient}
}

func CptsClientBuilder() *httpclient.HcHttpClientBuilder {
	builder := httpclient.NewHcHttpClientBuilder()
	return builder
}

// BatchUpdateTaskStatus 批量启停任务
//
// 批量启停任务
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) BatchUpdateTaskStatus(request *model.BatchUpdateTaskStatusRequest) (*model.BatchUpdateTaskStatusResponse, error) {
	requestDef := GenReqDefForBatchUpdateTaskStatus()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.BatchUpdateTaskStatusResponse), nil
	}
}

// BatchUpdateTaskStatusInvoker 批量启停任务
func (c *CptsClient) BatchUpdateTaskStatusInvoker(request *model.BatchUpdateTaskStatusRequest) *BatchUpdateTaskStatusInvoker {
	requestDef := GenReqDefForBatchUpdateTaskStatus()
	return &BatchUpdateTaskStatusInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateCase 创建用例（旧版）
//
// 创建用例（旧版）
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) CreateCase(request *model.CreateCaseRequest) (*model.CreateCaseResponse, error) {
	requestDef := GenReqDefForCreateCase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateCaseResponse), nil
	}
}

// CreateCaseInvoker 创建用例（旧版）
func (c *CptsClient) CreateCaseInvoker(request *model.CreateCaseRequest) *CreateCaseInvoker {
	requestDef := GenReqDefForCreateCase()
	return &CreateCaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateDirectory 创建目录
//
// 创建目录
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) CreateDirectory(request *model.CreateDirectoryRequest) (*model.CreateDirectoryResponse, error) {
	requestDef := GenReqDefForCreateDirectory()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateDirectoryResponse), nil
	}
}

// CreateDirectoryInvoker 创建目录
func (c *CptsClient) CreateDirectoryInvoker(request *model.CreateDirectoryRequest) *CreateDirectoryInvoker {
	requestDef := GenReqDefForCreateDirectory()
	return &CreateDirectoryInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateNewCase 创建用例
//
// 创建用例
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) CreateNewCase(request *model.CreateNewCaseRequest) (*model.CreateNewCaseResponse, error) {
	requestDef := GenReqDefForCreateNewCase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateNewCaseResponse), nil
	}
}

// CreateNewCaseInvoker 创建用例
func (c *CptsClient) CreateNewCaseInvoker(request *model.CreateNewCaseRequest) *CreateNewCaseInvoker {
	requestDef := GenReqDefForCreateNewCase()
	return &CreateNewCaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateNewTask 创建任务
//
// 创建任务
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) CreateNewTask(request *model.CreateNewTaskRequest) (*model.CreateNewTaskResponse, error) {
	requestDef := GenReqDefForCreateNewTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateNewTaskResponse), nil
	}
}

// CreateNewTaskInvoker 创建任务
func (c *CptsClient) CreateNewTaskInvoker(request *model.CreateNewTaskRequest) *CreateNewTaskInvoker {
	requestDef := GenReqDefForCreateNewTask()
	return &CreateNewTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateTask 创建任务（旧版）
//
// 创建任务（旧版）
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) CreateTask(request *model.CreateTaskRequest) (*model.CreateTaskResponse, error) {
	requestDef := GenReqDefForCreateTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateTaskResponse), nil
	}
}

// CreateTaskInvoker 创建任务（旧版）
func (c *CptsClient) CreateTaskInvoker(request *model.CreateTaskRequest) *CreateTaskInvoker {
	requestDef := GenReqDefForCreateTask()
	return &CreateTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateTemp 创建事务
//
// 创建事务
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) CreateTemp(request *model.CreateTempRequest) (*model.CreateTempResponse, error) {
	requestDef := GenReqDefForCreateTemp()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateTempResponse), nil
	}
}

// CreateTempInvoker 创建事务
func (c *CptsClient) CreateTempInvoker(request *model.CreateTempRequest) *CreateTempInvoker {
	requestDef := GenReqDefForCreateTemp()
	return &CreateTempInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateVariable 创建变量
//
// 创建变量
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) CreateVariable(request *model.CreateVariableRequest) (*model.CreateVariableResponse, error) {
	requestDef := GenReqDefForCreateVariable()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateVariableResponse), nil
	}
}

// CreateVariableInvoker 创建变量
func (c *CptsClient) CreateVariableInvoker(request *model.CreateVariableRequest) *CreateVariableInvoker {
	requestDef := GenReqDefForCreateVariable()
	return &CreateVariableInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DebugCase 调试用例
//
// 调试用例
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) DebugCase(request *model.DebugCaseRequest) (*model.DebugCaseResponse, error) {
	requestDef := GenReqDefForDebugCase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DebugCaseResponse), nil
	}
}

// DebugCaseInvoker 调试用例
func (c *CptsClient) DebugCaseInvoker(request *model.DebugCaseRequest) *DebugCaseInvoker {
	requestDef := GenReqDefForDebugCase()
	return &DebugCaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteCase 删除用例（旧版）
//
// 删除用例（旧版）
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) DeleteCase(request *model.DeleteCaseRequest) (*model.DeleteCaseResponse, error) {
	requestDef := GenReqDefForDeleteCase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteCaseResponse), nil
	}
}

// DeleteCaseInvoker 删除用例（旧版）
func (c *CptsClient) DeleteCaseInvoker(request *model.DeleteCaseRequest) *DeleteCaseInvoker {
	requestDef := GenReqDefForDeleteCase()
	return &DeleteCaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteDirectory 删除目录
//
// 删除目录
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) DeleteDirectory(request *model.DeleteDirectoryRequest) (*model.DeleteDirectoryResponse, error) {
	requestDef := GenReqDefForDeleteDirectory()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteDirectoryResponse), nil
	}
}

// DeleteDirectoryInvoker 删除目录
func (c *CptsClient) DeleteDirectoryInvoker(request *model.DeleteDirectoryRequest) *DeleteDirectoryInvoker {
	requestDef := GenReqDefForDeleteDirectory()
	return &DeleteDirectoryInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteNewCase 删除用例
//
// 删除用例
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) DeleteNewCase(request *model.DeleteNewCaseRequest) (*model.DeleteNewCaseResponse, error) {
	requestDef := GenReqDefForDeleteNewCase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteNewCaseResponse), nil
	}
}

// DeleteNewCaseInvoker 删除用例
func (c *CptsClient) DeleteNewCaseInvoker(request *model.DeleteNewCaseRequest) *DeleteNewCaseInvoker {
	requestDef := GenReqDefForDeleteNewCase()
	return &DeleteNewCaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteNewTask 删除任务
//
// 删除任务
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) DeleteNewTask(request *model.DeleteNewTaskRequest) (*model.DeleteNewTaskResponse, error) {
	requestDef := GenReqDefForDeleteNewTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteNewTaskResponse), nil
	}
}

// DeleteNewTaskInvoker 删除任务
func (c *CptsClient) DeleteNewTaskInvoker(request *model.DeleteNewTaskRequest) *DeleteNewTaskInvoker {
	requestDef := GenReqDefForDeleteNewTask()
	return &DeleteNewTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteTask 删除任务（旧版）
//
// 删除任务（旧版）
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) DeleteTask(request *model.DeleteTaskRequest) (*model.DeleteTaskResponse, error) {
	requestDef := GenReqDefForDeleteTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteTaskResponse), nil
	}
}

// DeleteTaskInvoker 删除任务（旧版）
func (c *CptsClient) DeleteTaskInvoker(request *model.DeleteTaskRequest) *DeleteTaskInvoker {
	requestDef := GenReqDefForDeleteTask()
	return &DeleteTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteTemp 删除事务
//
// 删除事务
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) DeleteTemp(request *model.DeleteTempRequest) (*model.DeleteTempResponse, error) {
	requestDef := GenReqDefForDeleteTemp()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteTempResponse), nil
	}
}

// DeleteTempInvoker 删除事务
func (c *CptsClient) DeleteTempInvoker(request *model.DeleteTempRequest) *DeleteTempInvoker {
	requestDef := GenReqDefForDeleteTemp()
	return &DeleteTempInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteVariable 删除全局变量
//
// 删除全局变量
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) DeleteVariable(request *model.DeleteVariableRequest) (*model.DeleteVariableResponse, error) {
	requestDef := GenReqDefForDeleteVariable()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteVariableResponse), nil
	}
}

// DeleteVariableInvoker 删除全局变量
func (c *CptsClient) DeleteVariableInvoker(request *model.DeleteVariableRequest) *DeleteVariableInvoker {
	requestDef := GenReqDefForDeleteVariable()
	return &DeleteVariableInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListProjectTestCase 查询用例树
//
// 查询用例树
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) ListProjectTestCase(request *model.ListProjectTestCaseRequest) (*model.ListProjectTestCaseResponse, error) {
	requestDef := GenReqDefForListProjectTestCase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListProjectTestCaseResponse), nil
	}
}

// ListProjectTestCaseInvoker 查询用例树
func (c *CptsClient) ListProjectTestCaseInvoker(request *model.ListProjectTestCaseRequest) *ListProjectTestCaseInvoker {
	requestDef := GenReqDefForListProjectTestCase()
	return &ListProjectTestCaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListTaskCases 获取任务关联的用例列表
//
// 获取任务关联的用例列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) ListTaskCases(request *model.ListTaskCasesRequest) (*model.ListTaskCasesResponse, error) {
	requestDef := GenReqDefForListTaskCases()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListTaskCasesResponse), nil
	}
}

// ListTaskCasesInvoker 获取任务关联的用例列表
func (c *CptsClient) ListTaskCasesInvoker(request *model.ListTaskCasesRequest) *ListTaskCasesInvoker {
	requestDef := GenReqDefForListTaskCases()
	return &ListTaskCasesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListVariables 查询全局变量
//
// 查询全局变量
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) ListVariables(request *model.ListVariablesRequest) (*model.ListVariablesResponse, error) {
	requestDef := GenReqDefForListVariables()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListVariablesResponse), nil
	}
}

// ListVariablesInvoker 查询全局变量
func (c *CptsClient) ListVariablesInvoker(request *model.ListVariablesRequest) *ListVariablesInvoker {
	requestDef := GenReqDefForListVariables()
	return &ListVariablesInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowAgentConfig 全链路压测探针获取配置信息
//
// 全链路压测探针获取配置信息
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) ShowAgentConfig(request *model.ShowAgentConfigRequest) (*model.ShowAgentConfigResponse, error) {
	requestDef := GenReqDefForShowAgentConfig()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowAgentConfigResponse), nil
	}
}

// ShowAgentConfigInvoker 全链路压测探针获取配置信息
func (c *CptsClient) ShowAgentConfigInvoker(request *model.ShowAgentConfigRequest) *ShowAgentConfigInvoker {
	requestDef := GenReqDefForShowAgentConfig()
	return &ShowAgentConfigInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowCase 查询用例
//
// 查询用例
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) ShowCase(request *model.ShowCaseRequest) (*model.ShowCaseResponse, error) {
	requestDef := GenReqDefForShowCase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowCaseResponse), nil
	}
}

// ShowCaseInvoker 查询用例
func (c *CptsClient) ShowCaseInvoker(request *model.ShowCaseRequest) *ShowCaseInvoker {
	requestDef := GenReqDefForShowCase()
	return &ShowCaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowHistoryRunInfo 查询PerfTest任务离线报告列表
//
// 查询PerfTest任务离线报告列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) ShowHistoryRunInfo(request *model.ShowHistoryRunInfoRequest) (*model.ShowHistoryRunInfoResponse, error) {
	requestDef := GenReqDefForShowHistoryRunInfo()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowHistoryRunInfoResponse), nil
	}
}

// ShowHistoryRunInfoInvoker 查询PerfTest任务离线报告列表
func (c *CptsClient) ShowHistoryRunInfoInvoker(request *model.ShowHistoryRunInfoRequest) *ShowHistoryRunInfoInvoker {
	requestDef := GenReqDefForShowHistoryRunInfo()
	return &ShowHistoryRunInfoInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowMergeCaseDetail 查询用例报告详情
//
// 查询用例报告详情
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) ShowMergeCaseDetail(request *model.ShowMergeCaseDetailRequest) (*model.ShowMergeCaseDetailResponse, error) {
	requestDef := GenReqDefForShowMergeCaseDetail()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowMergeCaseDetailResponse), nil
	}
}

// ShowMergeCaseDetailInvoker 查询用例报告详情
func (c *CptsClient) ShowMergeCaseDetailInvoker(request *model.ShowMergeCaseDetailRequest) *ShowMergeCaseDetailInvoker {
	requestDef := GenReqDefForShowMergeCaseDetail()
	return &ShowMergeCaseDetailInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowMergeReportLogsOutline 查询报告汇总数据
//
// 查询报告汇总数据
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) ShowMergeReportLogsOutline(request *model.ShowMergeReportLogsOutlineRequest) (*model.ShowMergeReportLogsOutlineResponse, error) {
	requestDef := GenReqDefForShowMergeReportLogsOutline()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowMergeReportLogsOutlineResponse), nil
	}
}

// ShowMergeReportLogsOutlineInvoker 查询报告汇总数据
func (c *CptsClient) ShowMergeReportLogsOutlineInvoker(request *model.ShowMergeReportLogsOutlineRequest) *ShowMergeReportLogsOutlineInvoker {
	requestDef := GenReqDefForShowMergeReportLogsOutline()
	return &ShowMergeReportLogsOutlineInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowMergeTaskCase 查询任务报告的用例列表
//
// 查询任务报告的用例列表
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) ShowMergeTaskCase(request *model.ShowMergeTaskCaseRequest) (*model.ShowMergeTaskCaseResponse, error) {
	requestDef := GenReqDefForShowMergeTaskCase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowMergeTaskCaseResponse), nil
	}
}

// ShowMergeTaskCaseInvoker 查询任务报告的用例列表
func (c *CptsClient) ShowMergeTaskCaseInvoker(request *model.ShowMergeTaskCaseRequest) *ShowMergeTaskCaseInvoker {
	requestDef := GenReqDefForShowMergeTaskCase()
	return &ShowMergeTaskCaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowReport 查询报告
//
// 查询报告
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) ShowReport(request *model.ShowReportRequest) (*model.ShowReportResponse, error) {
	requestDef := GenReqDefForShowReport()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowReportResponse), nil
	}
}

// ShowReportInvoker 查询报告
func (c *CptsClient) ShowReportInvoker(request *model.ShowReportRequest) *ShowReportInvoker {
	requestDef := GenReqDefForShowReport()
	return &ShowReportInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowTask 查询任务（旧版）
//
// 查询任务（旧版）
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) ShowTask(request *model.ShowTaskRequest) (*model.ShowTaskResponse, error) {
	requestDef := GenReqDefForShowTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowTaskResponse), nil
	}
}

// ShowTaskInvoker 查询任务（旧版）
func (c *CptsClient) ShowTaskInvoker(request *model.ShowTaskRequest) *ShowTaskInvoker {
	requestDef := GenReqDefForShowTask()
	return &ShowTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowTaskCaseAwChart 查询用例的AW曲线图
//
// 查询用例的AW曲线图
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) ShowTaskCaseAwChart(request *model.ShowTaskCaseAwChartRequest) (*model.ShowTaskCaseAwChartResponse, error) {
	requestDef := GenReqDefForShowTaskCaseAwChart()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowTaskCaseAwChartResponse), nil
	}
}

// ShowTaskCaseAwChartInvoker 查询用例的AW曲线图
func (c *CptsClient) ShowTaskCaseAwChartInvoker(request *model.ShowTaskCaseAwChartRequest) *ShowTaskCaseAwChartInvoker {
	requestDef := GenReqDefForShowTaskCaseAwChart()
	return &ShowTaskCaseAwChartInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowTaskSet 查询任务集
//
// 查询任务集
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) ShowTaskSet(request *model.ShowTaskSetRequest) (*model.ShowTaskSetResponse, error) {
	requestDef := GenReqDefForShowTaskSet()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowTaskSetResponse), nil
	}
}

// ShowTaskSetInvoker 查询任务集
func (c *CptsClient) ShowTaskSetInvoker(request *model.ShowTaskSetRequest) *ShowTaskSetInvoker {
	requestDef := GenReqDefForShowTaskSet()
	return &ShowTaskSetInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowTemp 查询事务
//
// 查询事务
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) ShowTemp(request *model.ShowTempRequest) (*model.ShowTempResponse, error) {
	requestDef := GenReqDefForShowTemp()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowTempResponse), nil
	}
}

// ShowTempInvoker 查询事务
func (c *CptsClient) ShowTempInvoker(request *model.ShowTempRequest) *ShowTempInvoker {
	requestDef := GenReqDefForShowTemp()
	return &ShowTempInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowTempSet 查询事务集
//
// 查询事务集
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) ShowTempSet(request *model.ShowTempSetRequest) (*model.ShowTempSetResponse, error) {
	requestDef := GenReqDefForShowTempSet()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowTempSetResponse), nil
	}
}

// ShowTempSetInvoker 查询事务集
func (c *CptsClient) ShowTempSetInvoker(request *model.ShowTempSetRequest) *ShowTempSetInvoker {
	requestDef := GenReqDefForShowTempSet()
	return &ShowTempSetInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateAgentHealthStatus 全链路压测探针上报健康状态
//
// 全链路压测探针上报健康状态
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) UpdateAgentHealthStatus(request *model.UpdateAgentHealthStatusRequest) (*model.UpdateAgentHealthStatusResponse, error) {
	requestDef := GenReqDefForUpdateAgentHealthStatus()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateAgentHealthStatusResponse), nil
	}
}

// UpdateAgentHealthStatusInvoker 全链路压测探针上报健康状态
func (c *CptsClient) UpdateAgentHealthStatusInvoker(request *model.UpdateAgentHealthStatusRequest) *UpdateAgentHealthStatusInvoker {
	requestDef := GenReqDefForUpdateAgentHealthStatus()
	return &UpdateAgentHealthStatusInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateCase 修改用例（旧版）
//
// 修改用例（旧版）
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) UpdateCase(request *model.UpdateCaseRequest) (*model.UpdateCaseResponse, error) {
	requestDef := GenReqDefForUpdateCase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateCaseResponse), nil
	}
}

// UpdateCaseInvoker 修改用例（旧版）
func (c *CptsClient) UpdateCaseInvoker(request *model.UpdateCaseRequest) *UpdateCaseInvoker {
	requestDef := GenReqDefForUpdateCase()
	return &UpdateCaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateDirectory 修改目录
//
// 修改目录
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) UpdateDirectory(request *model.UpdateDirectoryRequest) (*model.UpdateDirectoryResponse, error) {
	requestDef := GenReqDefForUpdateDirectory()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateDirectoryResponse), nil
	}
}

// UpdateDirectoryInvoker 修改目录
func (c *CptsClient) UpdateDirectoryInvoker(request *model.UpdateDirectoryRequest) *UpdateDirectoryInvoker {
	requestDef := GenReqDefForUpdateDirectory()
	return &UpdateDirectoryInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateNewCase 修改用例
//
// 修改用例
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) UpdateNewCase(request *model.UpdateNewCaseRequest) (*model.UpdateNewCaseResponse, error) {
	requestDef := GenReqDefForUpdateNewCase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateNewCaseResponse), nil
	}
}

// UpdateNewCaseInvoker 修改用例
func (c *CptsClient) UpdateNewCaseInvoker(request *model.UpdateNewCaseRequest) *UpdateNewCaseInvoker {
	requestDef := GenReqDefForUpdateNewCase()
	return &UpdateNewCaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateTask 修改任务（旧版）
//
// 修改任务（旧版）
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) UpdateTask(request *model.UpdateTaskRequest) (*model.UpdateTaskResponse, error) {
	requestDef := GenReqDefForUpdateTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateTaskResponse), nil
	}
}

// UpdateTaskInvoker 修改任务（旧版）
func (c *CptsClient) UpdateTaskInvoker(request *model.UpdateTaskRequest) *UpdateTaskInvoker {
	requestDef := GenReqDefForUpdateTask()
	return &UpdateTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateTaskRelatedTestCase 修改任务关联用例
//
// 修改任务关联用例
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) UpdateTaskRelatedTestCase(request *model.UpdateTaskRelatedTestCaseRequest) (*model.UpdateTaskRelatedTestCaseResponse, error) {
	requestDef := GenReqDefForUpdateTaskRelatedTestCase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateTaskRelatedTestCaseResponse), nil
	}
}

// UpdateTaskRelatedTestCaseInvoker 修改任务关联用例
func (c *CptsClient) UpdateTaskRelatedTestCaseInvoker(request *model.UpdateTaskRelatedTestCaseRequest) *UpdateTaskRelatedTestCaseInvoker {
	requestDef := GenReqDefForUpdateTaskRelatedTestCase()
	return &UpdateTaskRelatedTestCaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateTaskStatus 更新任务状态
//
// 更新任务状态
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) UpdateTaskStatus(request *model.UpdateTaskStatusRequest) (*model.UpdateTaskStatusResponse, error) {
	requestDef := GenReqDefForUpdateTaskStatus()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateTaskStatusResponse), nil
	}
}

// UpdateTaskStatusInvoker 更新任务状态
func (c *CptsClient) UpdateTaskStatusInvoker(request *model.UpdateTaskStatusRequest) *UpdateTaskStatusInvoker {
	requestDef := GenReqDefForUpdateTaskStatus()
	return &UpdateTaskStatusInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateTemp 修改事务
//
// 修改事务
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) UpdateTemp(request *model.UpdateTempRequest) (*model.UpdateTempResponse, error) {
	requestDef := GenReqDefForUpdateTemp()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateTempResponse), nil
	}
}

// UpdateTempInvoker 修改事务
func (c *CptsClient) UpdateTempInvoker(request *model.UpdateTempRequest) *UpdateTempInvoker {
	requestDef := GenReqDefForUpdateTemp()
	return &UpdateTempInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateVariable 修改变量
//
// 修改变量
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) UpdateVariable(request *model.UpdateVariableRequest) (*model.UpdateVariableResponse, error) {
	requestDef := GenReqDefForUpdateVariable()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateVariableResponse), nil
	}
}

// UpdateVariableInvoker 修改变量
func (c *CptsClient) UpdateVariableInvoker(request *model.UpdateVariableRequest) *UpdateVariableInvoker {
	requestDef := GenReqDefForUpdateVariable()
	return &UpdateVariableInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateProject 创建工程
//
// 创建工程
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) CreateProject(request *model.CreateProjectRequest) (*model.CreateProjectResponse, error) {
	requestDef := GenReqDefForCreateProject()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateProjectResponse), nil
	}
}

// CreateProjectInvoker 创建工程
func (c *CptsClient) CreateProjectInvoker(request *model.CreateProjectRequest) *CreateProjectInvoker {
	requestDef := GenReqDefForCreateProject()
	return &CreateProjectInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteProject 删除工程
//
// 删除工程
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) DeleteProject(request *model.DeleteProjectRequest) (*model.DeleteProjectResponse, error) {
	requestDef := GenReqDefForDeleteProject()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteProjectResponse), nil
	}
}

// DeleteProjectInvoker 删除工程
func (c *CptsClient) DeleteProjectInvoker(request *model.DeleteProjectRequest) *DeleteProjectInvoker {
	requestDef := GenReqDefForDeleteProject()
	return &DeleteProjectInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ListProjectSets 查询工程集
//
// 查询工程集
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) ListProjectSets(request *model.ListProjectSetsRequest) (*model.ListProjectSetsResponse, error) {
	requestDef := GenReqDefForListProjectSets()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListProjectSetsResponse), nil
	}
}

// ListProjectSetsInvoker 查询工程集
func (c *CptsClient) ListProjectSetsInvoker(request *model.ListProjectSetsRequest) *ListProjectSetsInvoker {
	requestDef := GenReqDefForListProjectSets()
	return &ListProjectSetsInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowProcess 查询导入进度
//
// 查询导入进度
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) ShowProcess(request *model.ShowProcessRequest) (*model.ShowProcessResponse, error) {
	requestDef := GenReqDefForShowProcess()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowProcessResponse), nil
	}
}

// ShowProcessInvoker 查询导入进度
func (c *CptsClient) ShowProcessInvoker(request *model.ShowProcessRequest) *ShowProcessInvoker {
	requestDef := GenReqDefForShowProcess()
	return &ShowProcessInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowProject 查询工程
//
// 查询工程
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) ShowProject(request *model.ShowProjectRequest) (*model.ShowProjectResponse, error) {
	requestDef := GenReqDefForShowProject()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowProjectResponse), nil
	}
}

// ShowProjectInvoker 查询工程
func (c *CptsClient) ShowProjectInvoker(request *model.ShowProjectRequest) *ShowProjectInvoker {
	requestDef := GenReqDefForShowProject()
	return &ShowProjectInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateProject 修改工程
//
// 修改工程
//
// Please refer to HUAWEI cloud API Explorer for details.
func (c *CptsClient) UpdateProject(request *model.UpdateProjectRequest) (*model.UpdateProjectResponse, error) {
	requestDef := GenReqDefForUpdateProject()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateProjectResponse), nil
	}
}

// UpdateProjectInvoker 修改工程
func (c *CptsClient) UpdateProjectInvoker(request *model.UpdateProjectRequest) *UpdateProjectInvoker {
	requestDef := GenReqDefForUpdateProject()
	return &UpdateProjectInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}
