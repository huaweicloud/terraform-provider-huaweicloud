---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_processes"
description: |-
  Use this data source to query DataArts Architecture processes within HuaweiCloud.
---

# huaweicloud_dataarts_architecture_processes

Use this data source to query DataArts Architecture processes within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_architecture_processes" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the processes are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the processes belong.

* `name` - (Optional, String) Specifies the name of processes.

* `parent_id` - (Optional, String) Specifies the parent ID of processes.

* `create_by` - (Optional, String) Specifies the creator of the processes.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `processes` - The list of processes that matched filter parameters.  
  The [processes](#dataarts_architecture_processes) structure is documented below.

<a name="dataarts_architecture_processes"></a>
The `processes` block supports:

* `id` - The ID of the process, in UUID format.

* `name` - The name of the process.

* `name_en` - The English name of the process.

* `description` - The description of the process.

* `owner` - The owner of the process.

* `parent_id` - The parent ID of process.

* `prev_id` - The previous ID of process.

* `next_id` - The next ID of process.

* `qualified_id` - The qualified ID of process.

* `created_at` - The creation time of the process, in RFC3339 format.

* `updated_at` - The latest update time of the process, in RFC3339 format.

* `created_by` - The creator of the process.

* `updated_by` - The last editor of the process.
