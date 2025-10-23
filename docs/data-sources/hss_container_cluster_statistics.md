---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_cluster_statistics"
description: |-
  Use this data source to get the HSS container cluster statistics within HuaweiCloud.
---

# huaweicloud_hss_container_cluster_statistics

Use this data source to get the HSS container cluster statistics within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_container_cluster_statistics" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `risk_cluster_num` - The number of clusters with risks.

* `app_vul_cluster_num` - The number of clusters with application vulnerabilities.

* `unscan_cluster_num` - The number of unscanned clusters.

* `all_cluster_num` - The total number of clusters.
