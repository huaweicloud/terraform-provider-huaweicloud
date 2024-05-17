---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_line_groups"
description: |-
  Use this data source to get the list of DNS line groups.
---

# huaweicloud_dns_line_groups

Use this data source to get the list of DNS line groups.

## Example Usage

```hcl
variable "line_group_name" {}

data "huaweicloud_dns_line_groups" "test" {
  name = var.line_group_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `line_id` - (Optional, String) Specifies the ID of the line group. Fuzzy search is supported.

* `name` - (Optional, String) Specifies the name of the line group. Fuzzy search is supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `groups` - The list of the line groups.

  The [groups](#groups_struct) structure is documented below.

<a name="groups_struct"></a>
The `groups` block supports:

* `id` - The ID of the line group.

* `name` - The name of the line group.

* `lines` - The list of the resolution line IDs corresponding to line group.

* `description` - The description of the line group.

* `status` - The current status of the line group.
  The valid values are **ACTIVE**, **ERROR**, **FREEZE** and **DISABLE**.

* `created_at` - The creation time of the line group.

* `updated_at` - The latest update time of the line group.
