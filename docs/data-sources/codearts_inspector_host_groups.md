---
subcategory: "CodeArts Inspector"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_inspector_host_groups"
description: |-
  Use this data source to get the list of CodeArts inspector host groups.
---

# huaweicloud_codearts_inspector_host_groups

Use this data source to get the list of CodeArts inspector host groups.

## Example Usage

```hcl
data "huaweicloud_codearts_inspector_host_groups" "test" {}
```

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `groups` - Indicates the group list.

  The [groups](#groups_struct) structure is documented below.

<a name="groups_struct"></a>
The `groups` block supports:

* `id` - Indicates the group ID.

* `name` - Indicates the group name.
