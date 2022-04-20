package v1

import (
	http_client "github.com/huaweicloud/huaweicloud-sdk-go-v3/core"

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

//创建用例
func (c *CptsClient) CreateCase(request *model.CreateCaseRequest) (*model.CreateCaseResponse, error) {
	requestDef := GenReqDefForCreateCase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateCaseResponse), nil
	}
}

//创建任务
func (c *CptsClient) CreateTask(request *model.CreateTaskRequest) (*model.CreateTaskResponse, error) {
	requestDef := GenReqDefForCreateTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateTaskResponse), nil
	}
}

//创建事务
func (c *CptsClient) CreateTemp(request *model.CreateTempRequest) (*model.CreateTempResponse, error) {
	requestDef := GenReqDefForCreateTemp()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateTempResponse), nil
	}
}

//创建变量
func (c *CptsClient) CreateVariable(request *model.CreateVariableRequest) (*model.CreateVariableResponse, error) {
	requestDef := GenReqDefForCreateVariable()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateVariableResponse), nil
	}
}

//调试用例
func (c *CptsClient) DebugCase(request *model.DebugCaseRequest) (*model.DebugCaseResponse, error) {
	requestDef := GenReqDefForDebugCase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DebugCaseResponse), nil
	}
}

//删除用例
func (c *CptsClient) DeleteCase(request *model.DeleteCaseRequest) (*model.DeleteCaseResponse, error) {
	requestDef := GenReqDefForDeleteCase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteCaseResponse), nil
	}
}

//删除任务
func (c *CptsClient) DeleteTask(request *model.DeleteTaskRequest) (*model.DeleteTaskResponse, error) {
	requestDef := GenReqDefForDeleteTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteTaskResponse), nil
	}
}

//删除事务
func (c *CptsClient) DeleteTemp(request *model.DeleteTempRequest) (*model.DeleteTempResponse, error) {
	requestDef := GenReqDefForDeleteTemp()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteTempResponse), nil
	}
}

//查询全局变量
func (c *CptsClient) ListVariables(request *model.ListVariablesRequest) (*model.ListVariablesResponse, error) {
	requestDef := GenReqDefForListVariables()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListVariablesResponse), nil
	}
}

//查询CPTS任务离线报告列表
func (c *CptsClient) ShowHistoryRunInfo(request *model.ShowHistoryRunInfoRequest) (*model.ShowHistoryRunInfoResponse, error) {
	requestDef := GenReqDefForShowHistoryRunInfo()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowHistoryRunInfoResponse), nil
	}
}

//查询报告
func (c *CptsClient) ShowReport(request *model.ShowReportRequest) (*model.ShowReportResponse, error) {
	requestDef := GenReqDefForShowReport()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowReportResponse), nil
	}
}

//查询任务
func (c *CptsClient) ShowTask(request *model.ShowTaskRequest) (*model.ShowTaskResponse, error) {
	requestDef := GenReqDefForShowTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowTaskResponse), nil
	}
}

//查询任务集
func (c *CptsClient) ShowTaskSet(request *model.ShowTaskSetRequest) (*model.ShowTaskSetResponse, error) {
	requestDef := GenReqDefForShowTaskSet()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowTaskSetResponse), nil
	}
}

//查询事务
func (c *CptsClient) ShowTemp(request *model.ShowTempRequest) (*model.ShowTempResponse, error) {
	requestDef := GenReqDefForShowTemp()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowTempResponse), nil
	}
}

//查询事务集
func (c *CptsClient) ShowTempSet(request *model.ShowTempSetRequest) (*model.ShowTempSetResponse, error) {
	requestDef := GenReqDefForShowTempSet()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowTempSetResponse), nil
	}
}

//修改用例
func (c *CptsClient) UpdateCase(request *model.UpdateCaseRequest) (*model.UpdateCaseResponse, error) {
	requestDef := GenReqDefForUpdateCase()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateCaseResponse), nil
	}
}

//修改任务
func (c *CptsClient) UpdateTask(request *model.UpdateTaskRequest) (*model.UpdateTaskResponse, error) {
	requestDef := GenReqDefForUpdateTask()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateTaskResponse), nil
	}
}

//更新任务状态
func (c *CptsClient) UpdateTaskStatus(request *model.UpdateTaskStatusRequest) (*model.UpdateTaskStatusResponse, error) {
	requestDef := GenReqDefForUpdateTaskStatus()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateTaskStatusResponse), nil
	}
}

//修改事务
func (c *CptsClient) UpdateTemp(request *model.UpdateTempRequest) (*model.UpdateTempResponse, error) {
	requestDef := GenReqDefForUpdateTemp()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateTempResponse), nil
	}
}

//修改变量
func (c *CptsClient) UpdateVariable(request *model.UpdateVariableRequest) (*model.UpdateVariableResponse, error) {
	requestDef := GenReqDefForUpdateVariable()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateVariableResponse), nil
	}
}

//创建工程
func (c *CptsClient) CreateProject(request *model.CreateProjectRequest) (*model.CreateProjectResponse, error) {
	requestDef := GenReqDefForCreateProject()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.CreateProjectResponse), nil
	}
}

//删除工程
func (c *CptsClient) DeleteProject(request *model.DeleteProjectRequest) (*model.DeleteProjectResponse, error) {
	requestDef := GenReqDefForDeleteProject()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.DeleteProjectResponse), nil
	}
}

//查询工程集
func (c *CptsClient) ListProjectSets(request *model.ListProjectSetsRequest) (*model.ListProjectSetsResponse, error) {
	requestDef := GenReqDefForListProjectSets()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ListProjectSetsResponse), nil
	}
}

//查询导入进度
func (c *CptsClient) ShowProcess(request *model.ShowProcessRequest) (*model.ShowProcessResponse, error) {
	requestDef := GenReqDefForShowProcess()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowProcessResponse), nil
	}
}

//查询工程
func (c *CptsClient) ShowProject(request *model.ShowProjectRequest) (*model.ShowProjectResponse, error) {
	requestDef := GenReqDefForShowProject()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.ShowProjectResponse), nil
	}
}

//修改工程
func (c *CptsClient) UpdateProject(request *model.UpdateProjectRequest) (*model.UpdateProjectResponse, error) {
	requestDef := GenReqDefForUpdateProject()

	if resp, err := c.HcClient.Sync(request, requestDef); err != nil {
		return nil, err
	} else {
		return resp.(*model.UpdateProjectResponse), nil
	}
}
