package v1

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/css/v1/model"
)

type AddIndependentNodeInvoker struct {
	*invoker.BaseInvoker
}

func (i *AddIndependentNodeInvoker) Invoke() (*model.AddIndependentNodeResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AddIndependentNodeResponse), nil
	}
}

type ChangeModeInvoker struct {
	*invoker.BaseInvoker
}

func (i *ChangeModeInvoker) Invoke() (*model.ChangeModeResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ChangeModeResponse), nil
	}
}

type ChangeSecurityGroupInvoker struct {
	*invoker.BaseInvoker
}

func (i *ChangeSecurityGroupInvoker) Invoke() (*model.ChangeSecurityGroupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ChangeSecurityGroupResponse), nil
	}
}

type CreateAiOpsInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateAiOpsInvoker) Invoke() (*model.CreateAiOpsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateAiOpsResponse), nil
	}
}

type CreateAutoCreatePolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateAutoCreatePolicyInvoker) Invoke() (*model.CreateAutoCreatePolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateAutoCreatePolicyResponse), nil
	}
}

type CreateBindPublicInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateBindPublicInvoker) Invoke() (*model.CreateBindPublicResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateBindPublicResponse), nil
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

type CreateClustersTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateClustersTagsInvoker) Invoke() (*model.CreateClustersTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateClustersTagsResponse), nil
	}
}

type CreateElbListenerInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateElbListenerInvoker) Invoke() (*model.CreateElbListenerResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateElbListenerResponse), nil
	}
}

type CreateLoadIkThesaurusInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateLoadIkThesaurusInvoker) Invoke() (*model.CreateLoadIkThesaurusResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateLoadIkThesaurusResponse), nil
	}
}

type CreateLogBackupInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateLogBackupInvoker) Invoke() (*model.CreateLogBackupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateLogBackupResponse), nil
	}
}

type CreateSnapshotInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateSnapshotInvoker) Invoke() (*model.CreateSnapshotResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateSnapshotResponse), nil
	}
}

type DeleteAiOpsInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteAiOpsInvoker) Invoke() (*model.DeleteAiOpsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteAiOpsResponse), nil
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

type DeleteClustersTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteClustersTagsInvoker) Invoke() (*model.DeleteClustersTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteClustersTagsResponse), nil
	}
}

type DeleteIkThesaurusInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteIkThesaurusInvoker) Invoke() (*model.DeleteIkThesaurusResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteIkThesaurusResponse), nil
	}
}

type DeleteSnapshotInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteSnapshotInvoker) Invoke() (*model.DeleteSnapshotResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteSnapshotResponse), nil
	}
}

type DownloadCertInvoker struct {
	*invoker.BaseInvoker
}

func (i *DownloadCertInvoker) Invoke() (*model.DownloadCertResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DownloadCertResponse), nil
	}
}

type EnableOrDisableElbInvoker struct {
	*invoker.BaseInvoker
}

func (i *EnableOrDisableElbInvoker) Invoke() (*model.EnableOrDisableElbResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.EnableOrDisableElbResponse), nil
	}
}

type ListAiOpsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListAiOpsInvoker) Invoke() (*model.ListAiOpsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListAiOpsResponse), nil
	}
}

type ListClustersDetailsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListClustersDetailsInvoker) Invoke() (*model.ListClustersDetailsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListClustersDetailsResponse), nil
	}
}

type ListClustersTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListClustersTagsInvoker) Invoke() (*model.ListClustersTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListClustersTagsResponse), nil
	}
}

type ListElbCertsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListElbCertsInvoker) Invoke() (*model.ListElbCertsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListElbCertsResponse), nil
	}
}

type ListElbsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListElbsInvoker) Invoke() (*model.ListElbsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListElbsResponse), nil
	}
}

type ListFlavorsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListFlavorsInvoker) Invoke() (*model.ListFlavorsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListFlavorsResponse), nil
	}
}

type ListImagesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListImagesInvoker) Invoke() (*model.ListImagesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListImagesResponse), nil
	}
}

type ListLogsJobInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListLogsJobInvoker) Invoke() (*model.ListLogsJobResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListLogsJobResponse), nil
	}
}

type ListSmnTopicsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListSmnTopicsInvoker) Invoke() (*model.ListSmnTopicsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListSmnTopicsResponse), nil
	}
}

type ListSnapshotsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListSnapshotsInvoker) Invoke() (*model.ListSnapshotsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListSnapshotsResponse), nil
	}
}

type ListYmlsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListYmlsInvoker) Invoke() (*model.ListYmlsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListYmlsResponse), nil
	}
}

type ListYmlsJobInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListYmlsJobInvoker) Invoke() (*model.ListYmlsJobResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListYmlsJobResponse), nil
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

type RestoreSnapshotInvoker struct {
	*invoker.BaseInvoker
}

func (i *RestoreSnapshotInvoker) Invoke() (*model.RestoreSnapshotResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RestoreSnapshotResponse), nil
	}
}

type RetryUpgradeTaskInvoker struct {
	*invoker.BaseInvoker
}

func (i *RetryUpgradeTaskInvoker) Invoke() (*model.RetryUpgradeTaskResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.RetryUpgradeTaskResponse), nil
	}
}

type ShowAutoCreatePolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowAutoCreatePolicyInvoker) Invoke() (*model.ShowAutoCreatePolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowAutoCreatePolicyResponse), nil
	}
}

type ShowClusterDetailInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowClusterDetailInvoker) Invoke() (*model.ShowClusterDetailResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowClusterDetailResponse), nil
	}
}

type ShowClusterTagInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowClusterTagInvoker) Invoke() (*model.ShowClusterTagResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowClusterTagResponse), nil
	}
}

type ShowElbDetailInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowElbDetailInvoker) Invoke() (*model.ShowElbDetailResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowElbDetailResponse), nil
	}
}

type ShowGetLogSettingInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowGetLogSettingInvoker) Invoke() (*model.ShowGetLogSettingResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowGetLogSettingResponse), nil
	}
}

type ShowIkThesaurusInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowIkThesaurusInvoker) Invoke() (*model.ShowIkThesaurusResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowIkThesaurusResponse), nil
	}
}

type ShowLogBackupInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowLogBackupInvoker) Invoke() (*model.ShowLogBackupResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowLogBackupResponse), nil
	}
}

type ShowVpcepConnectionInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowVpcepConnectionInvoker) Invoke() (*model.ShowVpcepConnectionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowVpcepConnectionResponse), nil
	}
}

type StartAutoSettingInvoker struct {
	*invoker.BaseInvoker
}

func (i *StartAutoSettingInvoker) Invoke() (*model.StartAutoSettingResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StartAutoSettingResponse), nil
	}
}

type StartLogAutoBackupPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *StartLogAutoBackupPolicyInvoker) Invoke() (*model.StartLogAutoBackupPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StartLogAutoBackupPolicyResponse), nil
	}
}

type StartLogsInvoker struct {
	*invoker.BaseInvoker
}

func (i *StartLogsInvoker) Invoke() (*model.StartLogsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StartLogsResponse), nil
	}
}

type StartPublicWhitelistInvoker struct {
	*invoker.BaseInvoker
}

func (i *StartPublicWhitelistInvoker) Invoke() (*model.StartPublicWhitelistResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StartPublicWhitelistResponse), nil
	}
}

type StartTargetClusterConnectivityTestInvoker struct {
	*invoker.BaseInvoker
}

func (i *StartTargetClusterConnectivityTestInvoker) Invoke() (*model.StartTargetClusterConnectivityTestResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StartTargetClusterConnectivityTestResponse), nil
	}
}

type StartVpecpInvoker struct {
	*invoker.BaseInvoker
}

func (i *StartVpecpInvoker) Invoke() (*model.StartVpecpResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StartVpecpResponse), nil
	}
}

type StopLogAutoBackupPolicyInvoker struct {
	*invoker.BaseInvoker
}

func (i *StopLogAutoBackupPolicyInvoker) Invoke() (*model.StopLogAutoBackupPolicyResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StopLogAutoBackupPolicyResponse), nil
	}
}

type StopLogsInvoker struct {
	*invoker.BaseInvoker
}

func (i *StopLogsInvoker) Invoke() (*model.StopLogsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StopLogsResponse), nil
	}
}

