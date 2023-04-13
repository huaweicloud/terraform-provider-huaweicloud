package v2

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/kafka/v2/model"
)

type BatchCreateOrDeleteKafkaTagInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchCreateOrDeleteKafkaTagInvoker) Invoke() (*model.BatchCreateOrDeleteKafkaTagResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchCreateOrDeleteKafkaTagResponse), nil
	}
}

type BatchDeleteGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchDeleteGroupInvoker) Invoke() (*model.BatchDeleteGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchDeleteGroupResponse), nil
	}
}

type BatchDeleteInstanceTopicInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchDeleteInstanceTopicInvoker) Invoke() (*model.BatchDeleteInstanceTopicResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchDeleteInstanceTopicResponse), nil
	}
}

type BatchDeleteInstanceUsersInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchDeleteInstanceUsersInvoker) Invoke() (*model.BatchDeleteInstanceUsersResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchDeleteInstanceUsersResponse), nil
	}
}

type BatchRestartOrDeleteInstancesInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchRestartOrDeleteInstancesInvoker) Invoke() (*model.BatchRestartOrDeleteInstancesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchRestartOrDeleteInstancesResponse), nil
	}
}

type CreateConnectorInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateConnectorInvoker) Invoke() (*model.CreateConnectorResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateConnectorResponse), nil
	}
}

type CreateInstanceByEngineInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateInstanceByEngineInvoker) Invoke() (*model.CreateInstanceByEngineResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateInstanceByEngineResponse), nil
	}
}

type CreateInstanceTopicInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateInstanceTopicInvoker) Invoke() (*model.CreateInstanceTopicResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateInstanceTopicResponse), nil
	}
}

type CreateInstanceUserInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateInstanceUserInvoker) Invoke() (*model.CreateInstanceUserResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateInstanceUserResponse), nil
	}
}

type CreatePartitionInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreatePartitionInvoker) Invoke() (*model.CreatePartitionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreatePartitionResponse), nil
	}
}

type CreatePostPaidInstanceInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreatePostPaidInstanceInvoker) Invoke() (*model.CreatePostPaidInstanceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreatePostPaidInstanceResponse), nil
	}
}

type CreateReassignmentTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateReassignmentTaskInvoker) Invoke() (*model.CreateReassignmentTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateReassignmentTaskResponse), nil
	}
}

type CreateSinkTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateSinkTaskInvoker) Invoke() (*model.CreateSinkTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateSinkTaskResponse), nil
	}
}

type DeleteBackgroundTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteBackgroundTaskInvoker) Invoke() (*model.DeleteBackgroundTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteBackgroundTaskResponse), nil
	}
}

type DeleteInstanceInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteInstanceInvoker) Invoke() (*model.DeleteInstanceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteInstanceResponse), nil
	}
}

type DeleteSinkTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteSinkTaskInvoker) Invoke() (*model.DeleteSinkTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteSinkTaskResponse), nil
	}
}

type ListAvailableZonesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAvailableZonesInvoker) Invoke() (*model.ListAvailableZonesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAvailableZonesResponse), nil
	}
}

type ListBackgroundTasksInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListBackgroundTasksInvoker) Invoke() (*model.ListBackgroundTasksResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListBackgroundTasksResponse), nil
	}
}

type ListEngineProductsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListEngineProductsInvoker) Invoke() (*model.ListEngineProductsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListEngineProductsResponse), nil
	}
}

type ListInstanceConsumerGroupsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListInstanceConsumerGroupsInvoker) Invoke() (*model.ListInstanceConsumerGroupsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListInstanceConsumerGroupsResponse), nil
	}
}

type ListInstanceTopicsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListInstanceTopicsInvoker) Invoke() (*model.ListInstanceTopicsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListInstanceTopicsResponse), nil
	}
}

type ListInstancesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListInstancesInvoker) Invoke() (*model.ListInstancesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListInstancesResponse), nil
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

type ListSinkTasksInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListSinkTasksInvoker) Invoke() (*model.ListSinkTasksResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListSinkTasksResponse), nil
	}
}

type ResetManagerPasswordInvoker struct {
	*invoker.BaseInvoker
}

func (i *ResetManagerPasswordInvoker) Invoke() (*model.ResetManagerPasswordResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ResetManagerPasswordResponse), nil
	}
}

type ResetMessageOffsetInvoker struct {
	*invoker.BaseInvoker
}

func (i *ResetMessageOffsetInvoker) Invoke() (*model.ResetMessageOffsetResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ResetMessageOffsetResponse), nil
	}
}

type ResetPasswordInvoker struct {
	*invoker.BaseInvoker
}

func (i *ResetPasswordInvoker) Invoke() (*model.ResetPasswordResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ResetPasswordResponse), nil
	}
}

type ResetUserPasswrodInvoker struct {
	*invoker.BaseInvoker
}

func (i *ResetUserPasswrodInvoker) Invoke() (*model.ResetUserPasswrodResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ResetUserPasswrodResponse), nil
	}
}

type ResizeEngineInstanceInvoker struct {
	*invoker.BaseInvoker
}

func (i *ResizeEngineInstanceInvoker) Invoke() (*model.ResizeEngineInstanceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ResizeEngineInstanceResponse), nil
	}
}

type ResizeInstanceInvoker struct {
	*invoker.BaseInvoker
}

