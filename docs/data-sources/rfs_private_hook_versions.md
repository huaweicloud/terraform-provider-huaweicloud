---
subcategory: "Resource Formation (RFS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rfs_private_hook_versions"
description: |-
  Use this datasource to get the list of private hook versions.
---

# huaweicloud_rfs_private_hook_versions

Use this datasource to get the list of private hook versions.

## Example Usage

```hcl
variable "hook_name" {}

data "huaweicloud_rfs_private_hook_versions" "test" {
  hook_name = var.hook_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `hook_name` - (Required, String) Specifies the name of the private hook.

* `hook_id` - (Optional, String) Specifies the unique ID of the private hook for strong matching.
  If it does not match the current hook, the API returns an error.

* `sort_key` - (Optional, String) Specifies the sort field, only supports **create_time**.

* `sort_dir` - (Optional, String) Specifies the ascending or descending order.
  Valid values are **asc** and **desc**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `versions` - The list of private hook versions. By default, these versions are sorted in descending order of the
  creation time.

  The [versions](#versions_struct) structure is documented below.

<a name="versions_struct"></a>
The `versions` block supports:

* `hook_name` - The name of the private hook.

* `hook_id` - The unique ID of the private hook.

* `hook_version` - The version number of the private hook.

* `hook_version_description` - The description of the private hook version.

* `create_time` - The creation time of a private hook version. It is represented in UTC
  format (YYYY-MM-DDTHH:mm:ss.SSSZ), such as **1970-01-01T00:00:00.000Z**.
