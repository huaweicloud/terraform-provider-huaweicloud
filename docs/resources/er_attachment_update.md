---
subcategory: "Enterprise Router (ER)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_er_attachment_update"
description: |-
  Use this resource to update the basic attachment information within HuaweiCloud.
---

# huaweicloud_er_attachment_update

Use this resource to update the basic attachment information within HuaweiCloud.

-> This resource is only a one-time action resource for updating the attachment. Deleting this resource will not restore
   the corresponding attachment configuration, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}
variable "attachment_id" {}
variable "new_attachment_name" {}

resource "huaweicloud_er_attachment_update" "test" {
  instance_id   = var.instance_id
  attachment_id = var.attachment_id
  name          = var.new_attachment_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the shared ER instance.

* `attachment_id` - (Required, String, NonUpdatable) Specifies the ID of the attachment to be accept or reject.

* `name` - (Optional, String) Specifies the new name of the attachment.
  The valid length is limited from `1` to `64` characters, only English letters, Chinese characters, digits,
  underscore (_), hyphens (-) and dots (.) allowed.

* `description` - (Optional, String) Specifies the new name of the attachment.
  The valid length is limited from `1` to `255`, and the angle brackets (< and >) are not allowed.
  The description do not restore with an empty value whether this parameter value is omitted or set empty.

-> At least one of parameter `name` and parameter `description` must be set.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
