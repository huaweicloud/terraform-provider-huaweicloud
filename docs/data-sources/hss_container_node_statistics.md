---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_node_statistics"
description: |-
  Use this data source to get HSS container node protection statistics within HuaweiCloud.
---

# huaweicloud_hss_container_node_statistics

Use this data source to get HSS container node protection statistics within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_container_node_statistics" "test" {
  enterprise_project_id = "0"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

  -> An enterprise project can be configured only after the enterprise project function is enabled.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `unprotected_num` - The number of unprotected servers.

* `protected_num` - The total number of protected nodes.

* `protected_num_on_demand` - The number of nodes protected on demand.

* `protected_num_packet_cycle` - The number of nodes protected by quota.

* `cluster_node_not_installed_num` - The number of uninstalled cluster nodes.

* `not_cluster_node_not_installed_num` - The number of uninstalled non-cluster nodes.
