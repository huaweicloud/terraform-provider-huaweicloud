---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_code_table"
description: ""
---

# huaweicloud_dataarts_architecture_code_table

Manages a DataArts Architecture code table resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "directory_id" {}
variable "name" {}
variable "code" {}

resource "huaweicloud_dataarts_architecture_code_table" "test" {
  workspace_id = var.workspace_id
  name         = var.name
  code         = var.code
  directory_id = var.directory_id

  fields {
    name = "field"
    code = "field_code"
    type = "BIGINT"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew)  Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the ID of DataArts Studio workspace.
  Changing this parameter creates a new resource.

* `name` - (Required, String) Specifies the name of the code table.

* `code` - (Required, String) Specifies the code of the code table.

* `directory_id` - (Required, String) Specifies the directory ID of the code table.

* `fields` - (Required, List) Specifies the fields information of the code table.
  The [fields](#DataArts_Architecture_Code_Table_Fields) structure is documented below.

* `description` - (Optional, String) Specifies the description of the code table.

<a name="DataArts_Architecture_Code_Table_Fields"></a>
The `fields` block supports:

* `name` - (Required, String) Specifies the name of a field.

* `code` - (Required, String) Specifies the code of a field.

* `type` - (Required, String) Specifies the type of a field. Valid values are: **BIGINT**, **BOOLEAN**, **DATE**,
  **DECIMAL**, **DOUBLE**, **STRING**, and **TIMESTAMP**.

* `description` - (Optional, String) Specifies the description of a field.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The time when the code table was created.

* `created_by` - The user who created the code table.

* `directory_path` - The directory path of the code table.

* `fields` - The fields information of the code table.
  The [fields](#DataArts_Architecture_Code_Table_Fields_Attribute) structure is documented below.

* `status` - The status of the code table. Valid values are: **DRAFT**, **PUBLISH_DEVELOPING**,
  **PUBLISHED**, **OFFLINE_DEVELOPING**, **OFFLINE** and **REJECT**.

* `updated_at` - The time when the code table was updated.

<a name="DataArts_Architecture_Code_Table_Fields_Attribute"></a>
The `fields` block supports:

* `id` - The ID of the field.

* `ordinal` - The ordinal of a field.

## Import

The DataArts Architecture code table resource can be imported using the `workspace_id` and `name`, separated by
a slash, e.g.

```bash
$ terraform import huaweicloud_dataarts_architecture_code_table.test <workspace_id>/<name>
```
