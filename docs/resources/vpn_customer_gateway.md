---
subcategory: "Virtual Private Network (VPN)"
---

# huaweicloud_vpn_customer_gateway

Manages a VPN customer gateway resource within HuaweiCloud.

## Example Usage

```HCL
variable "name" {}
variable "ip" {}

resource "huaweicloud_vpn_customer_gateway" "test" {
  name = var.name
  ip   = var.ip
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) The customer gateway name.

* `ip` - (Required, String, ForceNew) The IP address of the customer gateway.

  Changing this parameter will create a new resource.

* `route_mode` - (Optional, String, ForceNew) The route mode of the customer gateway. The value can be **static** and **bgp**.
  Defaults to **bgp**.

  Changing this parameter will create a new resource.

* `asn` - (Optional, Int, ForceNew) The BGP ASN number of the customer gateway, only works when the route_mode is
  **bgp**. The value ranges from **1** to **4294967295**, the default value is **65000**.

  Changing this parameter will create a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The create time.

* `updated_at` - The update time.

## Import

The customer gateway can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_vpn_customer_gateway.test 0ce123456a00f2591fabc00385ff1234
```
