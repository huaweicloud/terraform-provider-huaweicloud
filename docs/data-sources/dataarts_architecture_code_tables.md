---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_code_tables"
description: |-
  Use this data source to query DataArts Architecture code tables within HuaweiCloud.
---

# huaweicloud_dataarts_architecture_code_tables

Use this data source to query DataArts Architecture code tables within HuaweiCloud.

## Example Usage

### Query all code tables under a specified workspace

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_architecture_code_tables" "test" {
  workspace_id = var.workspace_id
}
```

### Query the code tables under a specified workspace and belongs a time range

```hcl
variable "workspace_id" {}
variable "begin_time" {
  default = "2026-01-01T00:00:00+08:00"
}
variable "end_time" {
  default = "2026-12-31T23:59:59+08:00"
}

data "huaweicloud_dataarts_architecture_code_tables" "test" {
  workspace_id = var.workspace_id
  begin_time   = var.begin_time
  end_time     = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the code tables are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the code tables belong.

* `name` - (Optional, String) Specifies the Chinese name of the code table to be exactly queried.

* `code` - (Optional, String) Specifies the English name of the code table to be exactly queried.

* `create_by` - (Optional, String) Specifies the creator name of the code table to be queried.

* `directory_id` - (Optional, String) Specifies the directory ID of the code table to be queried.

* `status` - (Optional, String) Specifies the status of the code table to be queried.

* `begin_time` - (Optional, String) Specifies the start time of the code table to be queried, in RFC3339 format.

* `end_time` - (Optional, String) Specifies the end time of the code table to be queried, in RFC3339 format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `code_tables` - The list of code tables that matched filter parameters.  
  The [code_tables](#dataarts_architecture_code_tables_attr) structure is documented below.

<a name="dataarts_architecture_code_tables_attr"></a>
The `code_tables` block supports:

* `id` - The ID of the code table.

* `name` - The name of the code table.

* `code` - The code of the code table.

* `directory_id` - The directory ID of the code table.

* `fields` - The fields information of the code table.  
  The [fields](#dataarts_architecture_code_tables_fields_attr) structure is documented below.

* `description` - The description of the code table.

* `directory_path` - The directory path of the code table.

* `created_by` - The user who created the code table.

* `created_at` - The time when the code table was created.

* `updated_at` - The time when the code table was updated.

* `status` - The status of the code table.
  + **DRAFT**
  + **PUBLISHED**
  + **OFFLINE**
  + **REJECT**

<a name="dataarts_architecture_code_tables_fields_attr"></a>
The `fields` block supports:

* `name` - The name of the field.

* `code` - The code of the field.

* `type` - The type of the field.

* `description` - The description of the field.

* `ordinal` - The ordinal of the field.

* `id` - The ID of the field.
