---
subcategory: "Cloud Phone (CPH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cph_phone_stop"
description: |-
  Manages a CPH phone stop resource within HuaweiCloud.
---

# huaweicloud_cph_phone_stop

Manages a CPH phone stop resource within HuaweiCloud.

## Example Usage

```hcl
variable "phone_id" {}

resource "huaweicloud_cph_phone_stop" "test" {
  phone_id = var.phone_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `phone_id` - (Required, String, NonUpdatable) Specifies the ID of the CPH phone.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
