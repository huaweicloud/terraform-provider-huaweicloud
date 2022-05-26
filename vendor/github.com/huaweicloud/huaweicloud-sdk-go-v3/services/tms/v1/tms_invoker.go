package v1

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/tms/v1/model"
)

type CreatePredefineTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreatePredefineTagsInvoker) Invoke() (*model.CreatePredefineTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreatePredefineTagsResponse), nil
	}
}

type DeletePredefineTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeletePredefineTagsInvoker) Invoke() (*model.DeletePredefineTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeletePredefineTagsResponse), nil
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

type ListPredefineTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListPredefineTagsInvoker) Invoke() (*model.ListPredefineTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListPredefineTagsResponse), nil
	}
}

type ShowApiVersionInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowApiVersionInvoker) Invoke() (*model.ShowApiVersionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowApiVersionResponse), nil
	}
}

type ShowTagQuotaInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowTagQuotaInvoker) Invoke() (*model.ShowTagQuotaResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowTagQuotaResponse), nil
	}
}

type UpdatePredefineTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdatePredefineTagsInvoker) Invoke() (*model.UpdatePredefineTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdatePredefineTagsResponse), nil
	}
}
