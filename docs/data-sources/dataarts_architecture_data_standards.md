---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_data_standards"
description: |-
  Use this data source to query DataArts Architecture data standards within HuaweiCloud.
---

# huaweicloud_dataarts_architecture_data_standards

Use this data source to query DataArts Architecture data standards within HuaweiCloud.

## Example Usage

### Query all data standards under a specified workspace

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_architecture_data_standards" "test" {
  workspace_id = var.workspace_id
}
```

### Query the data standards under a specified workspace and belongs a time range

```hcl
variable "workspace_id" {}
variable "begin_time" {
  default = "2026-01-01T00:00:00+08:00"
}
variable "end_time" {
  default = "2026-12-31T23:59:59+08:00"
}

data "huaweicloud_dataarts_architecture_data_standards" "test" {
  workspace_id = var.workspace_id
  begin_time   = var.begin_time
  end_time     = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the data standards are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the data standards belong.

* `name_ch` - (Optional, String) Specifies the Chinese name of the data standard to be exactly queried.

* `name_en` - (Optional, String) Specifies the English code of the data standard to be exactly queried.

* `directory_id` - (Optional, String) Specifies the directory ID of the data standard to be queried.

* `begin_time` - (Optional, String) Specifies the start time of the data standard to be queried, in RFC3339 format.

* `end_time` - (Optional, String) Specifies the end time of the data standard to be queried, in RFC3339 format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_standards` - The list of data standards that matched filter parameters.  
  The [data_standards](#dataarts_architecture_data_standards_attr) structure is documented below.

<a name="dataarts_architecture_data_standards_attr"></a>
The `data_standards` block supports:

* `id` - The ID of the data standard.

* `directory_id` - The directory ID of the data standard.

* `directory_path` - The directory path of the data standard.

* `values` - The value of data standard attributes.  
  The [values](#dataarts_architecture_data_standards_values_attr) structure is documented below.

* `status` - The status of the data standard.
  + **DRAFT**
  + **PUBLISHED**
  + **OFFLINE**
  + **REJECT**

* `created_by` - The name of data standard creator.

* `updated_by` - The name of data standard updater.

* `created_at` - The creation time of the data standard, in RFC3339 format.

* `updated_at` - The latest update time of the data standard, in RFC3339 format.

* `new_biz` - The biz info of the data standard.  
  The [new_biz](#dataarts_architecture_data_standards_new_biz_attr) structure is documented below.

<a name="dataarts_architecture_data_standards_values_attr"></a>
The `values` block supports:

* `fd_name` - The name of the data standard attribute.

* `fd_value` - The value of the data standard attribute.

* `id` - The ID of the data standard attribute.

* `fd_id` - The ID of the data standard attribute definition.

* `row_id` - The row ID of the data standard attribute.

* `directory_id` - The directory ID to which the attribute belongs.

* `status` - The status of the data standard attribute.
  + **DRAFT**
  + **PUBLISHED**
  + **OFFLINE**
  + **REJECT**

* `created_by` - The name of attribute creator.

* `updated_by` - The name of attribute updater.

* `created_at` - The creation time of the data standard attribute, in RFC3339 format.

* `updated_at` - The latest update time of the data standard attribute, in RFC3339 format.

<a name="dataarts_architecture_data_standards_new_biz_attr"></a>
The `new_biz` block supports:

* `id` - The ID of the new biz.

* `biz_type` - The type of the new biz.

* `biz_id` - The ID of the new biz.

* `biz_info` - The info of the new biz.

* `status` - The status of the new biz.

* `biz_version` - The version of the new biz.

* `created_at` - The time when the new biz was created, in RFC3339 format.

* `updated_at` - The time when the new biz was updated, in RFC3339 format.
