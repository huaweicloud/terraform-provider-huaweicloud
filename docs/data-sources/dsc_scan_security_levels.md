---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_security_levels"
description: |-
  Use this data source to get the list of DSC sensitive data security levels within HuaweiCloud.
---

# huaweicloud_dsc_security_levels

Use this data source to get the list of DSC sensitive data security levels within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_security_levels" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the security level name for fuzzy query.

* `category` - (Optional, String) Specifies the level source category. Valid values are **BUILT_IN**, **BUILT_IN_COPY**
  and **BUILT_SELF**.

* `is_deleted` - (Optional, String) Specifies whether to query all levels. Defaults to **false**, meaning only enabled
  levels are returned.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `security_levels` - The security level information list.

  The [security_levels](#security_levels_struct) structure is documented below.

<a name="security_levels_struct"></a>
The `security_levels` block supports:

* `level_id` - The security level ID.

* `project_id` - The project ID to which the security level belongs.

* `security_level_name` - The name of the security level.

* `color_number` - The color number of the security level, used for ui display.

* `security_level_desc` - The description of the security level.

* `used_count` - The number of recognition templates that use this security level.

* `category` - The creation type of the level.
  + **BUILT_IN**: Built-in rule.
  + **BUILT_IN_COPY**: Copy of built-in rule.
  + **BUILT_SELF**: User-defined rule.

* `create_time` - The creation timestamp of the security level.

* `update_time` - The last update timestamp of the security level.

* `sort_weight` - The sort weight of the security level.

* `is_deleted` - Whether the security level is disabled.
