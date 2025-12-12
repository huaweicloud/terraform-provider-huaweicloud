---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_resource_tags"
description: |-
    Use this data source to query resource tag information.
---

# huaweicloud_identityv5_resource_tags

Use this data source to query resource tag information.

## Example Usage

```hcl
variable "resource_id" {}

data "huaweicloud_identityv5_resource_tags" "test" {
  resource_id   = var.resource_id
  resource_type = "user"
}
```

## Argument Reference

The following arguments are supported:

* `resource_type` - (Required, String) Specifies the resource type, can be `trust agency` or `user`.

* `resource_id` - (Required, String) Specifies the resource id, length is 1 to 64 characters,
  only contains letters, numbers and `-` string.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `tags` - Indicates the custom tag list.
  The [tags](#Identityv5_tags) structure is documented below.

<a name="Identityv5_tags"></a>
The `tags` block supports:

* `tag_key` - Indicates the tag key, can contain letters, numbers, spaces and any combination of symbols
  `_`, `.`, `:`, `=`, `+`, `-`, `@`, but cannot start with space or begin with `_sys_`, length range [1,64].

* `tag_value` - Indicates the tag value, can contain letters, numbers, spaces and any combination of symbols
  `_`, `.`, `:`, `/`, `=`, `+`, `-`, `@`, can be empty string, length range [0,128].
