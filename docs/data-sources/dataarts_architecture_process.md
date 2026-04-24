---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_process"
description: |-
  Use this data source to get the list of the DataArts Architecture process architectures within HuaweiCloud.
---

# huaweicloud_dataarts_architecture_process

Use this data source to get the list of the DataArts Architecture process architectures within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_architecture_process" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the process architectures are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace where the process architectures are located.

* `name` - (Optional, String) Specifies the name or code of the process architecture for fuzzy matching.

* `parent_id` - (Optional, String) Specifies the parent directory ID.  
  An empty value means all, and "-1" means nodes under the root.

* `create_by` - (Optional, String) Specifies the creator of the process architecture.

* `owner` - (Optional, String) Specifies the owner of the process architecture.

* `begin_time` - (Optional, String) Specifies the left boundary of time filtering, in RFC3339 format.

* `end_time` - (Optional, String) Specifies the right boundary of time filtering, in RFC3339 format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `processes` - The list of process architectures that match the filter parameters.  
  The [processes](#dataarts_architecture_processes) structure is documented below.

<a name="dataarts_architecture_processes"></a>
The `processes` block supports:

* `id` - The ID of the process architecture.

* `name` - The name of the process architecture.

* `name_en` - The English name of the process architecture.

* `description` - The description of the process architecture.

* `guid` - The asset ID corresponding to the process architecture.

* `owner` - The owner of the process architecture.

* `parent_id` - The parent directory ID of the process architecture.

* `prev_id` - The previous node ID of the process architecture.

* `next_id` - The next node ID of the process architecture.

* `qualified_id` - The authentication ID of the process architecture, automatically generated.

* `create_by` - The creator of the process architecture.

* `update_by` - The updater of the process architecture.

* `create_time` - The creation time of the process architecture, in RFC3339 format.

* `update_time` - The update time of the process architecture, in RFC3339 format.

* `bizmetric_num` - The number of business metrics owned by the process architecture.

* `children_num` - The number of child processes owned by the process architecture.

* `children` - The list of child directories under the process architecture, in JSON format.
