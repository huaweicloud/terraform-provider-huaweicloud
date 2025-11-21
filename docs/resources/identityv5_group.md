---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_group"
description: |-
  Manages an IAM v5 group resource within HuaweiCloud.
---

# huaweicloud_identityv5_group

Manages an IAM v5 group resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_identityv5_group" "group" {
  group_name  = "test_group_name"
  description = "test description"
}
```

## Argument Reference

The following arguments are supported:

* `group_name` - (Required, String) Specifies the name of the group. The username consists of `1` to `128` characters.
  It can contain only uppercase letters, lowercase letters, digits, spaces, and special characters (-_) and cannot
  start with a digit.

* `description` - (Optional, String) Specifies the description of the group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - User group ID, with a length of 1 to 64 characters, consisting only of letters, numbers, and `-`.

* `created_at` - Indicates the time when the IAM group was created.

* `urn` - Indicates the uniform resource name.

## Import

The IAM v5 group can be imported using the `id`, e.g:

```bash
$ terraform import huaweicloud_identityv5_group.test <id>
```
