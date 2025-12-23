---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_soc_mappers"
description: |-
  Use this data source to get the list of the mappers.
---

# huaweicloud_secmaster_soc_mappers

Use this data source to get the list of the mappers.

## Example Usage

```hcl
variable "workspace_id" {}
variable "mapping_id" {}

data "huaweicloud_secmaster_soc_mappers" "test" {
  workspace_id = var.workspace_id
  mapping_id   = var.mapping_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `mapping_id` - (Required, String) Specifies the mapping ID.

* `name` - (Optional, String) Specifies the mapping name.

* `has_preprocess_rule` - (Optional, Bool) Specifies whether to configure preprocess rule.

* `start_time` - (Optional, String) Specifies the query start time.

* `end_time` - (Optional, String) Specifies the query end time.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The list of mappers.

  The [data](#mapper_data_struct) structure is documented below.

<a name="mapper_data_struct"></a>
The `data` block supports:

* `id` - The mapper ID.

* `name` - The mapper name.

* `project_id` - The project ID.

* `workspace_id` - The workspace ID.

* `dataclass_id` - The dataclass ID.

* `dataclass_name` - The dataclass name.

* `mapper_type_id` - The mapper type ID.

* `mapping_id` - The mapping ID.

* `create_time` - The creation time.

* `creator_id` - The creator ID.

* `creator_name` - The creator name.

* `update_time` - The update time.

* `modifier_id` - The modifier ID.

* `modifier_name` - The modifier name.
