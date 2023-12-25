package v2

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/oms/v2/model"
)

type BatchUpdateTasksInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchUpdateTasksInvoker) Invoke() (*model.BatchUpdateTasksResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchUpdateTasksResponse), nil
	}
}

type CheckPrefixInvoker struct {
	*invoker.BaseInvoker
}

func (i *CheckPrefixInvoker) Invoke() (*model.CheckPrefixResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CheckPrefixResponse), nil
	}
}

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

type CreateSyncTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateSyncTaskInvoker) Invoke() (*model.CreateSyncTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateSyncTaskResponse), nil
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

type DeleteSyncTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteSyncTaskInvoker) Invoke() (*model.DeleteSyncTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteSyncTaskResponse), nil
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

type ListSyncTaskStatisticInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListSyncTaskStatisticInvoker) Invoke() (*model.ListSyncTaskStatisticResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListSyncTaskStatisticResponse), nil
	}
}

type ListSyncTasksInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListSyncTasksInvoker) Invoke() (*model.ListSyncTasksResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListSyncTasksResponse), nil
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

type ShowBucketListInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowBucketListInvoker) Invoke() (*model.ShowBucketListResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowBucketListResponse), nil
	}
}

type ShowBucketObjectsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowBucketObjectsInvoker) Invoke() (*model.ShowBucketObjectsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowBucketObjectsResponse), nil
	}
}

type ShowBucketRegionInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowBucketRegionInvoker) Invoke() (*model.ShowBucketRegionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowBucketRegionResponse), nil
	}
}

type ShowCdnInfoInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowCdnInfoInvoker) Invoke() (*model.ShowCdnInfoResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowCdnInfoResponse), nil
	}
}

type ShowCloudTypeInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowCloudTypeInvoker) Invoke() (*model.ShowCloudTypeResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowCloudTypeResponse), nil
	}
}

type ShowRegionInfoInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowRegionInfoInvoker) Invoke() (*model.ShowRegionInfoResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowRegionInfoResponse), nil
	}
}

type ShowSyncTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowSyncTaskInvoker) Invoke() (*model.ShowSyncTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowSyncTaskResponse), nil
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

type StartSyncTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *StartSyncTaskInvoker) Invoke() (*model.StartSyncTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StartSyncTaskResponse), nil
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

type StopSyncTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *StopSyncTaskInvoker) Invoke() (*model.StopSyncTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StopSyncTaskResponse), nil
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
