---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_cluster_asset_statistics"
description: |-
  Use this data source to get HSS cluster asset statistics within HuaweiCloud.
---

# huaweicloud_hss_cluster_asset_statistics

Use this data source to get HSS cluster asset statistics within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_cluster_asset_statistics" "test" {
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

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `cluster_num` - The number of clusters.

* `work_load_num` - The number of workloads.

* `service_num` - The number of services.

* `pod_num` - The number of pods.
