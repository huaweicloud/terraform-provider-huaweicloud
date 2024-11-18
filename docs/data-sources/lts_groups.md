---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_groups"
description: |-
  Use this data source to get the list of LTS log groups.
---

# huaweicloud_lts_groups

Use this data source to get the list of LTS log groups.

## Example Usage

```hcl
data "huaweicloud_lts_groups" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `groups` - All log groups that match the filter parameters.

  The [groups](#groups_struct) structure is documented below.

<a name="groups_struct"></a>
The `groups` block supports:

* `id` - The log group ID.

* `name` - The log group name.

* `ttl_in_days` - The log expiration time(days).

* `tags` - The key/value pairs to associate with the log group.

* `enterprise_project_id` - The enterprise project ID to which the log group belongs.

* `created_at` - The creation time of the log group, in RFC3339 format.
