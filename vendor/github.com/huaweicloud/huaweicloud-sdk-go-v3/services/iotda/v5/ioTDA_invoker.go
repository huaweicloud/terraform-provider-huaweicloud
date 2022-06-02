package v5

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"
)

type CreateAccessCodeInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateAccessCodeInvoker) Invoke() (*model.CreateAccessCodeResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateAccessCodeResponse), nil
	}
}

type AddQueueInvoker struct {
	*invoker.BaseInvoker
}

func (i *AddQueueInvoker) Invoke() (*model.AddQueueResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AddQueueResponse), nil
	}
}

type BatchShowQueueInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchShowQueueInvoker) Invoke() (*model.BatchShowQueueResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchShowQueueResponse), nil
	}
}

type DeleteQueueInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteQueueInvoker) Invoke() (*model.DeleteQueueResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteQueueResponse), nil
	}
}

type ShowQueueInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowQueueInvoker) Invoke() (*model.ShowQueueResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowQueueResponse), nil
	}
}

type AddApplicationInvoker struct {
	*invoker.BaseInvoker
}

func (i *AddApplicationInvoker) Invoke() (*model.AddApplicationResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AddApplicationResponse), nil
	}
}

type DeleteApplicationInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteApplicationInvoker) Invoke() (*model.DeleteApplicationResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteApplicationResponse), nil
	}
}

type ShowApplicationInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowApplicationInvoker) Invoke() (*model.ShowApplicationResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowApplicationResponse), nil
	}
}

type ShowApplicationsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowApplicationsInvoker) Invoke() (*model.ShowApplicationsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowApplicationsResponse), nil
	}
}

type CreateAsyncCommandInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateAsyncCommandInvoker) Invoke() (*model.CreateAsyncCommandResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateAsyncCommandResponse), nil
	}
}

type ShowAsyncDeviceCommandInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAsyncDeviceCommandInvoker) Invoke() (*model.ShowAsyncDeviceCommandResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAsyncDeviceCommandResponse), nil
	}
}

type CreateBatchTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateBatchTaskInvoker) Invoke() (*model.CreateBatchTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateBatchTaskResponse), nil
	}
}

type ListBatchTasksInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListBatchTasksInvoker) Invoke() (*model.ListBatchTasksResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListBatchTasksResponse), nil
	}
}

type ShowBatchTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowBatchTaskInvoker) Invoke() (*model.ShowBatchTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowBatchTaskResponse), nil
	}
}

type DeleteBatchTaskFileInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteBatchTaskFileInvoker) Invoke() (*model.DeleteBatchTaskFileResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteBatchTaskFileResponse), nil
	}
}

type ListBatchTaskFilesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListBatchTaskFilesInvoker) Invoke() (*model.ListBatchTaskFilesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListBatchTaskFilesResponse), nil
	}
}

type AddCertificateInvoker struct {
	*invoker.BaseInvoker
}

func (i *AddCertificateInvoker) Invoke() (*model.AddCertificateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AddCertificateResponse), nil
	}
}

type CheckCertificateInvoker struct {
	*invoker.BaseInvoker
}

func (i *CheckCertificateInvoker) Invoke() (*model.CheckCertificateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CheckCertificateResponse), nil
	}
}

type DeleteCertificateInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteCertificateInvoker) Invoke() (*model.DeleteCertificateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteCertificateResponse), nil
	}
}

type ListCertificatesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListCertificatesInvoker) Invoke() (*model.ListCertificatesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListCertificatesResponse), nil
	}
}

type CreateCommandInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateCommandInvoker) Invoke() (*model.CreateCommandResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateCommandResponse), nil
	}
}

type AddDeviceGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *AddDeviceGroupInvoker) Invoke() (*model.AddDeviceGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AddDeviceGroupResponse), nil
	}
}

type CreateOrDeleteDeviceInGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateOrDeleteDeviceInGroupInvoker) Invoke() (*model.CreateOrDeleteDeviceInGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateOrDeleteDeviceInGroupResponse), nil
	}
}

type DeleteDeviceGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteDeviceGroupInvoker) Invoke() (*model.DeleteDeviceGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteDeviceGroupResponse), nil
	}
}

type ListDeviceGroupsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListDeviceGroupsInvoker) Invoke() (*model.ListDeviceGroupsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListDeviceGroupsResponse), nil
	}
}

type ShowDeviceGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowDeviceGroupInvoker) Invoke() (*model.ShowDeviceGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowDeviceGroupResponse), nil
	}
}

type ShowDevicesInGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowDevicesInGroupInvoker) Invoke() (*model.ShowDevicesInGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowDevicesInGroupResponse), nil
	}
}

type UpdateDeviceGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateDeviceGroupInvoker) Invoke() (*model.UpdateDeviceGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateDeviceGroupResponse), nil
	}
}

type AddDeviceInvoker struct {
	*invoker.BaseInvoker
}

func (i *AddDeviceInvoker) Invoke() (*model.AddDeviceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AddDeviceResponse), nil
	}
}

type DeleteDeviceInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteDeviceInvoker) Invoke() (*model.DeleteDeviceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteDeviceResponse), nil
	}
}

type FreezeDeviceInvoker struct {
	*invoker.BaseInvoker
}

func (i *FreezeDeviceInvoker) Invoke() (*model.FreezeDeviceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.FreezeDeviceResponse), nil
	}
}

type ListDevicesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListDevicesInvoker) Invoke() (*model.ListDevicesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListDevicesResponse), nil
	}
}

type ResetDeviceSecretInvoker struct {
	*invoker.BaseInvoker
}

