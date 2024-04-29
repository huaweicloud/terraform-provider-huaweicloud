---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_customer_gateway"
description: ""
---

# huaweicloud_vpn_customer_gateway

Manages a VPN customer gateway resource within HuaweiCloud.

## Example Usage

### Manages a common VPN customer gateway

```hcl
variable "name" {}
variable "ip" {}

resource "huaweicloud_vpn_customer_gateway" "test" {
  name = var.name
  ip   = var.ip
}
```

### Manages a VPN customer gateway with CA certificate

```hcl
variable "name" {}
variable "ip" {}
variable "certificate_content" {}

resource "huaweicloud_vpn_customer_gateway" "test" {
  name                = var.name
  ip                  = var.ip
  certificate_content = var.certificate_content
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

* `certificate_content` - (Optional, String)  The CA certificate content of the customer gateway.

* `tags` - (Optional, Map) Specifies the tags of the customer gateway.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `serial_number` - Indicates the serial number of the customer gateway certificate.

* `signature_algorithm` - Indicates the signature algorithm of the customer gateway certificate.

* `issuer` - Indicates the issuer of the customer gateway certificate.

* `subject` - Indicates the subject of the customer gateway certificate.

* `expire_time` - Indicates the expire time of the customer gateway certificate.

* `is_updatable` - Indicates whether the customer gateway certificate is updatable.

* `created_at` - The create time.

* `updated_at` - The update time.

## Import

The customer gateway can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_vpn_customer_gateway.test <id>
```
