package v1

import (
	http_client "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cpts/v1/model"
)

type CptsClient struct {
	HcClient *http_client.HcHttpClient
}

func NewCptsClient(hcClient *http_client.HcHttpClient) *CptsClient {
	return &CptsClient{HcClient: hcClient}
}

func CptsClientBuilder() *http_client.HcHttpClientBuilder {
	builder := http_client.NewHcHttpClientBuilder()
	return builder
}

// CreateCase 创建用例
//
// 创建用例
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *CptsClient) CreateCase(request *model.CreateCaseRequest) (*model.CreateCaseResponse, error) {
	requestDef := GenReqDefForCreateCase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateCaseResponse), nil
	}
}

// CreateCaseInvoker 创建用例
func (c *CptsClient) CreateCaseInvoker(request *model.CreateCaseRequest) *CreateCaseInvoker {
	requestDef := GenReqDefForCreateCase()
	return &CreateCaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateTask 创建任务
//
// 创建任务
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *CptsClient) CreateTask(request *model.CreateTaskRequest) (*model.CreateTaskResponse, error) {
	requestDef := GenReqDefForCreateTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateTaskResponse), nil
	}
}

// CreateTaskInvoker 创建任务
func (c *CptsClient) CreateTaskInvoker(request *model.CreateTaskRequest) *CreateTaskInvoker {
	requestDef := GenReqDefForCreateTask()
	return &CreateTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// CreateTemp 创建事务
//
// 创建事务
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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

// DeleteCase 删除用例
//
// 删除用例
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *CptsClient) DeleteCase(request *model.DeleteCaseRequest) (*model.DeleteCaseResponse, error) {
	requestDef := GenReqDefForDeleteCase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteCaseResponse), nil
	}
}

// DeleteCaseInvoker 删除用例
func (c *CptsClient) DeleteCaseInvoker(request *model.DeleteCaseRequest) *DeleteCaseInvoker {
	requestDef := GenReqDefForDeleteCase()
	return &DeleteCaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteTask 删除任务
//
// 删除任务
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *CptsClient) DeleteTask(request *model.DeleteTaskRequest) (*model.DeleteTaskResponse, error) {
	requestDef := GenReqDefForDeleteTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteTaskResponse), nil
	}
}

// DeleteTaskInvoker 删除任务
func (c *CptsClient) DeleteTaskInvoker(request *model.DeleteTaskRequest) *DeleteTaskInvoker {
	requestDef := GenReqDefForDeleteTask()
	return &DeleteTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// DeleteTemp 删除事务
//
// 删除事务
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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

// ListVariables 查询全局变量
//
// 查询全局变量
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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

// ShowHistoryRunInfo 查询CPTS任务离线报告列表
//
// 查询CPTS任务离线报告列表
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *CptsClient) ShowHistoryRunInfo(request *model.ShowHistoryRunInfoRequest) (*model.ShowHistoryRunInfoResponse, error) {
	requestDef := GenReqDefForShowHistoryRunInfo()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowHistoryRunInfoResponse), nil
	}
}

// ShowHistoryRunInfoInvoker 查询CPTS任务离线报告列表
func (c *CptsClient) ShowHistoryRunInfoInvoker(request *model.ShowHistoryRunInfoRequest) *ShowHistoryRunInfoInvoker {
	requestDef := GenReqDefForShowHistoryRunInfo()
	return &ShowHistoryRunInfoInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowReport 查询报告
//
// 查询报告
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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

// ShowTask 查询任务
//
// 查询任务
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *CptsClient) ShowTask(request *model.ShowTaskRequest) (*model.ShowTaskResponse, error) {
	requestDef := GenReqDefForShowTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowTaskResponse), nil
	}
}

// ShowTaskInvoker 查询任务
func (c *CptsClient) ShowTaskInvoker(request *model.ShowTaskRequest) *ShowTaskInvoker {
	requestDef := GenReqDefForShowTask()
	return &ShowTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// ShowTaskSet 查询任务集
//
// 查询任务集
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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

// UpdateCase 修改用例
//
// 修改用例
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *CptsClient) UpdateCase(request *model.UpdateCaseRequest) (*model.UpdateCaseResponse, error) {
	requestDef := GenReqDefForUpdateCase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateCaseResponse), nil
	}
}

// UpdateCaseInvoker 修改用例
func (c *CptsClient) UpdateCaseInvoker(request *model.UpdateCaseRequest) *UpdateCaseInvoker {
	requestDef := GenReqDefForUpdateCase()
	return &UpdateCaseInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateTask 修改任务
//
// 修改任务
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
func (c *CptsClient) UpdateTask(request *model.UpdateTaskRequest) (*model.UpdateTaskResponse, error) {
	requestDef := GenReqDefForUpdateTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateTaskResponse), nil
	}
}

// UpdateTaskInvoker 修改任务
func (c *CptsClient) UpdateTaskInvoker(request *model.UpdateTaskRequest) *UpdateTaskInvoker {
	requestDef := GenReqDefForUpdateTask()
	return &UpdateTaskInvoker{invoker.NewBaseInvoker(c.HcClient, request, requestDef)}
}

// UpdateTaskStatus 更新任务状态
//
// 更新任务状态
//
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
// 详细说明请参考华为云API Explorer。
// Please refer to Huawei cloud API Explorer for details.
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
