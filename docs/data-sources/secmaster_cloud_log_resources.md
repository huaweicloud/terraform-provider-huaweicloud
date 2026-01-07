---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_cloud_log_resources"
description: |-
  Use this data source to get the list of cloud log resources.
---

# huaweicloud_secmaster_cloud_log_resources

Use this data source to get the list of cloud log resources.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_cloud_log_resources" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `region_id` - (Optional, String) Specifies the region ID.

* `sort_key` - (Optional, String) Specifies sorting field.

* `sort_dir` - (Optional, String) Specifies sorting order.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `datasets` - The list of cloud log resources.

  The [datasets](#datasets_struct) structure is documented below.

* `exist` - Whether the resource exist.

* `workspaces` - The list of workspaces.

<a name="datasets_struct"></a>
The `datasets` block supports:

* `alert` - Whether the alarm is triggered.

* `allow_alert` - Whether to allow an alarm.

* `allow_lts` - Whether long-term storage is allowed.

* `create_time` - The creation time.

* `domain_id` - The account ID.

* `enable` - The enable status.

* `project_id` - The project ID.

* `region` - Whether it is regional-level data.

* `region_id` - The region ID.

* `success` - Whether the operation was successful.

* `total` - The total numbers.

* `update_time` - The update time.

* `workspace_id` - The workspace ID.
