---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_kubernetes_clusters_daemonsets"
description: |-
  Use this data source to get the list of HSS container Kubernetes cluster daemonsets within HuaweiCloud.
---

# huaweicloud_hss_container_kubernetes_clusters_daemonsets

Use this data source to get the list of HSS container Kubernetes cluster daemonsets within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_container_kubernetes_clusters_daemonsets" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `type` - (Optional, String) Specifies the cluster type.  
  The valid values are as follows:
  + **cce**: CCE cluster.
  + **ali**: Alibaba Cloud cluster.
  + **tencent**: Tencent Cloud cluster.
  + **azure**: Microsoft Azure cluster.
  + **aws**: Amazon cluster.
  + **self_built_hw**: Huawei Cloud self-built cluster.
  + **self_built_idc**: IDC self-built cluster.

* `show_cluster_log_status` - (Optional, Bool) Specifies whether the query results display the access status of
  cluster logs. The valid values are **true** and **false**, defaults to **false**.

* `access_status` - (Optional, Bool) Specifies to query based on the access status of the cluster,
  with **true** indicating access and **false** indicating no access.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_num` - The total number of connected clusters.

* `upgradeful_num` - The total number of clusters to be upgraded.

* `err_running_num` - The total number of abnormal clusters running.

* `err_access_num` - The total number of abnormal cluster connections.

* `data_list` - The data list.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `latest_version` - Is it the latest version.

* `agent_version` - The cluster agent version.

* `cluster_name` - The cluster name.

* `cluster_id` - The cluster ID.

* `namespace` - The namespace.

* `cluster_type` - The cluster type.

* `node_num` - The total number of nodes.

* `ds_info` - The daemonset status.

  The [ds_info](#ds_info_struct) structure is documented below.

* `cluster_status` - The cluster status.  
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

* `installed_status` - The Cluster DS installation status.  
  The valid values are as follows:
  + **installing**
  + **install_success**
  + **install_failed**
  + **partically_success**
  + **upgrade_success**
  + **upgrade_failed**
  + **upgrading**
  + **none**

* `access_status` - The cluster ANP access status.  
  The valid values are as follows:
  + **not_connect**
  + **connect_success**
  + **connect_fail**
  + **connect_disruption**

* `combined_status` - The combination state of cluster ANP and DS.  
  The valid values are as follows:
  + **accessing**
  + **access_error**
  + **running**
  + **run_error**
  + **upgrading**
  + **upgrade_error**

* `failed_message` - The reasons for cluster plugin access failure.

* `cluster_log_status` - The access status of cluster logs.  
  The valid values are as follows:
  + **success**
  + **partial_success**

* `invoked_service` - Call the service to identify the CCE free medical examination report.  
  The valid values are as follows:
  + **hss**
  + **cce**

* `registry_info` - The image warehouse information.

  The [registry_info](#registry_info_struct) structure is documented below.

<a name="ds_info_struct"></a>
The `ds_info` block supports:

* `desired_num` - The target quantity.

* `current_num` - The current quantity.

* `ready_num` - The ready quantity.

<a name="registry_info_struct"></a>
The `registry_info` block supports:

* `registry_type` - The image repository type.

* `registry_addr` - The image warehouse address.

* `registry_username` - The image repository username.

* `namespace` - The namespace.
