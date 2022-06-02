package v3

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/kps/v3/model"
)

type AssociateKeypairInvoker struct {
	*invoker.BaseInvoker
}

func (i *AssociateKeypairInvoker) Invoke() (*model.AssociateKeypairResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AssociateKeypairResponse), nil
	}
}

type CreateKeypairInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateKeypairInvoker) Invoke() (*model.CreateKeypairResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateKeypairResponse), nil
	}
}

type DeleteAllFailedTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteAllFailedTaskInvoker) Invoke() (*model.DeleteAllFailedTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteAllFailedTaskResponse), nil
	}
}

type DeleteFailedTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteFailedTaskInvoker) Invoke() (*model.DeleteFailedTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteFailedTaskResponse), nil
	}
}

type DeleteKeypairInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteKeypairInvoker) Invoke() (*model.DeleteKeypairResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteKeypairResponse), nil
	}
}

type DisassociateKeypairInvoker struct {
	*invoker.BaseInvoker
}

func (i *DisassociateKeypairInvoker) Invoke() (*model.DisassociateKeypairResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DisassociateKeypairResponse), nil
	}
}

type ListFailedTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListFailedTaskInvoker) Invoke() (*model.ListFailedTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListFailedTaskResponse), nil
	}
}

type ListKeypairDetailInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListKeypairDetailInvoker) Invoke() (*model.ListKeypairDetailResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListKeypairDetailResponse), nil
	}
}

type ListKeypairTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListKeypairTaskInvoker) Invoke() (*model.ListKeypairTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListKeypairTaskResponse), nil
	}
}

type ListKeypairsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListKeypairsInvoker) Invoke() (*model.ListKeypairsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListKeypairsResponse), nil
	}
}

type ListRunningTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListRunningTaskInvoker) Invoke() (*model.ListRunningTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListRunningTaskResponse), nil
	}
}

type UpdateKeypairDescriptionInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateKeypairDescriptionInvoker) Invoke() (*model.UpdateKeypairDescriptionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateKeypairDescriptionResponse), nil
	}
}
