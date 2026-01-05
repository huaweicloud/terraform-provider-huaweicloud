---
subcategory: "Bare Metal Server (BMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_bms_instance_restart"
description: |-
  Manages a BMS instance restart resource within HuaweiCloud.
---

# huaweicloud_bms_instance_restart

Manages a BMS instance restart resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_bms_instance_restart" "test" {
  type = "HARD"

  servers {
    id = var.instance_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `type` - (Required, String, NonUpdatable) Specifies the BMS reboot type. Value options:
  + **SOFT**: soft restart (invalid)
  + **HARD**: hard restart (default)

* `servers` - (Required, List, NonUpdatable) Specifies BMS IDs.
  The [servers](#servers_struct) structure is documented below.

<a name="servers_struct"></a>
The `servers` block supports:

* `id` - (Required, String, NonUpdatable) Specifies the BMS ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
