---
subcategory: "Cloud Application Engine (CAE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cae_certificate"
description: |-
  Manages a certificate resource within HuaweiCloud.
---

# huaweicloud_cae_certificate

Manages a certificate resource within HuaweiCloud.

## Example Usage

```hcl
variable "environment_id" {}
variable "certificate_name" {}
variable "certificate_content" {}
variable "certificate_private_key" {}

resource "huaweicloud_cae_certificate" "test" {
  environment_id = var.environment_id
  name           = var.certificate_name
  crt            = var.certificate_content
  key            = var.certificate_private_key
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `environment_id` - (Required, String, ForceNew) Specifies the ID of the CAE environment.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the certificate.
  Changing this creates a new resource.  
  The maximum length of the name is `64` characters, only lowercase letters, digits, hyphens (-) and dots (.) are
  allowed.  
  The name must start and end with a lowercase letter or a digit.

* `crt` - (Required, String) Specifies the content of the certificate.  
  Base64 format corresponding to PEM encoding.

* `key` - (Required, String) Specifies the private key of the certificate.  
  Base64 format corresponding to PEM encoding.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the ID of the enterprise project to which the
  certificate belongs.  
  Changing this creates a new resource.

-> If the `environment_id` belongs to the non-default enterprise project, this parameter is required and is only valid
   for enterprise users.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The creation time of the certificate, in RFC3339 format.

## Import

The certificate resource can be imported using `environment_id` and `name`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_cae_certificate.test <environment_id>/<name>
```

For the certificate with a non-default enterprise project ID, its enterprise project ID need to be specified
additionanlly when importing. All fields are separated by slashes (/), e.g.

```bash
$ terraform import huaweicloud_cae_certificate.test <environment_id>/<name>/<enterprise_project_id>
```
