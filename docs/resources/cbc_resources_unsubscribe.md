---
subcategory: "Cloud Business Center (CBC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbc_resources_unsubscribe"
description: |-
  Use this resource to unsubscribe resources within HuaweiCloud.
---

# huaweicloud_cbc_resources_unsubscribe

Use this resource to unsubscribe resources within HuaweiCloud.

-> This resource is only a one-time action resource for unsubscribing the specified resources. Deleting this resource
   will not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "resource_ids" {
  type = list(string)
}

resource "huaweicloud_cbc_resources_unsubscribe" "test" {
  resource_ids = var.resource_ids
}
```

## Argument Reference

The following arguments are supported:

* `resource_ids` - (Required, List) Specifies the IDs of the resource to be unsubscribed.
  Supports up to `10` resource IDs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