type StopPublicWhitelistInvoker struct {
	*invoker.BaseInvoker
}

func (i *StopPublicWhitelistInvoker) Invoke() (*model.StopPublicWhitelistResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StopPublicWhitelistResponse), nil
	}
}

type StopSnapshotInvoker struct {
	*invoker.BaseInvoker
}

func (i *StopSnapshotInvoker) Invoke() (*model.StopSnapshotResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StopSnapshotResponse), nil
	}
}

type StopVpecpInvoker struct {
	*invoker.BaseInvoker
}

func (i *StopVpecpInvoker) Invoke() (*model.StopVpecpResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StopVpecpResponse), nil
	}
}

type UpdateAzByInstanceTypeInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateAzByInstanceTypeInvoker) Invoke() (*model.UpdateAzByInstanceTypeResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateAzByInstanceTypeResponse), nil
	}
}

type UpdateBatchClustersTagsInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateBatchClustersTagsInvoker) Invoke() (*model.UpdateBatchClustersTagsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateBatchClustersTagsResponse), nil
	}
}

type UpdateClusterNameInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateClusterNameInvoker) Invoke() (*model.UpdateClusterNameResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateClusterNameResponse), nil
	}
}

type UpdateEsListenerInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateEsListenerInvoker) Invoke() (*model.UpdateEsListenerResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateEsListenerResponse), nil
	}
}

type UpdateExtendClusterInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateExtendClusterInvoker) Invoke() (*model.UpdateExtendClusterResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateExtendClusterResponse), nil
	}
}

type UpdateExtendInstanceStorageInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateExtendInstanceStorageInvoker) Invoke() (*model.UpdateExtendInstanceStorageResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateExtendInstanceStorageResponse), nil
	}
}

type UpdateFlavorInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateFlavorInvoker) Invoke() (*model.UpdateFlavorResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateFlavorResponse), nil
	}
}

type UpdateFlavorByTypeInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateFlavorByTypeInvoker) Invoke() (*model.UpdateFlavorByTypeResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateFlavorByTypeResponse), nil
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

type UpdateLogSettingInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateLogSettingInvoker) Invoke() (*model.UpdateLogSettingResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateLogSettingResponse), nil
	}
}

type UpdateOndemandClusterToPeriodInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateOndemandClusterToPeriodInvoker) Invoke() (*model.UpdateOndemandClusterToPeriodResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateOndemandClusterToPeriodResponse), nil
	}
}

type UpdatePublicBandWidthInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdatePublicBandWidthInvoker) Invoke() (*model.UpdatePublicBandWidthResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdatePublicBandWidthResponse), nil
	}
}

type UpdateShrinkClusterInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateShrinkClusterInvoker) Invoke() (*model.UpdateShrinkClusterResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateShrinkClusterResponse), nil
	}
}

type UpdateShrinkNodesInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateShrinkNodesInvoker) Invoke() (*model.UpdateShrinkNodesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateShrinkNodesResponse), nil
	}
}

type UpdateSnapshotSettingInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateSnapshotSettingInvoker) Invoke() (*model.UpdateSnapshotSettingResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateSnapshotSettingResponse), nil
	}
}

type UpdateUnbindPublicInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateUnbindPublicInvoker) Invoke() (*model.UpdateUnbindPublicResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateUnbindPublicResponse), nil
	}
}

type UpdateVpcepConnectionInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateVpcepConnectionInvoker) Invoke() (*model.UpdateVpcepConnectionResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateVpcepConnectionResponse), nil
	}
}

type UpdateVpcepWhitelistInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateVpcepWhitelistInvoker) Invoke() (*model.UpdateVpcepWhitelistResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateVpcepWhitelistResponse), nil
	}
}

type UpdateYmlsInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateYmlsInvoker) Invoke() (*model.UpdateYmlsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateYmlsResponse), nil
	}
}

type UpgradeCoreInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpgradeCoreInvoker) Invoke() (*model.UpgradeCoreResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpgradeCoreResponse), nil
	}
}

type UpgradeDetailInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpgradeDetailInvoker) Invoke() (*model.UpgradeDetailResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpgradeDetailResponse), nil
	}
}

