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

type UpdateApplicationInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateApplicationInvoker) Invoke() (*model.UpdateApplicationResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateApplicationResponse), nil
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

type CreateRoutingBacklogPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateRoutingBacklogPolicyInvoker) Invoke() (*model.CreateRoutingBacklogPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateRoutingBacklogPolicyResponse), nil
	}
}

type DeleteRoutingBacklogPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteRoutingBacklogPolicyInvoker) Invoke() (*model.DeleteRoutingBacklogPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteRoutingBacklogPolicyResponse), nil
	}
}

type ListRoutingBacklogPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListRoutingBacklogPolicyInvoker) Invoke() (*model.ListRoutingBacklogPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListRoutingBacklogPolicyResponse), nil
	}
}

type ShowRoutingBacklogPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowRoutingBacklogPolicyInvoker) Invoke() (*model.ShowRoutingBacklogPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowRoutingBacklogPolicyResponse), nil
	}
}

type UpdateRoutingBacklogPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateRoutingBacklogPolicyInvoker) Invoke() (*model.UpdateRoutingBacklogPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateRoutingBacklogPolicyResponse), nil
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

type DeleteBatchTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteBatchTaskInvoker) Invoke() (*model.DeleteBatchTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteBatchTaskResponse), nil
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

type RetryBatchTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *RetryBatchTaskInvoker) Invoke() (*model.RetryBatchTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RetryBatchTaskResponse), nil
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

type StopBatchTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *StopBatchTaskInvoker) Invoke() (*model.StopBatchTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StopBatchTaskResponse), nil
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

type UploadBatchTaskFileInvoker struct {
	*invoker.BaseInvoker
}

func (i *UploadBatchTaskFileInvoker) Invoke() (*model.UploadBatchTaskFileResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UploadBatchTaskFileResponse), nil
	}
}

type AddBridgeInvoker struct {
	*invoker.BaseInvoker
}

func (i *AddBridgeInvoker) Invoke() (*model.AddBridgeResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AddBridgeResponse), nil
	}
}

type DeleteBridgeInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteBridgeInvoker) Invoke() (*model.DeleteBridgeResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteBridgeResponse), nil
	}
}

type ListBridgesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListBridgesInvoker) Invoke() (*model.ListBridgesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListBridgesResponse), nil
	}
}

type ResetBridgeSecretInvoker struct {
	*invoker.BaseInvoker
}

func (i *ResetBridgeSecretInvoker) Invoke() (*model.ResetBridgeSecretResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ResetBridgeSecretResponse), nil
	}
}

type BroadcastMessageInvoker struct {
	*invoker.BaseInvoker
}

func (i *BroadcastMessageInvoker) Invoke() (*model.BroadcastMessageResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BroadcastMessageResponse), nil
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

type UpdateCertificateInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateCertificateInvoker) Invoke() (*model.UpdateCertificateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateCertificateResponse), nil
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

type CreateDeviceAuthorizerInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateDeviceAuthorizerInvoker) Invoke() (*model.CreateDeviceAuthorizerResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateDeviceAuthorizerResponse), nil
	}
}

type DeleteDeviceAuthorizerInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteDeviceAuthorizerInvoker) Invoke() (*model.DeleteDeviceAuthorizerResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteDeviceAuthorizerResponse), nil
	}
}

type ListDeviceAuthorizersInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListDeviceAuthorizersInvoker) Invoke() (*model.ListDeviceAuthorizersResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListDeviceAuthorizersResponse), nil
	}
}

type ShowDeviceAuthorizerInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowDeviceAuthorizerInvoker) Invoke() (*model.ShowDeviceAuthorizerResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowDeviceAuthorizerResponse), nil
	}
}

type UpdateDeviceAuthorizerInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateDeviceAuthorizerInvoker) Invoke() (*model.UpdateDeviceAuthorizerResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateDeviceAuthorizerResponse), nil
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

type ListDeviceGroupsByDeviceInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListDeviceGroupsByDeviceInvoker) Invoke() (*model.ListDeviceGroupsByDeviceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListDeviceGroupsByDeviceResponse), nil
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

type SearchDevicesInvoker struct {
	*invoker.BaseInvoker
}

func (i *SearchDevicesInvoker) Invoke() (*model.SearchDevicesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.SearchDevicesResponse), nil
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

type CreateDeviceProxyInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateDeviceProxyInvoker) Invoke() (*model.CreateDeviceProxyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateDeviceProxyResponse), nil
	}
}

type DeleteDeviceProxyInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteDeviceProxyInvoker) Invoke() (*model.DeleteDeviceProxyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteDeviceProxyResponse), nil
	}
}

type ListDeviceProxiesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListDeviceProxiesInvoker) Invoke() (*model.ListDeviceProxiesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListDeviceProxiesResponse), nil
	}
}

type ShowDeviceProxyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowDeviceProxyInvoker) Invoke() (*model.ShowDeviceProxyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowDeviceProxyResponse), nil
	}
}

type UpdateDeviceProxyInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateDeviceProxyInvoker) Invoke() (*model.UpdateDeviceProxyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateDeviceProxyResponse), nil
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

type CreateRoutingFlowControlPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateRoutingFlowControlPolicyInvoker) Invoke() (*model.CreateRoutingFlowControlPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateRoutingFlowControlPolicyResponse), nil
	}
}

type DeleteRoutingFlowControlPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteRoutingFlowControlPolicyInvoker) Invoke() (*model.DeleteRoutingFlowControlPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteRoutingFlowControlPolicyResponse), nil
	}
}

type ListRoutingFlowControlPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListRoutingFlowControlPolicyInvoker) Invoke() (*model.ListRoutingFlowControlPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListRoutingFlowControlPolicyResponse), nil
	}
}

type ShowRoutingFlowControlPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowRoutingFlowControlPolicyInvoker) Invoke() (*model.ShowRoutingFlowControlPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowRoutingFlowControlPolicyResponse), nil
	}
}

type UpdateRoutingFlowControlPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateRoutingFlowControlPolicyInvoker) Invoke() (*model.UpdateRoutingFlowControlPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateRoutingFlowControlPolicyResponse), nil
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

type CreateOtaPackageInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateOtaPackageInvoker) Invoke() (*model.CreateOtaPackageResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateOtaPackageResponse), nil
	}
}

type DeleteOtaPackageInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteOtaPackageInvoker) Invoke() (*model.DeleteOtaPackageResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteOtaPackageResponse), nil
	}
}

type ListOtaPackageInfoInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListOtaPackageInfoInvoker) Invoke() (*model.ListOtaPackageInfoResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListOtaPackageInfoResponse), nil
	}
}

type ShowOtaPackageInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowOtaPackageInvoker) Invoke() (*model.ShowOtaPackageResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowOtaPackageResponse), nil
	}
}

type BindDevicePolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *BindDevicePolicyInvoker) Invoke() (*model.BindDevicePolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BindDevicePolicyResponse), nil
	}
}

type CreateDevicePolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateDevicePolicyInvoker) Invoke() (*model.CreateDevicePolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateDevicePolicyResponse), nil
	}
}

type DeleteDevicePolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteDevicePolicyInvoker) Invoke() (*model.DeleteDevicePolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteDevicePolicyResponse), nil
	}
}

type ListDevicePoliciesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListDevicePoliciesInvoker) Invoke() (*model.ListDevicePoliciesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListDevicePoliciesResponse), nil
	}
}

type ShowDevicePolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowDevicePolicyInvoker) Invoke() (*model.ShowDevicePolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowDevicePolicyResponse), nil
	}
}

type ShowTargetsInDevicePolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowTargetsInDevicePolicyInvoker) Invoke() (*model.ShowTargetsInDevicePolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowTargetsInDevicePolicyResponse), nil
	}
}

type UnbindDevicePolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *UnbindDevicePolicyInvoker) Invoke() (*model.UnbindDevicePolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UnbindDevicePolicyResponse), nil
	}
}

type UpdateDevicePolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateDevicePolicyInvoker) Invoke() (*model.UpdateDevicePolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateDevicePolicyResponse), nil
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

type CreateProvisioningTemplateInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateProvisioningTemplateInvoker) Invoke() (*model.CreateProvisioningTemplateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateProvisioningTemplateResponse), nil
	}
}

type DeleteProvisioningTemplateInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteProvisioningTemplateInvoker) Invoke() (*model.DeleteProvisioningTemplateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteProvisioningTemplateResponse), nil
	}
}

type ListProvisioningTemplatesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListProvisioningTemplatesInvoker) Invoke() (*model.ListProvisioningTemplatesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListProvisioningTemplatesResponse), nil
	}
}

type ShowProvisioningTemplateInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowProvisioningTemplateInvoker) Invoke() (*model.ShowProvisioningTemplateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowProvisioningTemplateResponse), nil
	}
}

type UpdateProvisioningTemplateInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateProvisioningTemplateInvoker) Invoke() (*model.UpdateProvisioningTemplateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateProvisioningTemplateResponse), nil
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

type AddTunnelInvoker struct {
	*invoker.BaseInvoker
}

func (i *AddTunnelInvoker) Invoke() (*model.AddTunnelResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AddTunnelResponse), nil
	}
}

type CloseDeviceTunnelInvoker struct {
	*invoker.BaseInvoker
}

func (i *CloseDeviceTunnelInvoker) Invoke() (*model.CloseDeviceTunnelResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CloseDeviceTunnelResponse), nil
	}
}

type DeleteDeviceTunnelInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteDeviceTunnelInvoker) Invoke() (*model.DeleteDeviceTunnelResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteDeviceTunnelResponse), nil
	}
}

type ListDeviceTunnelsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListDeviceTunnelsInvoker) Invoke() (*model.ListDeviceTunnelsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListDeviceTunnelsResponse), nil
	}
}

type ShowDeviceTunnelInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowDeviceTunnelInvoker) Invoke() (*model.ShowDeviceTunnelResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowDeviceTunnelResponse), nil
	}
}
