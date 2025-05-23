---
subcategory: "SFS Turbo"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sfs_turbos_by_tags"
description: |-
  Use this data source to get list of SFS Turbo file systems by tags.
---

# huaweicloud_sfs_turbos_by_tags

Use this data source to get list of SFS Turbo file systems by tags.

## Example Usage

```hcl
data "huaweicloud_sfs_turbos_by_tags" "test" {
  action = "filter"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `action` - (Required, String) Specifies the operation type of listing fily systems by tags.
  The value can be **filter** or **count**.

* `without_any_tag` - (Optional, Bool) Specifies the resources to be queried contain no tags.
  If this parameter is set to **true**, all resources without specified tags are queried. In this case, the tags field
  is ignored. If this parameter is set to **false** or not specified, it does not take effect, meaning that all
  resources are returned or resources are filtered by `tags` or `matches`.

* `tags` - (Optional, List) Specifies the tags to filter the resources. The `tags` can contain a maximum of `20` keys.
  Each tag key can have a maximum of `20` tag values. The tag value corresponding to each tag key can be an empty array
  but the structure cannot be missing. Each tag key must be unique, and tag values of the same tag must be unique.
  The response returns resources containing all tags in this list. Keys in this list are in the AND relationship and
  values in each key-value structure are in the OR relationship. If no `tags` is specified, all data is returned.
  The [tags](#turbos_tags) structure is documented below.

* `matches` - (Optional, List) Specifies the matches condition to filter the resources.
  The [matches](#turbos_matches) structure is documented below.

  -> Currently, the `matches` size is `1`.

<a name="turbos_tags"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the key of the tags.
  The `key` can contain a maximum of `128` characters and cannot be left blank.

* `values` - (Required, List) Specifies the values list of the tags.
  Each value can contain a maximum of `255` characters. An empty list for values indicates any value.

<a name="turbos_matches"></a>
The `matches` block supports:

* `key` - (Required, String) Specifies the key of the matches.
  Currently, only **resource_name** is supported.

* `value` - (Required, String) Specifies the value of the matches.
  If the `value` ends with `*`, prefix search will be performed. e.g. `sfsturbo*` indicates all resources whose names
  start with **sfsturbo** will be returned.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - The list of the SFS Turbo file systems.
  The [resources](#turbos_resources) structure is documented below.

* `total_count` - The total count of the SFS Turbo file systems.

<a name="turbos_resources"></a>
The `resources` block supports:

* `resource_id` - The ID of the SFS Turbo file systems.

* `resource_name` - The name of the SFS Turbo file systems.

* `resource_detail` - The detail of the SFS Turbo file systems.

* `tags` - The tags list of the SFS Turbo file systems.
  The [tags](#turbos_tags_attr) structure is documented below.

<a name="turbos_tags_attr"></a>
The `tags` block supports:

* `key` - The key of the tags.

* `value` - The value of the tags.
