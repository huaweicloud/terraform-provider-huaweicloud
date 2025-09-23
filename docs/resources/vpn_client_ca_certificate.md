---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_client_ca_certificate"
description: |-
  Manages a VPN client CA certificate within HuaweiCloud.
---

# huaweicloud_vpn_client_ca_certificate

Manages a VPN client CA certificate within HuaweiCloud.

## Example Usage

```hcl
variable "vpn_server_id" {}
variable "name" {}
variable "content" {}

resource "huaweicloud_vpn_client_ca_certificate" "test" {
  vpn_server_id = var.vpn_server_id
  name          = var.name
  content       = var.content
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `vpn_server_id` - (Required, String, NonUpdatable) Specifies the VPN server ID.

* `name` - (Required, String) Specifies the name of client CA certificate.

* `content` - (Required, String, NonUpdatable) Specifies the content of client CA certificate.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `expiration_time` - The client CA certificate expiration time.

* `issuer` - The client CA certificate issuer.

* `serial_number` - The client CA certificate serial number.

* `signature_algorithm` - The signature algorithm of the client CA certificate.

* `subject` - The client CA certificate subject.

* `created_at` - The creation time.

* `updated_at` - The update time.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The client CA certificate can be imported using `vpn_server_id` and `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_vpn_client_ca_certificate.test <vpn_server_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attribute is `content`. It is generally recommended running
`terraform plan` after importing the resource. You can then decide if changes should be applied to the client CA
certificate, or the resource definition should be updated to align with the client CA certificate.
Also you can ignore changes as below.

```hcl
resource "huaweicloud_vpn_client_ca_certificate" "test" {
    ...

  lifecycle {
    ignore_changes = [
      content
    ]
  }
}
```
