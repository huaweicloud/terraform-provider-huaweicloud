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

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String) Specifies the workspace ID of DataArts Architecture. Changing this parameter will
  create a new resource.

* `directory_id` - (Required, String) Specifies the directory ID that the data standard belongs to.

* `values` - (Required, List) Specifies the value of data standard attributes.
The [values](#DataStandard_Value) structure is documented below.

<a name="DataStandard_Value"></a>
The `values` block supports:

* `fd_name` - (Required, String) Specifies the name of the data standard attribute.

* `fd_value` - (Optional, String) Specifies the value of the data standard attribute.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `values` - Indicates the value of data standard attributes.
  The [values](#DataStandard_Value) structure is documented below.

* `directory_path` - Indicates the path of the directory.

* `status` - Indicates the status of the data standard. The value can be: **DRAFT**, **PUBLISH_DEVELOPING**,
  **PUBLISHED**, **OFFLINE_DEVELOPING**, **OFFLINE**, **REJECT**.

* `created_by` - Indicates the name of creator.

* `updated_by` - Indicates the name of updater.

* `created_at` - Indicates the creation time of the data standard.

* `updated_at` - Indicates the latest update time of the data standard.

* `new_biz` - Specifies the biz info of manager.
  The [new_biz](#DataStandard_NewBiz) structure is documented below.

<a name="DataStandard_Value"></a>
The `values` block supports:

* `id` - Indicates the ID of the data standard attribute.

* `fd_id` - Indicates the ID of the data standard attribute definition.

* `directory_id` - Indicates the directory ID that the attribute belongs to.

* `row_id` - Indicates the ID of data standard.

* `status` - Indicates the status of the data standard. The value can be: **DRAFT**, **PUBLISH_DEVELOPING**,
  **PUBLISHED**, **OFFLINE_DEVELOPING**, **OFFLINE**, **REJECT**.

* `created_by` - Indicates the name of creator.

* `updated_by` - Indicates the name of updater.

* `created_at` - Indicates the creation time of the data standard.

* `updated_at` - Indicates the latest update time of the data standard.

<a name="DataStandard_NewBiz"></a>
The `new_biz` block supports:

* `id` - Indicates the ID of the new biz.

* `biz_type` - Indicates the type of the new biz.

* `biz_id` - Indicates the ID of data standard.

* `biz_info` - Indicates the info of the new biz.

* `status` - Indicates the status of the new biz.

* `biz_version` - Indicates the version of the new biz.

* `created_by` - Indicates the creation time of the new biz.

* `updated_at` - Indicates the latest update time of the new biz.

## Import

The DataArts Architecture data standard can be imported using the `workspace_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dataarts_architecture_data_standard.test <workspace_id>/<id>
```
