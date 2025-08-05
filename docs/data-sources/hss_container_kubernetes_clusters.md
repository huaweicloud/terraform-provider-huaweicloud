---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_kubernetes_clusters"
description: |-
  Use this data source to get the list of HSS container kubernetes clusters within HuaweiCloud.
---

# huaweicloud_hss_container_kubernetes_clusters

Use this data source to get the list of HSS container kubernetes clusters within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_container_kubernetes_clusters" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the HSS container kubernetes clusters.
  If omitted, the provider-level region will be used.

* `cluster_name` - (Optional, String) Specifies the cluster name.

* `load_agent_info` - (Optional, Bool) Specifies whether to load agent related information.
  The value can be **true** or **false**. Defaults to **false**.

* `scene` - (Optional, String) Specifies the query scenario type.  
  The valid values are as follows:
  + **cluster_risk**: Cluster risk scenario.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `last_update_time` - The latest update time.

* `total_num` - The total number of clusters.

* `cluster_info_list` - The list of cluster information.
  The [cluster_info_list](#cluster_info_list_struct) structure is documented below.

<a name="cluster_info_list_struct"></a>
The `cluster_info_list` block supports:

* `id` - The ID.

* `cluster_name` - The cluster name.

* `cluster_id` - The cluster ID.

* `cluster_type` - The cluster type.

* `status` - The cluster status.  
  The valid values are as follows:
  + **Available**: Indicating that the cluster is in a normal state.
  + **Unavailable**: Indicating cluster anomaly, manual deletion is required or contact the administrator for deletion.
  + **ScalingUp**: Indicating that the cluster is currently undergoing expansion.
  + **ScalingDown**: Indicating that the cluster is currently undergoing capacity reduction.
  + **Creating**: Indicating that the cluster is currently in the process of being created.
  + **Deleting**: Indicating that the cluster is in the process of being deleted.
  + **Upgrading**: Indicating that the cluster is currently undergoing an upgrade process.
  + **Resizing**: The cluster is currently undergoing specification changes.
  + **RollingBack**: Indicating that the cluster is currently in the process of rolling back.
  + **RollbackFailed**: Indicating a cluster rollback exception, please contact the administrator for a rollback retry.
  + **Empty**: The cluster has no resources.

* `version` - The cluster version.

* `total_nodes_number` - The total number of nodes.

* `active_nodes_number` - The normal number of nodes.

* `creation_timestamp` - Create timestamp.

* `agent_installed_num` - The number of installed agent nodes in the cluster.

* `agent_install_failed_num` - The number of failed installation nodes in the cluster.

* `agent_not_install_num` - The number of nodes without agent installed in the cluster.

* `agent_ds_install_status` - The installation status of agent ds in the cluster. When associating agent related
  information, it is necessary to also consider the `last_operate_time` time.  
  The valid values are as follows:
  + **NotInstall**
  + **Installed**

* `agent_ds_failed_reason` - The reason for operation failure.

* `last_operate_timestamp` - The latest operation timestamp, daemon script installation operation time, when the agent
  is still being installed within a `10` minutes interval.

* `last_scan_time` - The latest scan timestamp of the cluster.

* `sys_vul_num` - The number of system vulnerabilities in the cluster.

* `app_vul_num` - The number of application vulnerabilities in the cluster.

* `emg_vul_num` - The number of emergency vulnerabilities in the cluster.

* `risk_assess_num` - The number of risk assessment questions in the cluster.

* `sec_comp_num` - The number of security and compliance issues in the cluster.