func (i *ResizeInstanceInvoker) Invoke() (*model.ResizeInstanceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ResizeInstanceResponse), nil
	}
}

type RestartManagerInvoker struct {
	*invoker.BaseInvoker
}

func (i *RestartManagerInvoker) Invoke() (*model.RestartManagerResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RestartManagerResponse), nil
	}
}

type ShowBackgroundTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowBackgroundTaskInvoker) Invoke() (*model.ShowBackgroundTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowBackgroundTaskResponse), nil
	}
}

type ShowCesHierarchyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowCesHierarchyInvoker) Invoke() (*model.ShowCesHierarchyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowCesHierarchyResponse), nil
	}
}

type ShowClusterInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowClusterInvoker) Invoke() (*model.ShowClusterResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowClusterResponse), nil
	}
}

type ShowCoordinatorsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowCoordinatorsInvoker) Invoke() (*model.ShowCoordinatorsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowCoordinatorsResponse), nil
	}
}

type ShowEngineInstanceExtendProductInfoInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowEngineInstanceExtendProductInfoInvoker) Invoke() (*model.ShowEngineInstanceExtendProductInfoResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowEngineInstanceExtendProductInfoResponse), nil
	}
}

type ShowGroupsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowGroupsInvoker) Invoke() (*model.ShowGroupsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowGroupsResponse), nil
	}
}

type ShowInstanceInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowInstanceInvoker) Invoke() (*model.ShowInstanceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowInstanceResponse), nil
	}
}

type ShowInstanceExtendProductInfoInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowInstanceExtendProductInfoInvoker) Invoke() (*model.ShowInstanceExtendProductInfoResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowInstanceExtendProductInfoResponse), nil
	}
}

type ShowInstanceMessagesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowInstanceMessagesInvoker) Invoke() (*model.ShowInstanceMessagesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowInstanceMessagesResponse), nil
	}
}

type ShowInstanceTopicDetailInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowInstanceTopicDetailInvoker) Invoke() (*model.ShowInstanceTopicDetailResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowInstanceTopicDetailResponse), nil
	}
}

type ShowInstanceUsersInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowInstanceUsersInvoker) Invoke() (*model.ShowInstanceUsersResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowInstanceUsersResponse), nil
	}
}

type ShowKafkaProjectTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowKafkaProjectTagsInvoker) Invoke() (*model.ShowKafkaProjectTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowKafkaProjectTagsResponse), nil
	}
}

type ShowKafkaTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowKafkaTagsInvoker) Invoke() (*model.ShowKafkaTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowKafkaTagsResponse), nil
	}
}

type ShowKafkaTopicPartitionDiskusageInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowKafkaTopicPartitionDiskusageInvoker) Invoke() (*model.ShowKafkaTopicPartitionDiskusageResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowKafkaTopicPartitionDiskusageResponse), nil
	}
}

type ShowMaintainWindowsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowMaintainWindowsInvoker) Invoke() (*model.ShowMaintainWindowsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowMaintainWindowsResponse), nil
	}
}

type ShowMessagesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowMessagesInvoker) Invoke() (*model.ShowMessagesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowMessagesResponse), nil
	}
}

type ShowPartitionBeginningMessageInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowPartitionBeginningMessageInvoker) Invoke() (*model.ShowPartitionBeginningMessageResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowPartitionBeginningMessageResponse), nil
	}
}

type ShowPartitionEndMessageInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowPartitionEndMessageInvoker) Invoke() (*model.ShowPartitionEndMessageResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowPartitionEndMessageResponse), nil
	}
}

type ShowPartitionMessageInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowPartitionMessageInvoker) Invoke() (*model.ShowPartitionMessageResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowPartitionMessageResponse), nil
	}
}

type ShowSinkTaskDetailInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowSinkTaskDetailInvoker) Invoke() (*model.ShowSinkTaskDetailResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowSinkTaskDetailResponse), nil
	}
}

type ShowTopicAccessPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowTopicAccessPolicyInvoker) Invoke() (*model.ShowTopicAccessPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowTopicAccessPolicyResponse), nil
	}
}

type UpdateInstanceInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateInstanceInvoker) Invoke() (*model.UpdateInstanceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateInstanceResponse), nil
	}
}

type UpdateInstanceAutoCreateTopicInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateInstanceAutoCreateTopicInvoker) Invoke() (*model.UpdateInstanceAutoCreateTopicResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateInstanceAutoCreateTopicResponse), nil
	}
}

type UpdateInstanceCrossVpcIpInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateInstanceCrossVpcIpInvoker) Invoke() (*model.UpdateInstanceCrossVpcIpResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateInstanceCrossVpcIpResponse), nil
	}
}

type UpdateInstanceTopicInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateInstanceTopicInvoker) Invoke() (*model.UpdateInstanceTopicResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateInstanceTopicResponse), nil
	}
}

type UpdateSinkTaskQuotaInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateSinkTaskQuotaInvoker) Invoke() (*model.UpdateSinkTaskQuotaResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateSinkTaskQuotaResponse), nil
	}
}

type UpdateTopicAccessPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateTopicAccessPolicyInvoker) Invoke() (*model.UpdateTopicAccessPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateTopicAccessPolicyResponse), nil
	}
}

type UpdateTopicReplicaInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateTopicReplicaInvoker) Invoke() (*model.UpdateTopicReplicaResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateTopicReplicaResponse), nil
	}
}
