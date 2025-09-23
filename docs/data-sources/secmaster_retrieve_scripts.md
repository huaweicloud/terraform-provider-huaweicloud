---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_retrieve_scripts"
description: |-
  Use this data source to get the list of SecMaster retrieve scripts within HuaweiCloud.
---

# huaweicloud_secmaster_retrieve_scripts

Use this data source to get the list of SecMaster retrieve scripts within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_retrieve_scripts" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `table_id` - (Optional, String) Specifies the table ID.

* `script_name` - (Optional, String) Specifies the script name.

* `sort_key` - (Optional, String) Specifies the attribute fields for sorting.

* `sort_dir` - (Optional, String) Specifies the sorting order. Supported values are **ASC** and **DESC**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `records` - The retrieve scripts list.

  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `records` block supports:

* `script_id` - The script ID.

* `project_id` - The project ID.

* `workspace_id` - The workspace ID.

* `script_name` - The script name.

* `table_id` - The table ID.

* `category` - The script classification. The value can be **RETRIEVE** or **ANALYSIS**.

* `directory` - The script directory group name, with a length between `1` and `256` characters.

* `description` - The relevant description information of the script, the length is between `1` and `1,024` characters.

* `script` - The script content, with a length between `1` and `10,240` characters.

* `create_by` - The created by.

* `create_time` - The creation time, millisecond timestamp.

* `update_by` - The last updated by.

* `update_time` - The update time, millisecond timestamp.
