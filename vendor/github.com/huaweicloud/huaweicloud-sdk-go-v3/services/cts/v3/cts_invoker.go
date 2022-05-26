package v3

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cts/v3/model"
)

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
