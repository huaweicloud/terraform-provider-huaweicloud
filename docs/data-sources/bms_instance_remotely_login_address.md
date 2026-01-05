---
subcategory: "Bare Metal Server (BMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_bms_instance_remotely_login_address"
description: |-
  Use this data source to get the address for remotely logging in to a BMS.
---

# huaweicloud_bms_instance_remotely_login_address

Use this data source to get the address for remotely logging in to a BMS.

## Example Usage

```hcl
variable "server_id" {}

data "huaweicloud_bms_instance_remotely_login_address" "demo" {
  server_id = var.server_id

  remote_console {
    protocol = "serial"
    type     = "serial"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `server_id` - (Required, String) Specifies the BMS ID.

* `remote_console` - (Required, List) Specifies the action to obtain the address for remotely logging in to a BMS.
  The [remote_console](#remote_console_struct) structure is documented below.

<a name="remote_console_struct"></a>
The `remote_console` block supports:

* `protocol` - (Required, String) Specifies the remote login protocol. Set it to **serial**.

* `type` - (Required, String) Specifies the remote login mode. Set it to **serial**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `remote_console` - Indicates the address for remotely logging in to a BMS.
  The [remote_console](#remote_console_attribute) structure is documented below.

<a name="remote_console_attribute"></a>
The `remote_console` block supports:

* `url` - Indicates the remote login URL.
