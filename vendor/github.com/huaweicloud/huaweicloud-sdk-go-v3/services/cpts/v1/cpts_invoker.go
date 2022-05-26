package v1

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cpts/v1/model"
)

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
