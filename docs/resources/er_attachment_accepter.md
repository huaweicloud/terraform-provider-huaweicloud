---
subcategory: "Enterprise Router (ER)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_er_attachment_accepter"
description: |-
  Use this resource to accept or reject the shared attachment within HuaweiCloud.
---

# huaweicloud_er_attachment_accepter

Use this resource to accept or reject the shared attachment within HuaweiCloud.

-> This resource is only a one-time action resource for operating the attachment. Deleting this resource
   will not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}
variable "attachment_id" {}

resource "huaweicloud_er_attachment_accepter" "test" {
  instance_id   = var.instance_id
  attachment_id = var.attachment_id
  action        = "accept"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the shared ER instance.

* `attachment_id` - (Required, String, NonUpdatable) Specifies the ID of the attachment to be accept or reject.

* `action` - (Required, String, NonUpdatable) Specifies the action type.  
  The valid values are as follows:
  + **accept**
  + **reject**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
