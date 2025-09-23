---
subcategory: "Cloud Phone (CPH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cph_phone_property"
description: |-
  Manages a CPH phone property resource within HuaweiCloud.
---

# huaweicloud_cph_phone_property

Manages a CPH phone property resource within HuaweiCloud.

## Example Usage

```hcl
variable "phone_id" {}

resource "huaweicloud_cph_phone_property" "test" {
  phones {
    phone_id = var.phone_id
    property = jsonencode({
      "com.cph.mainkeys":0,
      "disable.status.bar":0,
      "ro.permission.changed":0,
      "ro.horizontal.screen":0,
      "ro.install.auto":0,
      "ro.com.cph.sfs_enable":0,
      "ro.product.manufacturer":"Huawei",
      "ro.product.name":"monbox",
      "ro.com.cph.notification_disable":0
    })
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `phones` - (Required, List, NonUpdatable) Specifies the CPH phones.
  The [phones](#cph_phones) structure is documented below.

<a name="cph_phones"></a>
The `phones` block supports:

* `phone_id` - (Required, String, NonUpdatable) Specifies the phone ID.

* `property` - (Required, String, NonUpdatable) Specifies the phone property, the format is json string.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
