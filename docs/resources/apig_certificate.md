---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_certificate"
description: ""
---

# huaweicloud_apig_certificate

Manages an APIG SSL certificate resource within HuaweiCloud.

## Example Usage

### Manages a global SSL certificate

```hcl
variable "certificate_name" {}
variable "certificate_content" {
  type    = string
  default = "'-----BEGIN CERTIFICATE-----THIS IS YOUR CERT CONTENT-----END CERTIFICATE-----'"
}
variable "certificate_private_key" {
  type    = string
  default = "'-----BEGIN PRIVATE KEY-----THIS IS YOUR PRIVATE KEY-----END PRIVATE KEY-----'"
}

resource "huaweicloud_apig_certificate" "test" {
  name        = var.certificate_name
  content     = var.certificate_content
  private_key = var.certificate_private_key
}
```

### Manages a local SSL certificate in a specified dedicated APIG instance

```hcl
variable "certificate_name" {}
variable "certificate_content" {
  type    = string
  default = "'-----BEGIN CERTIFICATE-----THIS IS YOUR CERT CONTENT-----END CERTIFICATE-----'"
}
variable "certificate_private_key" {
  type    = string
  default = "'-----BEGIN PRIVATE KEY-----THIS IS YOUR PRIVATE KEY-----END PRIVATE KEY-----'"
}
variable "dedicated_instance_id" {}

resource "huaweicloud_apig_certificate" "test" {
  name        = var.certificate_name
  content     = var.certificate_content
  private_key = var.certificate_private_key
  type        = "instance"
  instance_id = var.dedicated_instance_id
}
```

### Manages a local SSL certificate (with the ROOT CA certificate)

```hcl
variable "certificate_name" {}
variable "certificate_content" {
  type    = string
  default = "'-----BEGIN CERTIFICATE-----THIS IS YOUR CERT CONTENT-----END CERTIFICATE-----'"
}
variable "certificate_private_key" {
  type    = string
  default = "'-----BEGIN PRIVATE KEY-----THIS IS YOUR PRIVATE KEY-----END PRIVATE KEY-----'"
}
variable "root_ca_certificate_content" {
  type    = string
  default = "'-----BEGIN CERTIFICATE-----THIS IS YOUR CERT CONTENT-----END CERTIFICATE-----'"
}
variable "dedicated_instance_id" {}

resource "huaweicloud_apig_certificate" "test" {
  name            = var.certificate_name
  content         = var.certificate_content
  private_key     = var.certificate_private_key
  trusted_root_ca = var.root_ca_certificate_content
  type            = "instance"
  instance_id     = var.dedicated_instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the certificate is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `name` - (Required, String) Specifies the certificate name.  
  The valid length is limited from `4` to `50`, only Chinese and English letters, digits and underscores (_) are
  allowed. The name must start with a Chinese or English letter.  

* `content` - (Required, String) Specifies the certificate content.

* `private_key` - (Required, String) Specifies the private key of the certificate.

-> Read these documentations to learn [how to apply a private SSL certificate](https://support.huaweicloud.com/intl/en-us/tg-ccm/ccm_01_0025.html).

* `type` - (Optional, String, ForceNew) Specifies the certificate type. The valid values are as follows:
  + **instance**
  + **global**

  Defaults to **global**. Changing this will create a new resource.

* `instance_id` - (Optional, String, ForceNew) Specifies the dedicated instance ID to which the certificate belongs.  
  Required if `type` is **instance**.
  Changing this will create a new resource.

* `trusted_root_ca` - (Optional, String) Specifies the trusted **ROOT CA** certificate.

-> Currently, the ROOT CA parameter only certificates of type `instance` are support.
   Read this documentation to learn [how to purchase a private ROOT CA certificate](https://support.huaweicloud.com/intl/en-us/tg-ccm/ccm_01_0016.html).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The certificate ID.

* `effected_at` - The effective time of the certificate, in RFC3339 format (YYYY-MM-DDThh:mm:ssZ).

* `expires_at` - The expiration time of the certificate, in RFC3339 format (YYYY-MM-DDThh:mm:ssZ).

* `signature_algorithm` - What signature algorithm the certificate uses.

* `sans` - The SAN (Subject Alternative Names) of the certificate.

## Import

Certificates can be imported using their `id`, e.g.

```bash
$ terraform import huaweicloud_apig_certificate.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `content`, `private_key` and `trusted_root_ca`.
It is generally recommended running `terraform plan` after importing a certificate.
You can then decide if changes should be applied to the certificate, or the resource definition should be updated to
align with the certificate. Also you can ignore changes as below.

```hcl
resource "huaweicloud_apig_certificate" "test" {
  ...

  lifecycle {
    ignore_changes = [
      content, private_key, trusted_root_ca,
    ]
  }
}
```
