---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_cluster_risk_affected_resources"
description: |-
  Use this data source to get the list of HSS container cluster risk affected resources within HuaweiCloud.
---

# huaweicloud_hss_container_cluster_risk_affected_resources

Use this data source to get the list of HSS container cluster risk affected resources within HuaweiCloud.

## Example Usage

```hcl
variable "risk_id" {}

data "huaweicloud_hss_container_cluster_risk_affected_resources" "test" {
  risk_id = var.risk_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `risk_id` - (Required, String) Specifies the risk ID.

* `cluster_id` - (Optional, String) Specifies the cluster ID.

* `resource_name` - (Optional, String) Specifies the resource name.

* `resource_type` - (Optional, String) Specifies the resource type.

* `namespace` - (Optional, String) Specifies the namespace of the resource.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data_list` - The list of affected resource information.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `resource_name` - The resource name.

* `resource_id` - The resource ID.

* `resource_type` - The resource type.

* `namespace` - The namespace of the resource.

* `hit_rule` - The rule that detected the risk in the resource.

* `hit_path_list` - The list of paths where risks exist in the resource.

* `first_scan_time` - The first scan time.

* `last_scan_time` - The last scan time.
