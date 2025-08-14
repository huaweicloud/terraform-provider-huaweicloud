---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_kubernetes_clusters_configs"
description: |-
  Use this data source to get HSS container kubernetes clusters configuration information within HuaweiCloud.
---

# huaweicloud_hss_container_kubernetes_clusters_configs

Use this data source to get HSS container kubernetes clusters configuration information within HuaweiCloud.

## Example Usage

```hcl
variable "enterprise_project_id" {}
variable "cluster_id" {}
variable "cluster_name" {}
  

data "huaweicloud_hss_container_kubernetes_clusters_configs" "test" {
  enterprise_project_id = var.enterprise_project_id

  cluster_info_list {
    cluster_id   = var.cluster_id
    cluster_name = var.cluster_name
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_info_list` - (Required, List) Specifies the cluster information list.
  The [cluster_info_list](#cluster_info_list_struct) structure is documented below.

* `cluster_id_list` - (Optional, List) Specifies the cluster ID list.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project that the server belongs to.
  The value **0** indicates the default enterprise project. To query servers in all enterprise projects, set this
  parameter to **all_granted_eps**. If you have only the permission on an enterprise project, you need to transfer the
  enterprise project ID to query the server in the enterprise project. Otherwise, an error is reported due to insufficient
  permission.

  -> An enterprise project can be configured only after the enterprise project function is enabled.

<a name="cluster_info_list_struct"></a>
The `cluster_info_list` block supports:

* `cluster_id` - (Required, String) Specifies the cluster ID.

* `cluster_name` - (Required, String) Specifies the cluster name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data_list` - The list of cluster configurations.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block contains:

* `cluster_id` - The cluster ID.

* `protect_node_num` - The number of protected nodes.

* `protect_interrupt_node_num` - The number of nodes with interrupted protection.

* `protect_degradation_node_num` - The number of nodes with degraded protection.

* `unprotect_node_num` - The number of unprotected nodes.

* `node_total_num` - The total number of nodes.

* `cluster_name` - The cluster name.

* `charging_mode` - The charging mode. The valid values are **on_demand** and **free**.

* `prefer_packet_cycle` - Whether to prefer packet cycle. The value `0` indicates false, and `1` indicates true.

* `protect_type` - The protection type.

* `protect_status` - The protection status. The valid values are:
  + **protecting**: Protecting
  + **part_protect**: Partial protection
  + **creating**: Creating
  + **error_protect**: Protection error
  + **unprotect**: Unprotected
  + **wait_protect**: Waiting for protection

* `cluster_type` - The cluster type.

* `fail_reason` - The failure reason.
