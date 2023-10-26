package v3

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cts/v3/model"
)

type BatchCreateResourceTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchCreateResourceTagsInvoker) Invoke() (*model.BatchCreateResourceTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchCreateResourceTagsResponse), nil
	}
}

type BatchDeleteResourceTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchDeleteResourceTagsInvoker) Invoke() (*model.BatchDeleteResourceTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchDeleteResourceTagsResponse), nil
	}
}

type CheckObsBucketsInvoker struct {
	*invoker.BaseInvoker
}

func (i *CheckObsBucketsInvoker) Invoke() (*model.CheckObsBucketsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CheckObsBucketsResponse), nil
	}
}

type CreateNotificationInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateNotificationInvoker) Invoke() (*model.CreateNotificationResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateNotificationResponse), nil
	}
}

type CreateTrackerInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateTrackerInvoker) Invoke() (*model.CreateTrackerResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateTrackerResponse), nil
	}
}

type DeleteNotificationInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteNotificationInvoker) Invoke() (*model.DeleteNotificationResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteNotificationResponse), nil
	}
}

type DeleteTrackerInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteTrackerInvoker) Invoke() (*model.DeleteTrackerResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteTrackerResponse), nil
	}
}

type ListNotificationsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListNotificationsInvoker) Invoke() (*model.ListNotificationsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListNotificationsResponse), nil
	}
}

type ListOperationsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListOperationsInvoker) Invoke() (*model.ListOperationsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListOperationsResponse), nil
	}
}

type ListQuotasInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListQuotasInvoker) Invoke() (*model.ListQuotasResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListQuotasResponse), nil
	}
}

type ListTraceResourcesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListTraceResourcesInvoker) Invoke() (*model.ListTraceResourcesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListTraceResourcesResponse), nil
	}
}

type ListTracesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListTracesInvoker) Invoke() (*model.ListTracesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListTracesResponse), nil
	}
}

type ListTrackersInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListTrackersInvoker) Invoke() (*model.ListTrackersResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListTrackersResponse), nil
	}
}

type ListUserResourcesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListUserResourcesInvoker) Invoke() (*model.ListUserResourcesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListUserResourcesResponse), nil
	}
}

type UpdateNotificationInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateNotificationInvoker) Invoke() (*model.UpdateNotificationResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateNotificationResponse), nil
	}
}

type UpdateTrackerInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateTrackerInvoker) Invoke() (*model.UpdateTrackerResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateTrackerResponse), nil
	}
}
