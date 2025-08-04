---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_kubernetes_cluster_daemonset"
description: |-
  Manages an HSS container kubernetes cluster daemonset resource within HuaweiCloud.
---

# huaweicloud_hss_container_kubernetes_cluster_daemonset

Manages an HSS container kubernetes cluster daemonset resource within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}
variable "cluster_name" {}
variable "auto_upgrade" {}
variable "runtime_name" {}

resource "huaweicloud_hss_container_kubernetes_cluster_daemonset" "test" {
  cluster_id   = var.cluster_id
  cluster_name = var.cluster_name
  auto_upgrade = var.auto_upgrade

  runtime_info {
    runtime_name = var.runtime_name
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the CCE cluster ID.

* `cluster_name` - (Required, String, NonUpdatable) Specifies the CCE cluster name.

* `auto_upgrade` - (Required, Bool, NonUpdatable) Specifies whether to enable automatic agent upgrade.

* `runtime_info` - (Required, List, NonUpdatable) Specifies the container runtime configuration.
  The [runtime_info](#runtime_info_struct) structure is documented below.

* `schedule_info` - (Optional, List) Specifies the node scheduling information.
  The [schedule_info](#schedule_info_struct) structure is documented below.

* `agent_version` - (Optional, String) Specifies the agent version.

* `invoked_service` - (Optional, String) Specifies the calling service.  
  The valid values are as follows:
  + **hss**
  + **cce**

  Defaults to **hss**.

* `charging_mode` - (Optional, String) Specifies the payment mode.  
  The valid values are as follows:
  + **on_demand**: On-demand.
  + **free_security_check**: Free safety medical examination.

* `cce_protection_type` - (Optional, String) Specifies the CCE protection type.  
  The valid values are as follows:
  + **cluster_level**: Cluster level protection.
  + **node_level**: Node level protection.

* `prefer_packet_cycle` - (Optional, Bool) Specifies whether to prioritize the use of package cycle quotas.

-> The `invoked_service`, `charging_mode`, `cce_protection_type`, and `prefer_packet_cycle` parameters are used in
   the CCE integrated protection call scenario.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the hosts under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

<a name="runtime_info_struct"></a>
The `runtime_info` block supports:

* `runtime_name` - (Required, String, NonUpdatable) Specifies the runtime name.  
  The valid values are as follows:
  + **crio_endpoint**: CRIO.
  + **containerd_endpoint**: Containerd.
  + **docker_endpoint**: Docker.
  + **isulad_endpoint**: Isulad.
  + **podman_endpoint**: Podman.

* `runtime_path` - (Optional, String, NonUpdatable) Specifies the runtime path.  

<a name="schedule_info_struct"></a>
The `schedule_info` block supports:

* `node_selector` - (Optional, List) Specifies the node selector.

* `pod_tolerances` - (Optional, List) Specifies the pod tolerance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID same as `cluster_id`.

* `yaml_content` - The original yaml.

* `node_num` - The total number of nodes.

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

* `ds_info` - The ds status.
  The [ds_info](#ds_info_struct) structure is documented below.

* `installed_status` - The cluster ds installation status.  
  The valid values are as follows:
  + **installing**
  + **install_success**
  + **install_failed**
  + **partically_success**
  + **upgrade_success**
  + **upgrade_failed**
  + **upgrading**
  + **none**

<a name="ds_info_struct"></a>
The `ds_info` block supports:

* `desired_num` - The target number.

* `current_num` - The current quantity.

* `ready_num` - The ready quantity.

## Timeouts

This resource provides the following timeouts configuration options:

* `delete` - Default is 3 minutes.

## Import

The HSS container kubernetes cluster daemonset can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_hss_container_kubernetes_cluster_daemonset.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `cluster_id`, `cluster_name`,
`auto_upgrade`, `agent_version`, `invoked_service`, `charging_mode`, `cce_protection_type`, `prefer_packet_cycle`,
`enterprise_project_id`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition
should be updated to align with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_hss_container_kubernetes_cluster_daemonset" "test" { 
  ...

  lifecycle {
    ignore_changes = [
      cluster_id, cluster_name, auto_upgrade, agent_version, invoked_service, charging_mode, cce_protection_type,
      prefer_packet_cycle, enterprise_project_id,
    ]
  }
}
```
