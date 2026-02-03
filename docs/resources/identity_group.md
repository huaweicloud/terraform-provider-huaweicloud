---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_group"
description: |-
  Manages an IAM user group resource within HuaweiCloud.
---

# huaweicloud_identity_group

Manages an IAM user group resource within HuaweiCloud.

-> You *must* have admin privileges to use this resource.

## Example Usage

```hcl
variable "group_name" {}
variable "group_description" {}

resource "huaweicloud_identity_group" "test" {
  name        = var.group_name
  description = var.group_description
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the name of the group.  
  The valid length is limited from `1` to `128`, only Chinese characters, letters, digits, spaces, hyphens (-) and
  underscores (_) are allowed.

* `description` - (Optional, String) Specifies the description of the group.  
  The valid length is limited from `0` to `255`, cannot includes these special characters: `@#$%^&*<>\`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Import

Groups can be imported using their `id`, e.g.

```bash
$ terraform import huaweicloud_identity_group.test <id>
```
