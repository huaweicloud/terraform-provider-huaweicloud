---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_group"
description: |-
  Manages an IAM v5 user group resource within HuaweiCloud.
---

# huaweicloud_identityv5_group

Manages an IAM v5 user group resource within HuaweiCloud.

## Example Usage

```hcl
variable "user_group_name" {}

resource "huaweicloud_identityv5_group" "test" {
  group_name  = var.user_group_name
}
```

## Argument Reference

The following arguments are supported:

* `group_name` - (Required, String) Specifies the name of the user group.  
  The name consists of `1` to `128` characters.  
  Only Chinese characters, letters, digits, spaces, hyphens (-), and
  underscores (_) are allowed.

* `description` - (Optional, String) Specifies the description of the user group.  
  The maximum length is `255` characters, cannot includes these special characters: `@#$%^&*<>\`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also the user group ID.

* `created_at` - The creation time of the user group.

* `urn` - The uniform resource name of the user group.

## Import

The IAM v5 group can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_identityv5_group.test <id>
```