type StartKibanaPublicInvoker struct {
	*invoker.BaseInvoker
}

func (i *StartKibanaPublicInvoker) Invoke() (*model.StartKibanaPublicResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StartKibanaPublicResponse), nil
	}
}

type StopPublicKibanaWhitelistInvoker struct {
	*invoker.BaseInvoker
}

func (i *StopPublicKibanaWhitelistInvoker) Invoke() (*model.StopPublicKibanaWhitelistResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StopPublicKibanaWhitelistResponse), nil
	}
}

type UpdateAlterKibanaInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateAlterKibanaInvoker) Invoke() (*model.UpdateAlterKibanaResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateAlterKibanaResponse), nil
	}
}

type UpdateCloseKibanaInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateCloseKibanaInvoker) Invoke() (*model.UpdateCloseKibanaResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateCloseKibanaResponse), nil
	}
}

type UpdatePublicKibanaWhitelistInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdatePublicKibanaWhitelistInvoker) Invoke() (*model.UpdatePublicKibanaWhitelistResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdatePublicKibanaWhitelistResponse), nil
	}
}

type AddFavoriteInvoker struct {
	*invoker.BaseInvoker
}

func (i *AddFavoriteInvoker) Invoke() (*model.AddFavoriteResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.AddFavoriteResponse), nil
	}
}

type CreateCnfInvoker struct {
	*invoker.BaseInvoker
}

func (i *CreateCnfInvoker) Invoke() (*model.CreateCnfResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.CreateCnfResponse), nil
	}
}

type DeleteConfInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteConfInvoker) Invoke() (*model.DeleteConfResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteConfResponse), nil
	}
}

type DeleteConfigInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteConfigInvoker) Invoke() (*model.DeleteConfigResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteConfigResponse), nil
	}
}

type DeleteTemplateInvoker struct {
	*invoker.BaseInvoker
}

func (i *DeleteTemplateInvoker) Invoke() (*model.DeleteTemplateResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.DeleteTemplateResponse), nil
	}
}

type ListActionsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListActionsInvoker) Invoke() (*model.ListActionsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListActionsResponse), nil
	}
}

type ListCertsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListCertsInvoker) Invoke() (*model.ListCertsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListCertsResponse), nil
	}
}

type ListConfsInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListConfsInvoker) Invoke() (*model.ListConfsResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListConfsResponse), nil
	}
}

type ListPipelinesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListPipelinesInvoker) Invoke() (*model.ListPipelinesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListPipelinesResponse), nil
	}
}

type ListTemplatesInvoker struct {
	*invoker.BaseInvoker
}

func (i *ListTemplatesInvoker) Invoke() (*model.ListTemplatesResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ListTemplatesResponse), nil
	}
}

type ShowGetConfDetailInvoker struct {
	*invoker.BaseInvoker
}

func (i *ShowGetConfDetailInvoker) Invoke() (*model.ShowGetConfDetailResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.ShowGetConfDetailResponse), nil
	}
}

type StartConnectivityTestInvoker struct {
	*invoker.BaseInvoker
}

func (i *StartConnectivityTestInvoker) Invoke() (*model.StartConnectivityTestResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StartConnectivityTestResponse), nil
	}
}

type StartPipelineInvoker struct {
	*invoker.BaseInvoker
}

func (i *StartPipelineInvoker) Invoke() (*model.StartPipelineResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StartPipelineResponse), nil
	}
}

type StopHotPipelineInvoker struct {
	*invoker.BaseInvoker
}

func (i *StopHotPipelineInvoker) Invoke() (*model.StopHotPipelineResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StopHotPipelineResponse), nil
	}
}

type StopPipelineInvoker struct {
	*invoker.BaseInvoker
}

func (i *StopPipelineInvoker) Invoke() (*model.StopPipelineResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.StopPipelineResponse), nil
	}
}

type UpdateCnfInvoker struct {
	*invoker.BaseInvoker
}

func (i *UpdateCnfInvoker) Invoke() (*model.UpdateCnfResponse, error) {
	if result, err := i.BaseInvoker.Invoke(); err != nil {
		return nil, err
	} else {
		return result.(*model.UpdateCnfResponse), nil
	}
}
