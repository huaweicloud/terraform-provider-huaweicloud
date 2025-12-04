---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_pipes"
description: |-
  Use this data source to query the SecMaster pipes within HuaweiCloud.
---

# huaweicloud_secmaster_pipes

Use this data source to query the SecMaster pipes within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_pipes" "example" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the pipes.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to query pipes.

* `dataspace_id` - (Optional, String) Specifies the dataspace ID to filter pipes.

* `pipe_id` - (Optional, String) Specifies the pipe ID to filter pipes.

* `pipe_name` - (Optional, String) Specifies the pipe name to filter pipes.

* `sort_dir` - (Optional, String) Specifies the sort direction.

* `sort_key` - (Optional, String) Specifies the field to sort the results by.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The list of pipes that match the query criteria.

  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `records` block supports:

* `create_by` - The creator of the pipe.

* `create_time` - The creation timestamp of the pipe.

* `dataspace_id` - The ID of the associated dataspace.

* `dataspace_name` - The name of the associated dataspace.

* `description` - The description of the pipe.

* `domain_id` - The domain ID of the pipe.

* `pipe_id` - The unique identifier of the pipe.

* `pipe_name` - The name of the pipe.

* `pipe_type` - The type of the pipe. Valid values are **system-defined** and **user-defined**.

* `project_id` - The project ID of the pipe.

* `shards` - The number of shards for the pipe.

* `storage_period` - The storage period in days for the pipe. The unit is day.

* `update_by` - The last updater of the pipe.

* `update_time` - The last update timestamp of the pipe.
