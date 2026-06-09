---
subcategory: "VPN"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_gateway_certificates"
description: |-
  Use this data source to get the list of VPN gateway certificates within HuaweiCloud.
---

# huaweicloud_vpn_gateway_certificates

Use this data source to get the list of VPN gateway certificates within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_vpn_gateway_certificates" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the certificates.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `certificates` - The list of VPN gateway certificates.
  The [certificates](#vpn_gateway_certificates) structure is documented below.

<a name="vpn_gateway_certificates"></a>
The `certificates` block supports:

* `id` - The VPN gateway certificate ID.

* `name` - The VPN gateway certificate name.

* `project_id` - The project ID of the tenant.

* `vgw_id` - The VPN gateway ID.

* `status` - The status of the gateway certificate.
  The valid values are as follows:
  + **BOUND** - Associated
  + **FAULT** - Association failed
  + **BINDING** - Associating

* `issuer` - The issuer of the SM signature certificate.

* `signature_algorithm` - The signature algorithm of the SM signature certificate.

* `certificate_serial_number` - The serial number of the SM signature certificate.

* `certificate_subject` - The subject of the signature certificate.

* `certificate_expire_time` - The expiration time of the SM signature certificate.

* `certificate_chain_serial_number` - The serial number of the CA certificate.

* `certificate_chain_subject` - The subject of the CA certificate.

* `certificate_chain_expire_time` - The expiration time of the CA certificate.

* `enc_certificate_serial_number` - The serial number of the SM encryption certificate.

* `enc_certificate_subject` - The subject of the encryption certificate.

* `enc_certificate_expire_time` - The expiration time of the SM encryption certificate.

* `created_at` - The creation time.

* `updated_at` - The update time.
