---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_kubernetes_cluster_protection_enable"
description: |-
  Manage a resource to enable HSS container Kubernetes cluster (CCE) protection within HuaweiCloud.
---

# huaweicloud_hss_container_kubernetes_cluster_protection_enable

Manage a resource to enable HSS container Kubernetes cluster (CCE) protection within HuaweiCloud.

-> This resource is only a one-time action resource for HSS container Kubernetes cluster protection. Deleting this resource
   will not disable the protection, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "cluster_name" {}
variable "cluster_type" {}
variable "cluster_id" {}
variable "charging_mode" {}
variable "cce_protection_type" {}
variable "enterprise_project_id" {}

resource "huaweicloud_hss_container_kubernetes_cluster_protection_enable" "test" {
  cluster_name          = var.cluster_name
  cluster_id            = var.cluster_id
  cluster_type          = var.cluster_type
  charging_mode         = var.charging_mode
  cce_protection_type   = var.cce_protection_type
  enterprise_project_id = var.enterprise_project_id
  prefer_packet_cycle   = false

  # Define your own successful protection status based on your business scenario. Please read the description of the
  # following field `monitor protection statuses` carefully before use.
  monitor_protection_statuses = ["protecting", "part_protect", "unprotect", "wait_protect"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `cluster_name` - (Required, String, NonUpdatable) Specifies the name of the cluster to protect.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the ID of the cluster.

* `cluster_type` - (Optional, String, NonUpdatable) Specifies the type of CCE cluster. Options:
  + **existing**: Existing cluster.
  + **adding**: New cluster.

* `charging_mode` - (Optional, String, NonUpdatable) Specifies the charging mode. Options:
  + **on_demand**: On-demand.
  + **free_security_check**: Free security check.

* `cce_protection_type` - (Optional, String, NonUpdatable) Specifies the CCE protection type. Options:
  + **cluster_level**: Cluster-level protection.
  + **node_level**: Node-level protection.

* `prefer_packet_cycle` - (Optional, Bool, NonUpdatable) Specifies whether to prioritize the use of packet cycle quota.
  Defaults to **false**.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the ID of the enterprise project that the resource
  belongs to. The value **0** indicates the default enterprise project. To query resources in all enterprise projects,
  set this parameter to **all_granted_eps**. If you have only the permission on an enterprise project, you need to
  transfer the enterprise project ID to query the resource in the enterprise project.
  Otherwise, an error is reported due to insufficient permission.

* `monitor_protection_statuses` - (Optional, List, NonUpdatable) Specifies the protection statuses to monitor.
  Valid values are **protecting**, **part_protect**, **creating**, **unprotect**, and **wait_protect**.
  The **protecting** status will be monitored by default if this field is not configured.

  -> For example, if you configure the value of this field to **["protecting", "wait_protect"]**, the resource will rotate
  the cluster protection configuration information.
  <br/>1. If the protection status is **protecting** or **wait_protect**, the resource creation is successful.
  <br/>2. If the protection status is **error_protect**, the resource creation fails.
  <br/>3. Other protection states will continue to rotate until the timeout.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID same as `cluster_id`.

* `protect_status` - The protection status of the cluster. Represents the status value matched by the
  parameter `monitor_protection_statuses`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