func (i *ResetDeviceSecretInvoker) Invoke() (*model.ResetDeviceSecretResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ResetDeviceSecretResponse), nil
	}
}

type ResetFingerprintInvoker struct {
	*invoker.BaseInvoker
}

func (i *ResetFingerprintInvoker) Invoke() (*model.ResetFingerprintResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ResetFingerprintResponse), nil
	}
}

type ShowDeviceInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowDeviceInvoker) Invoke() (*model.ShowDeviceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowDeviceResponse), nil
	}
}

type UnfreezeDeviceInvoker struct {
	*invoker.BaseInvoker
}

func (i *UnfreezeDeviceInvoker) Invoke() (*model.UnfreezeDeviceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UnfreezeDeviceResponse), nil
	}
}

type UpdateDeviceInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateDeviceInvoker) Invoke() (*model.UpdateDeviceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateDeviceResponse), nil
	}
}

type ShowDeviceShadowInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowDeviceShadowInvoker) Invoke() (*model.ShowDeviceShadowResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowDeviceShadowResponse), nil
	}
}

type UpdateDeviceShadowDesiredDataInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateDeviceShadowDesiredDataInvoker) Invoke() (*model.UpdateDeviceShadowDesiredDataResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateDeviceShadowDesiredDataResponse), nil
	}
}

type CreateMessageInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateMessageInvoker) Invoke() (*model.CreateMessageResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateMessageResponse), nil
	}
}

type ListDeviceMessagesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListDeviceMessagesInvoker) Invoke() (*model.ListDeviceMessagesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListDeviceMessagesResponse), nil
	}
}

type ShowDeviceMessageInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowDeviceMessageInvoker) Invoke() (*model.ShowDeviceMessageResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowDeviceMessageResponse), nil
	}
}

type CreateProductInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateProductInvoker) Invoke() (*model.CreateProductResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateProductResponse), nil
	}
}

type DeleteProductInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteProductInvoker) Invoke() (*model.DeleteProductResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteProductResponse), nil
	}
}

type ListProductsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListProductsInvoker) Invoke() (*model.ListProductsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListProductsResponse), nil
	}
}

type ShowProductInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowProductInvoker) Invoke() (*model.ShowProductResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowProductResponse), nil
	}
}

type UpdateProductInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateProductInvoker) Invoke() (*model.UpdateProductResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateProductResponse), nil
	}
}

type ListPropertiesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListPropertiesInvoker) Invoke() (*model.ListPropertiesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListPropertiesResponse), nil
	}
}

type UpdatePropertiesInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdatePropertiesInvoker) Invoke() (*model.UpdatePropertiesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdatePropertiesResponse), nil
	}
}

type CreateRoutingRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateRoutingRuleInvoker) Invoke() (*model.CreateRoutingRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateRoutingRuleResponse), nil
	}
}

type CreateRuleActionInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateRuleActionInvoker) Invoke() (*model.CreateRuleActionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateRuleActionResponse), nil
	}
}

type DeleteRoutingRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteRoutingRuleInvoker) Invoke() (*model.DeleteRoutingRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteRoutingRuleResponse), nil
	}
}

type DeleteRuleActionInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteRuleActionInvoker) Invoke() (*model.DeleteRuleActionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteRuleActionResponse), nil
	}
}

type ListRoutingRulesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListRoutingRulesInvoker) Invoke() (*model.ListRoutingRulesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListRoutingRulesResponse), nil
	}
}

type ListRuleActionsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListRuleActionsInvoker) Invoke() (*model.ListRuleActionsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListRuleActionsResponse), nil
	}
}

type ShowRoutingRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowRoutingRuleInvoker) Invoke() (*model.ShowRoutingRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowRoutingRuleResponse), nil
	}
}

type ShowRuleActionInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowRuleActionInvoker) Invoke() (*model.ShowRuleActionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowRuleActionResponse), nil
	}
}

type UpdateRoutingRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateRoutingRuleInvoker) Invoke() (*model.UpdateRoutingRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateRoutingRuleResponse), nil
	}
}

type UpdateRuleActionInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateRuleActionInvoker) Invoke() (*model.UpdateRuleActionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateRuleActionResponse), nil
	}
}

type ChangeRuleStatusInvoker struct {
	*invoker.BaseInvoker
}

func (i *ChangeRuleStatusInvoker) Invoke() (*model.ChangeRuleStatusResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ChangeRuleStatusResponse), nil
	}
}

type CreateRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateRuleInvoker) Invoke() (*model.CreateRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateRuleResponse), nil
	}
}

type DeleteRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteRuleInvoker) Invoke() (*model.DeleteRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteRuleResponse), nil
	}
}

type ListRulesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListRulesInvoker) Invoke() (*model.ListRulesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListRulesResponse), nil
	}
}

type ShowRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowRuleInvoker) Invoke() (*model.ShowRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowRuleResponse), nil
	}
}

type UpdateRuleInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateRuleInvoker) Invoke() (*model.UpdateRuleResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateRuleResponse), nil
	}
}

type ListResourcesByTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListResourcesByTagsInvoker) Invoke() (*model.ListResourcesByTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListResourcesByTagsResponse), nil
	}
}

type TagDeviceInvoker struct {
	*invoker.BaseInvoker
}

func (i *TagDeviceInvoker) Invoke() (*model.TagDeviceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.TagDeviceResponse), nil
	}
}

type UntagDeviceInvoker struct {
	*invoker.BaseInvoker
}

func (i *UntagDeviceInvoker) Invoke() (*model.UntagDeviceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UntagDeviceResponse), nil
	}
}
