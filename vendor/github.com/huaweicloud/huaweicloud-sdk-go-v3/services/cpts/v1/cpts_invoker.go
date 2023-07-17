package v1

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cpts/v1/model"
)

type BatchUpdateTaskStatusInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchUpdateTaskStatusInvoker) Invoke() (*model.BatchUpdateTaskStatusResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchUpdateTaskStatusResponse), nil
	}
}

type CreateCaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateCaseInvoker) Invoke() (*model.CreateCaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateCaseResponse), nil
	}
}

type CreateDirectoryInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateDirectoryInvoker) Invoke() (*model.CreateDirectoryResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateDirectoryResponse), nil
	}
}

type CreateNewCaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateNewCaseInvoker) Invoke() (*model.CreateNewCaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateNewCaseResponse), nil
	}
}

type CreateNewTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateNewTaskInvoker) Invoke() (*model.CreateNewTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateNewTaskResponse), nil
	}
}

type CreateTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateTaskInvoker) Invoke() (*model.CreateTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateTaskResponse), nil
	}
}

type CreateTempInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateTempInvoker) Invoke() (*model.CreateTempResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateTempResponse), nil
	}
}

type CreateVariableInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateVariableInvoker) Invoke() (*model.CreateVariableResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateVariableResponse), nil
	}
}

type DebugCaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *DebugCaseInvoker) Invoke() (*model.DebugCaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DebugCaseResponse), nil
	}
}

type DeleteCaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteCaseInvoker) Invoke() (*model.DeleteCaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteCaseResponse), nil
	}
}

type DeleteDirectoryInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteDirectoryInvoker) Invoke() (*model.DeleteDirectoryResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteDirectoryResponse), nil
	}
}

type DeleteNewCaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteNewCaseInvoker) Invoke() (*model.DeleteNewCaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteNewCaseResponse), nil
	}
}

type DeleteNewTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteNewTaskInvoker) Invoke() (*model.DeleteNewTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteNewTaskResponse), nil
	}
}

type DeleteTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteTaskInvoker) Invoke() (*model.DeleteTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteTaskResponse), nil
	}
}

type DeleteTempInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteTempInvoker) Invoke() (*model.DeleteTempResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteTempResponse), nil
	}
}

type DeleteVariableInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteVariableInvoker) Invoke() (*model.DeleteVariableResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteVariableResponse), nil
	}
}

type ListProjectTestCaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListProjectTestCaseInvoker) Invoke() (*model.ListProjectTestCaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListProjectTestCaseResponse), nil
	}
}

type ListTaskCasesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListTaskCasesInvoker) Invoke() (*model.ListTaskCasesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListTaskCasesResponse), nil
	}
}

type ListVariablesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListVariablesInvoker) Invoke() (*model.ListVariablesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListVariablesResponse), nil
	}
}

type ShowAgentConfigInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAgentConfigInvoker) Invoke() (*model.ShowAgentConfigResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAgentConfigResponse), nil
	}
}

type ShowCaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowCaseInvoker) Invoke() (*model.ShowCaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowCaseResponse), nil
	}
}

type ShowHistoryRunInfoInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowHistoryRunInfoInvoker) Invoke() (*model.ShowHistoryRunInfoResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowHistoryRunInfoResponse), nil
	}
}

type ShowMergeCaseDetailInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowMergeCaseDetailInvoker) Invoke() (*model.ShowMergeCaseDetailResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowMergeCaseDetailResponse), nil
	}
}

type ShowMergeReportLogsOutlineInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowMergeReportLogsOutlineInvoker) Invoke() (*model.ShowMergeReportLogsOutlineResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowMergeReportLogsOutlineResponse), nil
	}
}

type ShowMergeTaskCaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowMergeTaskCaseInvoker) Invoke() (*model.ShowMergeTaskCaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowMergeTaskCaseResponse), nil
	}
}

type ShowReportInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowReportInvoker) Invoke() (*model.ShowReportResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowReportResponse), nil
	}
}

type ShowTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowTaskInvoker) Invoke() (*model.ShowTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowTaskResponse), nil
	}
}

type ShowTaskCaseAwChartInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowTaskCaseAwChartInvoker) Invoke() (*model.ShowTaskCaseAwChartResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowTaskCaseAwChartResponse), nil
	}
}

type ShowTaskSetInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowTaskSetInvoker) Invoke() (*model.ShowTaskSetResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowTaskSetResponse), nil
	}
}

type ShowTempInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowTempInvoker) Invoke() (*model.ShowTempResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowTempResponse), nil
	}
}

type ShowTempSetInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowTempSetInvoker) Invoke() (*model.ShowTempSetResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowTempSetResponse), nil
	}
}

type UpdateAgentHealthStatusInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateAgentHealthStatusInvoker) Invoke() (*model.UpdateAgentHealthStatusResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateAgentHealthStatusResponse), nil
	}
}

type UpdateCaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateCaseInvoker) Invoke() (*model.UpdateCaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateCaseResponse), nil
	}
}

type UpdateDirectoryInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateDirectoryInvoker) Invoke() (*model.UpdateDirectoryResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateDirectoryResponse), nil
	}
}

type UpdateNewCaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateNewCaseInvoker) Invoke() (*model.UpdateNewCaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateNewCaseResponse), nil
	}
}

type UpdateTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateTaskInvoker) Invoke() (*model.UpdateTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateTaskResponse), nil
	}
}

type UpdateTaskRelatedTestCaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateTaskRelatedTestCaseInvoker) Invoke() (*model.UpdateTaskRelatedTestCaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateTaskRelatedTestCaseResponse), nil
	}
}

type UpdateTaskStatusInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateTaskStatusInvoker) Invoke() (*model.UpdateTaskStatusResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateTaskStatusResponse), nil
	}
}

type UpdateTempInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateTempInvoker) Invoke() (*model.UpdateTempResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateTempResponse), nil
	}
}

type UpdateVariableInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateVariableInvoker) Invoke() (*model.UpdateVariableResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateVariableResponse), nil
	}
}

type CreateProjectInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateProjectInvoker) Invoke() (*model.CreateProjectResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateProjectResponse), nil
	}
}

type DeleteProjectInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteProjectInvoker) Invoke() (*model.DeleteProjectResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteProjectResponse), nil
	}
}

type ListProjectSetsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListProjectSetsInvoker) Invoke() (*model.ListProjectSetsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListProjectSetsResponse), nil
	}
}

type ShowProcessInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowProcessInvoker) Invoke() (*model.ShowProcessResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowProcessResponse), nil
	}
}

type ShowProjectInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowProjectInvoker) Invoke() (*model.ShowProjectResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowProjectResponse), nil
	}
}

type UpdateProjectInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateProjectInvoker) Invoke() (*model.UpdateProjectResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateProjectResponse), nil
	}
}
