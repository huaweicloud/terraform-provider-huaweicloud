package v3

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cce/v3/model"
)

type AddNodeInvoker struct {
	*invoker.BaseInvoker
}

func (i *AddNodeInvoker) Invoke() (*model.AddNodeResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AddNodeResponse), nil
	}
}

type AddNodesToNodePoolInvoker struct {
	*invoker.BaseInvoker
}

func (i *AddNodesToNodePoolInvoker) Invoke() (*model.AddNodesToNodePoolResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AddNodesToNodePoolResponse), nil
	}
}

type AwakeClusterInvoker struct {
	*invoker.BaseInvoker
}

func (i *AwakeClusterInvoker) Invoke() (*model.AwakeClusterResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AwakeClusterResponse), nil
	}
}

type BatchCreateClusterTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchCreateClusterTagsInvoker) Invoke() (*model.BatchCreateClusterTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchCreateClusterTagsResponse), nil
	}
}

type BatchDeleteClusterTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchDeleteClusterTagsInvoker) Invoke() (*model.BatchDeleteClusterTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchDeleteClusterTagsResponse), nil
	}
}

type BatchSyncNodesInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchSyncNodesInvoker) Invoke() (*model.BatchSyncNodesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchSyncNodesResponse), nil
	}
}

type ContinueUpgradeClusterTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *ContinueUpgradeClusterTaskInvoker) Invoke() (*model.ContinueUpgradeClusterTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ContinueUpgradeClusterTaskResponse), nil
	}
}

type CreateAddonInstanceInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateAddonInstanceInvoker) Invoke() (*model.CreateAddonInstanceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateAddonInstanceResponse), nil
	}
}

type CreateCloudPersistentVolumeClaimsInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateCloudPersistentVolumeClaimsInvoker) Invoke() (*model.CreateCloudPersistentVolumeClaimsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateCloudPersistentVolumeClaimsResponse), nil
	}
}

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

type CreateClusterMasterSnapshotInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateClusterMasterSnapshotInvoker) Invoke() (*model.CreateClusterMasterSnapshotResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateClusterMasterSnapshotResponse), nil
	}
}

type CreateKubernetesClusterCertInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateKubernetesClusterCertInvoker) Invoke() (*model.CreateKubernetesClusterCertResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateKubernetesClusterCertResponse), nil
	}
}

type CreateNodeInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateNodeInvoker) Invoke() (*model.CreateNodeResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateNodeResponse), nil
	}
}

type CreateNodePoolInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateNodePoolInvoker) Invoke() (*model.CreateNodePoolResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateNodePoolResponse), nil
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

type CreatePostCheckInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreatePostCheckInvoker) Invoke() (*model.CreatePostCheckResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreatePostCheckResponse), nil
	}
}

type CreatePreCheckInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreatePreCheckInvoker) Invoke() (*model.CreatePreCheckResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreatePreCheckResponse), nil
	}
}

type CreateReleaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateReleaseInvoker) Invoke() (*model.CreateReleaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateReleaseResponse), nil
	}
}

type CreateUpgradeWorkFlowInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateUpgradeWorkFlowInvoker) Invoke() (*model.CreateUpgradeWorkFlowResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateUpgradeWorkFlowResponse), nil
	}
}

type DeleteAddonInstanceInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteAddonInstanceInvoker) Invoke() (*model.DeleteAddonInstanceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteAddonInstanceResponse), nil
	}
}

type DeleteChartInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteChartInvoker) Invoke() (*model.DeleteChartResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteChartResponse), nil
	}
}

type DeleteCloudPersistentVolumeClaimsInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteCloudPersistentVolumeClaimsInvoker) Invoke() (*model.DeleteCloudPersistentVolumeClaimsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteCloudPersistentVolumeClaimsResponse), nil
	}
}

type DeleteClusterInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteClusterInvoker) Invoke() (*model.DeleteClusterResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteClusterResponse), nil
	}
}

type DeleteNodeInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteNodeInvoker) Invoke() (*model.DeleteNodeResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteNodeResponse), nil
	}
}

type DeleteNodePoolInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteNodePoolInvoker) Invoke() (*model.DeleteNodePoolResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteNodePoolResponse), nil
	}
}

type DeleteReleaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteReleaseInvoker) Invoke() (*model.DeleteReleaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteReleaseResponse), nil
	}
}

