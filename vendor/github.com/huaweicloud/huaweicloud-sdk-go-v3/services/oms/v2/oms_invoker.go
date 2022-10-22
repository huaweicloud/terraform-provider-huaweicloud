package v2

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/oms/v2/model"
)

type CreateSyncEventsInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateSyncEventsInvoker) Invoke() (*model.CreateSyncEventsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateSyncEventsResponse), nil
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

type ListTasksInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListTasksInvoker) Invoke() (*model.ListTasksResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListTasksResponse), nil
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

type StartTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *StartTaskInvoker) Invoke() (*model.StartTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StartTaskResponse), nil
	}
}

type StopTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *StopTaskInvoker) Invoke() (*model.StopTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StopTaskResponse), nil
	}
}

type UpdateBandwidthPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateBandwidthPolicyInvoker) Invoke() (*model.UpdateBandwidthPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateBandwidthPolicyResponse), nil
	}
}

type CreateTaskGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateTaskGroupInvoker) Invoke() (*model.CreateTaskGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateTaskGroupResponse), nil
	}
}

type DeleteTaskGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteTaskGroupInvoker) Invoke() (*model.DeleteTaskGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteTaskGroupResponse), nil
	}
}

type ListTaskGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListTaskGroupInvoker) Invoke() (*model.ListTaskGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListTaskGroupResponse), nil
	}
}

type RetryTaskGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *RetryTaskGroupInvoker) Invoke() (*model.RetryTaskGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RetryTaskGroupResponse), nil
	}
}

type ShowTaskGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowTaskGroupInvoker) Invoke() (*model.ShowTaskGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowTaskGroupResponse), nil
	}
}

type StartTaskGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *StartTaskGroupInvoker) Invoke() (*model.StartTaskGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StartTaskGroupResponse), nil
	}
}

type StopTaskGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *StopTaskGroupInvoker) Invoke() (*model.StopTaskGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StopTaskGroupResponse), nil
	}
}

type UpdateTaskGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateTaskGroupInvoker) Invoke() (*model.UpdateTaskGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateTaskGroupResponse), nil
	}
}

type ListApiVersionsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListApiVersionsInvoker) Invoke() (*model.ListApiVersionsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListApiVersionsResponse), nil
	}
}

type ShowApiInfoInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowApiInfoInvoker) Invoke() (*model.ShowApiInfoResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowApiInfoResponse), nil
	}
}
