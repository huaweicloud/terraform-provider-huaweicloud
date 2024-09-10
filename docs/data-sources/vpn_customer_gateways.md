---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_customer_gateways"
description: ""
---

# huaweicloud_vpn_customer_gateways

Use this data source to get a list of VPN customer gateways.

## Example Usage

```hcl
variable "route_mode" {}
variable "name" {}
variable "asn" {}
variable "ip" {}

data "huaweicloud_vpn_customer_gateways" "services" {
  route_mode = var.route_mode
  name       = var.name
  asn        = var.asn
  ip         = var.ip
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the VPN customer gateways.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the customer gateway name.

* `ip` - (Optional, String) Specifies the IP address of the customer gateway.

* `route_mode` - (Optional, String) Specifies the route mode of the customer gateway. The value can be **static** and **bgp**.

* `asn` - (Optional, Int) Specifies the BGP ASN number of the customer gateway, only works when the route_mode is
  **bgp**. The value ranges from `1` to `4,294,967,295`.

* `customer_gateway_id` - (Optional, String) Specifies the customer gateway ID used as the query filter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data resource ID.

* `customer_gateways` - All resource customer gateways that match the filter parameters.
  The [customer_gateways](#customer_Gateways) structure is documented below.

<a name="customer_Gateways"></a>
The `customer_gateways` block supports:

* `id` - Indicates the ID of the customer gateway.

* `name` - Indicates the name of the customer gateway.

* `ip` - Indicates the IP of the customer gateway.

* `route_mode` - Indicates the route mode of the customer gateway.

* `asn` - Indicates the asn of the customer gateway.

* `id_type` - Indicates the id_type of the customer gateway.

* `id_value` - Indicates the id_value of the customer gateway.

* `created_at` - The created time.

* `updated_at` - The last updated time.

* `ca_certificate` - Indicates the ca certificate information of the customer gateway.
  The [ca_certificate](#ca_Certificate) structure is documented below.

<a name="ca_Certificate"></a>
The `ca_certificate` block supports:

* `serial_number` - Indicates the serial number of the customer gateway certificate.

* `signature_algorithm` - Indicates the signature algorithm of the customer gateway certificate.

* `issuer` - Indicates the issuer of the customer gateway certificate.

* `subject` - Indicates the subject of the customer gateway certificate.

* `expire_time` - Indicates the expire time of the customer gateway certificate.

* `is_updatable` - Indicates whether the customer gateway certificate is updatable.
