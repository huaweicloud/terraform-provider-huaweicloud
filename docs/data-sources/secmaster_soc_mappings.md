---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_soc_mappings"
description: |-
  Use this data source to get the list of the mappings.
---

# huaweicloud_secmaster_soc_mappings

Use this data source to get the list of the mappings.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_soc_mappings" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `name` - (Optional, String) Specifies the mapping name.

* `status` - (Optional, String) Specifies the mapping status.

* `start_time` - (Optional, String) Specifies the query start time.

* `end_time` - (Optional, String) Specifies the query end time.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The list of catalogues that match the query criteria.

  The [data](#mapping_data_struct) structure is documented below.

<a name="mapping_data_struct"></a>
The `data` block supports:

* `id` - The mapping ID.

* `name` - The mapping name.

* `project_id` - The project ID.

* `workspace_id` - The workspace ID.

* `dataclass_id` - The dataclass ID.

* `dataclass_name` - The dataclass name.

* `classifier_id` - The classifier ID.

* `status` - The mapping status.

* `complete_degree` - The complete degree.

* `instance_num` - The number of the instances.

* `description` - The mapping description.

* `update_time` - The update time.

* `create_time` - The creation time.

* `creator_id` - The creator ID.

* `creator_name` - The creator name.

* `modifier_id` - The modifier ID.

* `modifier_name` - The modifier name.