type DownloadChartInvoker struct {
	*invoker.BaseInvoker
}

func (i *DownloadChartInvoker) Invoke() (*model.DownloadChartResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DownloadChartResponse), nil
	}
}

type HibernateClusterInvoker struct {
	*invoker.BaseInvoker
}

func (i *HibernateClusterInvoker) Invoke() (*model.HibernateClusterResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.HibernateClusterResponse), nil
	}
}

type ListAddonInstancesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAddonInstancesInvoker) Invoke() (*model.ListAddonInstancesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAddonInstancesResponse), nil
	}
}

type ListAddonTemplatesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAddonTemplatesInvoker) Invoke() (*model.ListAddonTemplatesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAddonTemplatesResponse), nil
	}
}

type ListChartsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListChartsInvoker) Invoke() (*model.ListChartsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListChartsResponse), nil
	}
}

type ListClusterMasterSnapshotTasksInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListClusterMasterSnapshotTasksInvoker) Invoke() (*model.ListClusterMasterSnapshotTasksResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListClusterMasterSnapshotTasksResponse), nil
	}
}

type ListClusterUpgradeFeatureGatesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListClusterUpgradeFeatureGatesInvoker) Invoke() (*model.ListClusterUpgradeFeatureGatesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListClusterUpgradeFeatureGatesResponse), nil
	}
}

type ListClusterUpgradePathsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListClusterUpgradePathsInvoker) Invoke() (*model.ListClusterUpgradePathsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListClusterUpgradePathsResponse), nil
	}
}

type ListClustersInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListClustersInvoker) Invoke() (*model.ListClustersResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListClustersResponse), nil
	}
}

type ListNodePoolsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListNodePoolsInvoker) Invoke() (*model.ListNodePoolsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListNodePoolsResponse), nil
	}
}

type ListNodesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListNodesInvoker) Invoke() (*model.ListNodesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListNodesResponse), nil
	}
}

type ListPartitionsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListPartitionsInvoker) Invoke() (*model.ListPartitionsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListPartitionsResponse), nil
	}
}

type ListPreCheckTasksInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListPreCheckTasksInvoker) Invoke() (*model.ListPreCheckTasksResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListPreCheckTasksResponse), nil
	}
}

type ListReleasesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListReleasesInvoker) Invoke() (*model.ListReleasesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListReleasesResponse), nil
	}
}

type ListUpgradeClusterTasksInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListUpgradeClusterTasksInvoker) Invoke() (*model.ListUpgradeClusterTasksResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListUpgradeClusterTasksResponse), nil
	}
}

type ListUpgradeWorkFlowsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListUpgradeWorkFlowsInvoker) Invoke() (*model.ListUpgradeWorkFlowsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListUpgradeWorkFlowsResponse), nil
	}
}

type MigrateNodeInvoker struct {
	*invoker.BaseInvoker
}

func (i *MigrateNodeInvoker) Invoke() (*model.MigrateNodeResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.MigrateNodeResponse), nil
	}
}

type PauseUpgradeClusterTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *PauseUpgradeClusterTaskInvoker) Invoke() (*model.PauseUpgradeClusterTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.PauseUpgradeClusterTaskResponse), nil
	}
}

type RemoveNodeInvoker struct {
	*invoker.BaseInvoker
}

func (i *RemoveNodeInvoker) Invoke() (*model.RemoveNodeResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RemoveNodeResponse), nil
	}
}

type ResetNodeInvoker struct {
	*invoker.BaseInvoker
}

func (i *ResetNodeInvoker) Invoke() (*model.ResetNodeResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ResetNodeResponse), nil
	}
}

type ResizeClusterInvoker struct {
	*invoker.BaseInvoker
}

func (i *ResizeClusterInvoker) Invoke() (*model.ResizeClusterResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ResizeClusterResponse), nil
	}
}

type RetryUpgradeClusterTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *RetryUpgradeClusterTaskInvoker) Invoke() (*model.RetryUpgradeClusterTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RetryUpgradeClusterTaskResponse), nil
	}
}

type RevokeKubernetesClusterCertInvoker struct {
	*invoker.BaseInvoker
}

func (i *RevokeKubernetesClusterCertInvoker) Invoke() (*model.RevokeKubernetesClusterCertResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RevokeKubernetesClusterCertResponse), nil
	}
}

type RollbackAddonInstanceInvoker struct {
	*invoker.BaseInvoker
}

