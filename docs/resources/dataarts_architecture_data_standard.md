---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_data_standard"
description: ""
---

# huaweicloud_dataarts_architecture_data_standard

Manages DataArts Architecture data standard resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "directory_id" {}

resource "huaweicloud_dataarts_architecture_data_standard" "test"{
  workspace_id = var.workspace_id
  directory_id = var.directory_id

  values {
    fd_name  = "nameCh"
    fd_value = "test_name"
  }

  values {
    fd_name  = "nameEn"
    fd_value = "test_name"
  }

  values {
    fd_name  = "description"
    fd_value = "this is a test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the data standard is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the ID of the workspace to which the data standard belongs.  
  Changing this parameter will create a new resource.

* `directory_id` - (Required, String) Specifies the directory ID to which the data standard belongs.

* `values` - (Required, List) Specifies the value of data standard attributes.  
  The [values](#dataarts_architecture_data_standard_values) structure is documented below.

<a name="dataarts_architecture_data_standard_values"></a>
The `values` block supports:

* `fd_name` - (Required, String) Specifies the name of the data standard attribute.

* `fd_value` - (Optional, String) Specifies the value of the data standard attribute.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `directory_path` - The directory path of the data standard.

* `status` - The status of the data standard.
  + **DRAFT**
  + **PUBLISHED**
  + **OFFLINE**
  + **REJECT**

* `created_by` - The name of data standard creator.

* `updated_by` - The name of data standard updater.

* `created_at` - The creation time of the data standard, in RFC3339 format.

* `updated_at` - The latest update time of the data standard, in RFC3339 format.

* `values` - The value of data standard attributes.  
  The [values](#dataarts_architecture_data_standard_values_attr) structure is documented below.

* `new_biz` - The biz info of the data standard.  
  The [new_biz](#dataarts_architecture_data_standard_new_biz_attr) structure is documented below.

<a name="dataarts_architecture_data_standard_values_attr"></a>
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

<a name="dataarts_architecture_data_standard_new_biz_attr"></a>
The `new_biz` block supports:

* `id` - The ID of the new biz.

* `biz_type` - The type of the new biz.

* `biz_id` - The ID of the new biz.

* `biz_info` - The info of the new biz.

* `status` - The status of the new biz.

* `biz_version` - The version of the new biz.

* `created_at` - The time when the new biz was created, in RFC3339 format.

* `updated_at` - The time when the new biz was updated, in RFC3339 format.

## Import

The DataArts Architecture data standard can be imported using the `workspace_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dataarts_architecture_data_standard.test <workspace_id>/<id>
```
