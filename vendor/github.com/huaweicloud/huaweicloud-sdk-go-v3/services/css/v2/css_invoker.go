package v2

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/css/v2/model"
)

type CreateClusterInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateClusterInvoker) Invoke() (*model.CreateClusterResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateClusterResponse), nil
	}
}

type RestartClusterInvoker struct {
	*invoker.BaseInvoker
}

func (i *RestartClusterInvoker) Invoke() (*model.RestartClusterResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RestartClusterResponse), nil
	}
}

type RollingRestartInvoker struct {
	*invoker.BaseInvoker
}

func (i *RollingRestartInvoker) Invoke() (*model.RollingRestartResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RollingRestartResponse), nil
	}
}

type StartAutoCreateSnapshotsInvoker struct {
	*invoker.BaseInvoker
}

func (i *StartAutoCreateSnapshotsInvoker) Invoke() (*model.StartAutoCreateSnapshotsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StartAutoCreateSnapshotsResponse), nil
	}
}

type StopAutoCreateSnapshotsInvoker struct {
	*invoker.BaseInvoker
}

func (i *StopAutoCreateSnapshotsInvoker) Invoke() (*model.StopAutoCreateSnapshotsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StopAutoCreateSnapshotsResponse), nil
	}
}
