data "huaweicloud_cce_clusters" "test" {
  count = var.cluster_id == "" ? 1 : 0

  name = var.cluster_name
}

data "huaweicloud_cce_addon_template" "test" {
  cluster_id = var.cluster_id != "" ? var.cluster_id : try(data.huaweicloud_cce_clusters.test[0].clusters[0].id, null)
  name       = var.addon_template_name
  # ST.002 Disable
  version    = var.addon_version
}

data "huaweicloud_identity_projects" "test" {
  count = var.project_id == "" ? 1 : 0

  name = var.region_name
  # ST.002 Enable
}

locals {
  original_custom = jsondecode(data.huaweicloud_cce_addon_template.test.spec).parameters.custom

  merged_custom = jsonencode({
    # Compared with the template, the fields that need to be added in custom are: cluster_id and tenant_id
    cluster_id                         = var.cluster_id
    tenant_id                          = var.project_id != "" ? var.project_id : try(data.huaweicloud_identity_projects.test[0].projects[0].id, "")
    # The values of the remaining custom fields are all retained
    annotations                        = local.original_custom.annotations
    coresTotal                         = local.original_custom.coresTotal
    expander                           = local.original_custom.expander
    ignoreDaemonSetsUtilization        = local.original_custom.ignoreDaemonSetsUtilization
    ignoreLocalVolumeNodeAffinity      = local.original_custom.ignoreLocalVolumeNodeAffinity
    initialNodeGroupBackoffDuration    = local.original_custom.initialNodeGroupBackoffDuration
    logLevel                           = local.original_custom.logLevel
    maxEmptyBulkDeleteFlag             = local.original_custom.maxEmptyBulkDeleteFlag
    maxNodeGroupBackoffDuration        = local.original_custom.maxNodeGroupBackoffDuration
    maxNodeGroupBinPackingDuration     = local.original_custom.maxNodeGroupBinPackingDuration
    maxNodeProvisionTime               = local.original_custom.maxNodeProvisionTime
    maxNodesTotal                      = local.original_custom.maxNodesTotal
    memoryTotal                        = local.original_custom.memoryTotal
    multiAZBalance                     = local.original_custom.multiAZBalance
    multiAZEnabled                     = local.original_custom.multiAZEnabled
    newEphemeralVolumesPodScaleUpDelay = local.original_custom.newEphemeralVolumesPodScaleUpDelay
    node_match_expressions             = local.original_custom.node_match_expressions
    podDisruptionBudget                = local.original_custom.podDisruptionBudget
    resetUnNeededWhenScaleUp           = local.original_custom.resetUnNeededWhenScaleUp
    scaleDownDelayAfterAdd             = local.original_custom.scaleDownDelayAfterAdd
    scaleDownDelayAfterDelete          = local.original_custom.scaleDownDelayAfterDelete
    scaleDownDelayAfterFailure         = local.original_custom.scaleDownDelayAfterFailure
    scaleDownEnabled                   = local.original_custom.scaleDownEnabled
    scaleDownUnneededTime              = local.original_custom.scaleDownUnneededTime
    scaleDownUtilizationThreshold      = local.original_custom.scaleDownUtilizationThreshold
    scaleUpCpuUtilizationThreshold     = local.original_custom.scaleUpCpuUtilizationThreshold
    scaleUpMemUtilizationThreshold     = local.original_custom.scaleUpMemUtilizationThreshold
    scaleUpUnscheduledPodEnabled       = local.original_custom.scaleUpUnscheduledPodEnabled
    scaleUpUtilizationEnabled          = local.original_custom.scaleUpUtilizationEnabled
    scanInterval                       = local.original_custom.scanInterval
    skipNodesWithCustomControllerPods  = local.original_custom.skipNodesWithCustomControllerPods
    tolerations                        = local.original_custom.tolerations
    unremovableNodeRecheckTimeout      = local.original_custom.unremovableNodeRecheckTimeout
  })
}

resource "huaweicloud_cce_addon" "test" {
  cluster_id    = var.cluster_id
  template_name = var.addon_template_name
  version       = var.addon_version

  values {
    basic_json  = jsonencode(jsondecode(data.huaweicloud_cce_addon_template.test.spec).basic)
    custom_json = local.merged_custom
    flavor_json = jsonencode(jsondecode(data.huaweicloud_cce_addon_template.test.spec).parameters.flavor1)
  }
}
