---
subcategory: "Resource Formation (RFS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rfs_private_modules"
description: |-
  Use this datasource to get the list of private modules.
---

# huaweicloud_rfs_private_modules

Use this datasource to get the list of private modules.

## Example Usage

```hcl
data "huaweicloud_rfs_private_modules" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `sort_key` - (Optional, String) Specifies the sort field, only supports **create_time**.

* `sort_dir` - (Optional, String) Specifies the ascending or descending order.
  Valid values are **asc** and **desc**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `modules` - The list of private modules. By default, these modules are sorted in ascending order of the creation time.

  The [modules](#modules_struct) structure is documented below.

<a name="modules_struct"></a>
The `modules` block supports:

* `module_name` - The name of the private module.

* `module_id` - The unique ID of the private module.

* `module_description` - The description of the private module.

* `create_time` - The creation time of a private module. It is represented in UTC format (YYYY-MM-DDTHH:mm:ss.SSSZ),
  such as **1970-01-01T00:00:00.000Z**.

* `update_time` - The update time of a private module. It is represented in UTC format (YYYY-MM-DDTHH:mm:ss.SSSZ),
  such as **1970-01-01T00:00:00.000Z**.
