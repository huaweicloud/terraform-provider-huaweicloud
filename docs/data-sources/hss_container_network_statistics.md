---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_network_statistics"
description: |-
  Use this data source to get the network statistics of HSS container clusters within HuaweiCloud.
---

# huaweicloud_hss_container_network_statistics

Use this data source to get the network statistics of HSS container clusters within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_container_network_statistics" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `protected_cluster_total_num` - The number of unprotected clusters.

* `cluster_total_num` - The total number of clusters.

* `namespace_total_num` - The total number of namespaces.

* `network_policy_total_num` - The total number of network policies.
