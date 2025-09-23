---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_change_delete"
description: |-
  Manages a COC change delete resource within HuaweiCloud.
---

# huaweicloud_coc_change_delete

Manages a COC change delete resource within HuaweiCloud.

~> Deleting change delete resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
variable "ticket_id" {}

resource "huaweicloud_coc_change_delete" "test" {
  ticket_type = "change"
  ticket_id   = var.ticket_id
}
```

## Argument Reference

The following arguments are supported:

* `ticket_type` - (Required, String, NonUpdatable) Specifies the type of work order that needs to be operated.
  The value **change** needs to be passed to indicate that the work order to be deleted is a change order.

* `ticket_id` - (Required, String, NonUpdatable) Specifies the work order number that needs to be deleted.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which equals to `ticket_id`.
