---
subcategory: "Resource Formation (RFS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rfs_private_hooks"
description: |-
  Use this data source to get the list of RFS private hooks.
---

# huaweicloud_rfs_private_hooks

Use this data source to get the list of RFS private hooks.

## Example Usage

```hcl
data "huaweicloud_rfs_private_hooks" "test" {}
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

* `hooks` - The list of private hooks. By default, these hooks are sorted in descending order of the creation time.

  The [hooks](#hooks_struct) structure is documented below.

<a name="hooks_struct"></a>
The `hooks` block supports:

* `hook_id` - The ID of the private hook.

* `hook_name` - The name of the private hook.

* `hook_description` - The description of the private hook.

* `default_version` - The default version of the private hook.

* `configuration` - The configuration of the hook.

  The [configuration](#configuration_struct) structure is documented below.

* `create_time` - The time when the hook was created, in RFC3339 format (yyyy-mm-ddTHH:MM:SSZ),
  such as **1970-01-01T00:00:00Z**.

* `update_time` - The time when the hook was last updated, in RFC3339 format (yyyy-mm-ddTHH:MM:SSZ),
  such as **1970-01-01T00:00:00Z**.

<a name="configuration_struct"></a>
The `configuration` block supports:

* `target_stacks` - The target stacks.

* `failure_mode` - The failure mode.
