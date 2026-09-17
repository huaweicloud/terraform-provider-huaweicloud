---
subcategory: "DRS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_object_column_info"
description: |-
  Use this data source to query the object column info (column mapping, column filtering) of a DRS job within HuaweiCloud.
---

# huaweicloud_drs_object_column_info

Use this data source to query the object column info (column mapping, column filtering) of a DRS job within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "job_id" {}

data "huaweicloud_drs_object_column_info" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the object column info.  
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the ID of the DRS job to query object column info.

* `object_id` - (Optional, String) Specifies the object ID.  
  If not specified, database-level objects are queried by default.

* `is_refresh` - (Optional, Bool) Specifies whether to force refresh the query result.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `object_with_column_infos` - The list of the objects related to column info.  
  The [object_with_column_infos](#drs_object_column_info_object_with_column_infos) structure is documented below.

<a name="drs_object_column_info_object_with_column_infos"></a>
The `object_with_column_infos` block supports:

* `id` - The ID of the node.

* `parent_id` - The parent node ID.

* `type` - The node type.

* `name` - The node name.

* `alias_name` - The node alias name.

* `notices` - The prompt messages, such as too many tables under a database.

* `extend_info` - The extended information.

* `is_support_expand` - Whether to support expand query.

* `has_column_info` - Whether the object has column info.

* `is_preset` - Whether the object is preset.

* `token_count` - The token count.

* `is_sent` - Whether the object has been sent to node.

* `sent_alias_name` - The alias name sent to node.
