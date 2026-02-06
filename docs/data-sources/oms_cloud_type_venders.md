---
subcategory: "Object Storage Migration Service (OMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_oms_cloud_type_venders"
description: |-
  Use this data source to query OMS supported cloud venders.
---

# huaweicloud_oms_cloud_type_venders

Use this data source to query OMS supported cloud venders.

## Example Usage

```hcl
variable "type" {}

data "huaweicloud_oms_cloud_type_venders" "test" {
  type = var.type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `type` - (Required, String) Specifies the connection end type.
  The value can be **src** (source) or **dst** (destination).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` -  The data source ID.

* `venders` - The list of supported cloud venders.
