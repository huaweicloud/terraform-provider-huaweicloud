---
subcategory: "Enterprise Switch (ESW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_esw_connection_vport_bind"
description: |-
  Manages an ESW connection vport bind resource within HuaweiCloud.
---

# huaweicloud_esw_connection_vport_bind

Manages an ESW connection vport bind resource within HuaweiCloud.

## Example Usage

```hcl
variable "connection_id" {}
variable "port_id" {}

resource "huaweicloud_esw_connection_vport_bind" "test" {
  connection_id = var.connection_id
  port_id       = var.port_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the ESW connection vport bind resource. If
  omitted, the provider-level region will be used. Changing this creates a new resource.

* `connection_id` - (Required, String, NonUpdatable) Specifies the ID of the connection.

* `port_id` - (Required, String, NonUpdatable) Specifies the ID of the port.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format `<connection_id>/<port_id>`.
