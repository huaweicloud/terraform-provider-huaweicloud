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
variable "id_value" {}

resource "huaweicloud_vpn_customer_gateway" "test" {
  name     = var.name
  id_value = var.id_value
}
```

### Manages a VPN customer gateway with CA certificate

```hcl
variable "name" {}
variable "id_value" {}
variable "certificate_content" {}

resource "huaweicloud_vpn_customer_gateway" "test" {
  name                = var.name
  id_value            = var.id_value
  certificate_content = var.certificate_content
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) The customer gateway name.  
  The valid length is limited from `1` to `64`, only letters, digits, hyphens (-) and underscores (_) are allowed.

* `id_value` - (Required, String, ForceNew) Specifies the identifier of a customer gateway.
  When `id_type` is set to **ip**, the value is an IPv4 address in dotted decimal notation, for example, 192.168.45.7.
  When `id_type` is set to **fqdn**, the value is a string of characters that can contain uppercase letters, lowercase letters,
  digits, and special characters. Spaces and the following special characters are not supported: & < > [ ] \ ?.

  Changing this parameter will create a new resource.

* `id_type` - (Optional, String, ForceNew) Specifies the identifier type of a customer gateway.
  The value can be **ip** or **fqdn**. The default value is **ip**.

* `asn` - (Optional, Int, ForceNew) The BGP ASN number of the customer gateway.
  The value ranges from `1` to `4,294,967,295`, the default value is `65,000`.
  Set this parameter to `0` when `id_type` is set to **fqdn**.

  Changing this parameter will create a new resource.

* `certificate_content` - (Optional, String)  The CA certificate content of the customer gateway.

* `tags` - (Optional, Map) Specifies the tags of the customer gateway.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `certificate_id` - Indicates the ID of the customer gateway certificate.

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

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attribute is `certificate_content`. It is generally recommended
running `terraform plan` after importing the resource. You can then decide if changes should be applied to the instance,
or the resource definition should be updated to align with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_vpn_customer_gateway" "test" {
    ...
  lifecycle {
    ignore_changes = [
      certificate_content,
    ]
  }
}
```
