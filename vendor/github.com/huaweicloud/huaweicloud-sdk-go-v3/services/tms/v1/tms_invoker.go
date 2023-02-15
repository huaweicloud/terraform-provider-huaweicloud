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

type CreateResourceTagInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateResourceTagInvoker) Invoke() (*model.CreateResourceTagResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateResourceTagResponse), nil
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

type DeleteResourceTagInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteResourceTagInvoker) Invoke() (*model.DeleteResourceTagResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteResourceTagResponse), nil
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

type ListProvidersInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListProvidersInvoker) Invoke() (*model.ListProvidersResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListProvidersResponse), nil
	}
}

type ListResourceInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListResourceInvoker) Invoke() (*model.ListResourceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListResourceResponse), nil
	}
}

type ListTagKeysInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListTagKeysInvoker) Invoke() (*model.ListTagKeysResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListTagKeysResponse), nil
	}
}

type ListTagValuesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListTagValuesInvoker) Invoke() (*model.ListTagValuesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListTagValuesResponse), nil
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

type ShowResourceTagInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowResourceTagInvoker) Invoke() (*model.ShowResourceTagResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowResourceTagResponse), nil
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
