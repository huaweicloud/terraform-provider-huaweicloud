---
subcategory: "Tag Management Service (TMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_tms_associated_resource_status"
description: |-
  Manages the TMS associated resource tag inheritance status within HuaweiCloud.
---

# huaweicloud_tms_associated_resource_status

Manages the TMS associated resource tag inheritance status within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_tms_associated_resource_status" "test" {
  status = "enable"
}
```

## Argument Reference

The following arguments are supported:

* `status` - (Required, String) Specifies whether the associated resource tag inheritance is enabled.
  The valid values are as follows:
  + **enable**
  + **disable**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The associated resource status can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_tms_associated_resource_status.test enable
```