func (i *RollbackAddonInstanceInvoker) Invoke() (*model.RollbackAddonInstanceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RollbackAddonInstanceResponse), nil
	}
}

type ScaleNodePoolInvoker struct {
	*invoker.BaseInvoker
}

func (i *ScaleNodePoolInvoker) Invoke() (*model.ScaleNodePoolResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ScaleNodePoolResponse), nil
	}
}

type ShowAddonInstanceInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAddonInstanceInvoker) Invoke() (*model.ShowAddonInstanceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAddonInstanceResponse), nil
	}
}

type ShowChartInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowChartInvoker) Invoke() (*model.ShowChartResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowChartResponse), nil
	}
}

type ShowChartValuesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowChartValuesInvoker) Invoke() (*model.ShowChartValuesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowChartValuesResponse), nil
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

type ShowClusterConfigInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowClusterConfigInvoker) Invoke() (*model.ShowClusterConfigResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowClusterConfigResponse), nil
	}
}

type ShowClusterConfigurationDetailsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowClusterConfigurationDetailsInvoker) Invoke() (*model.ShowClusterConfigurationDetailsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowClusterConfigurationDetailsResponse), nil
	}
}

type ShowClusterEndpointsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowClusterEndpointsInvoker) Invoke() (*model.ShowClusterEndpointsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowClusterEndpointsResponse), nil
	}
}

type ShowClusterSupportConfigurationInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowClusterSupportConfigurationInvoker) Invoke() (*model.ShowClusterSupportConfigurationResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowClusterSupportConfigurationResponse), nil
	}
}

type ShowClusterUpgradeInfoInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowClusterUpgradeInfoInvoker) Invoke() (*model.ShowClusterUpgradeInfoResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowClusterUpgradeInfoResponse), nil
	}
}

type ShowJobInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowJobInvoker) Invoke() (*model.ShowJobResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowJobResponse), nil
	}
}

type ShowNodeInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowNodeInvoker) Invoke() (*model.ShowNodeResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowNodeResponse), nil
	}
}

type ShowNodePoolInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowNodePoolInvoker) Invoke() (*model.ShowNodePoolResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowNodePoolResponse), nil
	}
}

type ShowNodePoolConfigurationDetailsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowNodePoolConfigurationDetailsInvoker) Invoke() (*model.ShowNodePoolConfigurationDetailsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowNodePoolConfigurationDetailsResponse), nil
	}
}

type ShowNodePoolConfigurationsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowNodePoolConfigurationsInvoker) Invoke() (*model.ShowNodePoolConfigurationsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowNodePoolConfigurationsResponse), nil
	}
}

type ShowPartitionInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowPartitionInvoker) Invoke() (*model.ShowPartitionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowPartitionResponse), nil
	}
}

type ShowPreCheckInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowPreCheckInvoker) Invoke() (*model.ShowPreCheckResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowPreCheckResponse), nil
	}
}

type ShowQuotasInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowQuotasInvoker) Invoke() (*model.ShowQuotasResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowQuotasResponse), nil
	}
}

type ShowReleaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowReleaseInvoker) Invoke() (*model.ShowReleaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowReleaseResponse), nil
	}
}

type ShowReleaseHistoryInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowReleaseHistoryInvoker) Invoke() (*model.ShowReleaseHistoryResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowReleaseHistoryResponse), nil
	}
}

type ShowUpgradeClusterTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowUpgradeClusterTaskInvoker) Invoke() (*model.ShowUpgradeClusterTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowUpgradeClusterTaskResponse), nil
	}
}

type ShowUpgradeWorkFlowInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowUpgradeWorkFlowInvoker) Invoke() (*model.ShowUpgradeWorkFlowResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowUpgradeWorkFlowResponse), nil
	}
}

type ShowUserChartsQuotasInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowUserChartsQuotasInvoker) Invoke() (*model.ShowUserChartsQuotasResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowUserChartsQuotasResponse), nil
	}
}

type SyncNodeInvoker struct {
	*invoker.BaseInvoker
}

func (i *SyncNodeInvoker) Invoke() (*model.SyncNodeResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.SyncNodeResponse), nil
	}
}

type SyncNodePoolInvoker struct {
	*invoker.BaseInvoker
}

func (i *SyncNodePoolInvoker) Invoke() (*model.SyncNodePoolResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.SyncNodePoolResponse), nil
	}
}

type UpdateAddonInstanceInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateAddonInstanceInvoker) Invoke() (*model.UpdateAddonInstanceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateAddonInstanceResponse), nil
	}
}

type UpdateChartInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateChartInvoker) Invoke() (*model.UpdateChartResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateChartResponse), nil
	}
}

type UpdateClusterInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateClusterInvoker) Invoke() (*model.UpdateClusterResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateClusterResponse), nil
	}
}

type UpdateClusterEipInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateClusterEipInvoker) Invoke() (*model.UpdateClusterEipResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateClusterEipResponse), nil
	}
}

type UpdateClusterLogConfigInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateClusterLogConfigInvoker) Invoke() (*model.UpdateClusterLogConfigResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateClusterLogConfigResponse), nil
	}
}

type UpdateNodeInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateNodeInvoker) Invoke() (*model.UpdateNodeResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateNodeResponse), nil
	}
}

type UpdateNodePoolInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateNodePoolInvoker) Invoke() (*model.UpdateNodePoolResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateNodePoolResponse), nil
	}
}

type UpdateNodePoolConfigurationInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateNodePoolConfigurationInvoker) Invoke() (*model.UpdateNodePoolConfigurationResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateNodePoolConfigurationResponse), nil
	}
}

type UpdatePartitionInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdatePartitionInvoker) Invoke() (*model.UpdatePartitionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdatePartitionResponse), nil
	}
}

type UpdateReleaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateReleaseInvoker) Invoke() (*model.UpdateReleaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateReleaseResponse), nil
	}
}

type UpgradeClusterInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpgradeClusterInvoker) Invoke() (*model.UpgradeClusterResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpgradeClusterResponse), nil
	}
}

type UpgradeNodePoolInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpgradeNodePoolInvoker) Invoke() (*model.UpgradeNodePoolResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpgradeNodePoolResponse), nil
	}
}

type UpgradeWorkFlowUpdateInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpgradeWorkFlowUpdateInvoker) Invoke() (*model.UpgradeWorkFlowUpdateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpgradeWorkFlowUpdateResponse), nil
	}
}

type UploadChartInvoker struct {
	*invoker.BaseInvoker
}

func (i *UploadChartInvoker) Invoke() (*model.UploadChartResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UploadChartResponse), nil
	}
}

type ShowVersionInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowVersionInvoker) Invoke() (*model.ShowVersionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowVersionResponse), nil
	}
}

type BatchCreateAutopilotClusterTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchCreateAutopilotClusterTagsInvoker) Invoke() (*model.BatchCreateAutopilotClusterTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchCreateAutopilotClusterTagsResponse), nil
	}
}

type BatchDeleteAutopilotClusterTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *BatchDeleteAutopilotClusterTagsInvoker) Invoke() (*model.BatchDeleteAutopilotClusterTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.BatchDeleteAutopilotClusterTagsResponse), nil
	}
}

type CreateAutopilotAddonInstanceInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateAutopilotAddonInstanceInvoker) Invoke() (*model.CreateAutopilotAddonInstanceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateAutopilotAddonInstanceResponse), nil
	}
}

type CreateAutopilotClusterInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateAutopilotClusterInvoker) Invoke() (*model.CreateAutopilotClusterResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateAutopilotClusterResponse), nil
	}
}

type CreateAutopilotClusterMasterSnapshotInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateAutopilotClusterMasterSnapshotInvoker) Invoke() (*model.CreateAutopilotClusterMasterSnapshotResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateAutopilotClusterMasterSnapshotResponse), nil
	}
}

type CreateAutopilotKubernetesClusterCertInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateAutopilotKubernetesClusterCertInvoker) Invoke() (*model.CreateAutopilotKubernetesClusterCertResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateAutopilotKubernetesClusterCertResponse), nil
	}
}

type CreateAutopilotPostCheckInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateAutopilotPostCheckInvoker) Invoke() (*model.CreateAutopilotPostCheckResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateAutopilotPostCheckResponse), nil
	}
}

type CreateAutopilotPreCheckInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateAutopilotPreCheckInvoker) Invoke() (*model.CreateAutopilotPreCheckResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateAutopilotPreCheckResponse), nil
	}
}

type CreateAutopilotReleaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateAutopilotReleaseInvoker) Invoke() (*model.CreateAutopilotReleaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateAutopilotReleaseResponse), nil
	}
}

type CreateAutopilotUpgradeWorkFlowInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateAutopilotUpgradeWorkFlowInvoker) Invoke() (*model.CreateAutopilotUpgradeWorkFlowResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateAutopilotUpgradeWorkFlowResponse), nil
	}
}

type DeleteAutopilotAddonInstanceInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteAutopilotAddonInstanceInvoker) Invoke() (*model.DeleteAutopilotAddonInstanceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteAutopilotAddonInstanceResponse), nil
	}
}

type DeleteAutopilotChartInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteAutopilotChartInvoker) Invoke() (*model.DeleteAutopilotChartResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteAutopilotChartResponse), nil
	}
}

type DeleteAutopilotClusterInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteAutopilotClusterInvoker) Invoke() (*model.DeleteAutopilotClusterResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteAutopilotClusterResponse), nil
	}
}

type DeleteAutopilotReleaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteAutopilotReleaseInvoker) Invoke() (*model.DeleteAutopilotReleaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteAutopilotReleaseResponse), nil
	}
}

type DownloadAutopilotChartInvoker struct {
	*invoker.BaseInvoker
}

func (i *DownloadAutopilotChartInvoker) Invoke() (*model.DownloadAutopilotChartResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DownloadAutopilotChartResponse), nil
	}
}

type ListAutopilotAddonInstancesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAutopilotAddonInstancesInvoker) Invoke() (*model.ListAutopilotAddonInstancesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAutopilotAddonInstancesResponse), nil
	}
}

type ListAutopilotAddonTemplatesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAutopilotAddonTemplatesInvoker) Invoke() (*model.ListAutopilotAddonTemplatesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAutopilotAddonTemplatesResponse), nil
	}
}

type ListAutopilotChartsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAutopilotChartsInvoker) Invoke() (*model.ListAutopilotChartsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAutopilotChartsResponse), nil
	}
}

type ListAutopilotClusterMasterSnapshotTasksInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAutopilotClusterMasterSnapshotTasksInvoker) Invoke() (*model.ListAutopilotClusterMasterSnapshotTasksResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAutopilotClusterMasterSnapshotTasksResponse), nil
	}
}

type ListAutopilotClusterUpgradeFeatureGatesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAutopilotClusterUpgradeFeatureGatesInvoker) Invoke() (*model.ListAutopilotClusterUpgradeFeatureGatesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAutopilotClusterUpgradeFeatureGatesResponse), nil
	}
}

type ListAutopilotClusterUpgradePathsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAutopilotClusterUpgradePathsInvoker) Invoke() (*model.ListAutopilotClusterUpgradePathsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAutopilotClusterUpgradePathsResponse), nil
	}
}

type ListAutopilotClustersInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAutopilotClustersInvoker) Invoke() (*model.ListAutopilotClustersResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAutopilotClustersResponse), nil
	}
}

type ListAutopilotPreCheckTasksInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAutopilotPreCheckTasksInvoker) Invoke() (*model.ListAutopilotPreCheckTasksResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAutopilotPreCheckTasksResponse), nil
	}
}

type ListAutopilotReleasesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAutopilotReleasesInvoker) Invoke() (*model.ListAutopilotReleasesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAutopilotReleasesResponse), nil
	}
}

type ListAutopilotUpgradeClusterTasksInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAutopilotUpgradeClusterTasksInvoker) Invoke() (*model.ListAutopilotUpgradeClusterTasksResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAutopilotUpgradeClusterTasksResponse), nil
	}
}

type ListAutopilotUpgradeWorkFlowsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAutopilotUpgradeWorkFlowsInvoker) Invoke() (*model.ListAutopilotUpgradeWorkFlowsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAutopilotUpgradeWorkFlowsResponse), nil
	}
}

type RetryAutopilotUpgradeClusterTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *RetryAutopilotUpgradeClusterTaskInvoker) Invoke() (*model.RetryAutopilotUpgradeClusterTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RetryAutopilotUpgradeClusterTaskResponse), nil
	}
}

type RollbackAutopilotAddonInstanceInvoker struct {
	*invoker.BaseInvoker
}

func (i *RollbackAutopilotAddonInstanceInvoker) Invoke() (*model.RollbackAutopilotAddonInstanceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RollbackAutopilotAddonInstanceResponse), nil
	}
}

type ShowAutopilotAddonInstanceInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAutopilotAddonInstanceInvoker) Invoke() (*model.ShowAutopilotAddonInstanceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAutopilotAddonInstanceResponse), nil
	}
}

type ShowAutopilotChartInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAutopilotChartInvoker) Invoke() (*model.ShowAutopilotChartResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAutopilotChartResponse), nil
	}
}

type ShowAutopilotChartValuesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAutopilotChartValuesInvoker) Invoke() (*model.ShowAutopilotChartValuesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAutopilotChartValuesResponse), nil
	}
}

type ShowAutopilotClusterInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAutopilotClusterInvoker) Invoke() (*model.ShowAutopilotClusterResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAutopilotClusterResponse), nil
	}
}

type ShowAutopilotClusterEndpointsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAutopilotClusterEndpointsInvoker) Invoke() (*model.ShowAutopilotClusterEndpointsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAutopilotClusterEndpointsResponse), nil
	}
}

type ShowAutopilotClusterUpgradeInfoInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAutopilotClusterUpgradeInfoInvoker) Invoke() (*model.ShowAutopilotClusterUpgradeInfoResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAutopilotClusterUpgradeInfoResponse), nil
	}
}

type ShowAutopilotJobInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAutopilotJobInvoker) Invoke() (*model.ShowAutopilotJobResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAutopilotJobResponse), nil
	}
}

type ShowAutopilotPreCheckInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAutopilotPreCheckInvoker) Invoke() (*model.ShowAutopilotPreCheckResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAutopilotPreCheckResponse), nil
	}
}

type ShowAutopilotQuotasInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAutopilotQuotasInvoker) Invoke() (*model.ShowAutopilotQuotasResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAutopilotQuotasResponse), nil
	}
}

type ShowAutopilotReleaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAutopilotReleaseInvoker) Invoke() (*model.ShowAutopilotReleaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAutopilotReleaseResponse), nil
	}
}

type ShowAutopilotReleaseHistoryInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAutopilotReleaseHistoryInvoker) Invoke() (*model.ShowAutopilotReleaseHistoryResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAutopilotReleaseHistoryResponse), nil
	}
}

type ShowAutopilotUpgradeClusterTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAutopilotUpgradeClusterTaskInvoker) Invoke() (*model.ShowAutopilotUpgradeClusterTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAutopilotUpgradeClusterTaskResponse), nil
	}
}

type ShowAutopilotUpgradeWorkFlowInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAutopilotUpgradeWorkFlowInvoker) Invoke() (*model.ShowAutopilotUpgradeWorkFlowResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAutopilotUpgradeWorkFlowResponse), nil
	}
}

type ShowAutopilotUserChartsQuotasInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAutopilotUserChartsQuotasInvoker) Invoke() (*model.ShowAutopilotUserChartsQuotasResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAutopilotUserChartsQuotasResponse), nil
	}
}

type UpdateAutopilotAddonInstanceInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateAutopilotAddonInstanceInvoker) Invoke() (*model.UpdateAutopilotAddonInstanceResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateAutopilotAddonInstanceResponse), nil
	}
}

type UpdateAutopilotChartInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateAutopilotChartInvoker) Invoke() (*model.UpdateAutopilotChartResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateAutopilotChartResponse), nil
	}
}

type UpdateAutopilotClusterInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateAutopilotClusterInvoker) Invoke() (*model.UpdateAutopilotClusterResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateAutopilotClusterResponse), nil
	}
}

type UpdateAutopilotClusterEipInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateAutopilotClusterEipInvoker) Invoke() (*model.UpdateAutopilotClusterEipResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateAutopilotClusterEipResponse), nil
	}
}

type UpdateAutopilotReleaseInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateAutopilotReleaseInvoker) Invoke() (*model.UpdateAutopilotReleaseResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateAutopilotReleaseResponse), nil
	}
}

type UpgradeAutopilotClusterInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpgradeAutopilotClusterInvoker) Invoke() (*model.UpgradeAutopilotClusterResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpgradeAutopilotClusterResponse), nil
	}
}

type UpgradeAutopilotWorkFlowUpdateInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpgradeAutopilotWorkFlowUpdateInvoker) Invoke() (*model.UpgradeAutopilotWorkFlowUpdateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpgradeAutopilotWorkFlowUpdateResponse), nil
	}
}

type UploadAutopilotChartInvoker struct {
	*invoker.BaseInvoker
}

func (i *UploadAutopilotChartInvoker) Invoke() (*model.UploadAutopilotChartResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UploadAutopilotChartResponse), nil
	}
}
