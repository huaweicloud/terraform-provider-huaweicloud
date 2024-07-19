---
subcategory: "Cloud Certificate Manager (CCM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ccm_certificate_import"
description: |-
  Using this resource to import an existing SSL certificate to CCM service within HuaweiCloud, and support pushing the
  SSL certificate to other HuaweiCloud services.
---

# huaweicloud_ccm_certificate_import

Using this resource to import an existing SSL certificate to CCM service within HuaweiCloud, and support pushing the
SSL certificate to other HuaweiCloud services.

## Example Usage

### Load the certificate contents from the local files

```hcl
resource "huaweicloud_ccm_certificate_import" "test" {
  name              = "certificate_name"
  certificate       = file("your_directory/xxx_ca.crt")
  certificate_chain = file("your_directory/xxx_ca_chain.crt")
  private_key       = file("your_directory/xxx_server.key")
}
```

### Write the contents of the certificate into the Terraform script

```hcl
variable "certificate" {}
variable "certificate_chain" {}
variable "private_key" {}

resource "huaweicloud_ccm_certificate_import" "test" {
  name              = "certificate_name"
  certificate       = var.certificate
  certificate_chain = var.certificate_chain
  private_key       = var.private_key
}
```

### Push the SSL certificate to another HuaweiCloud service

```hcl
resource "huaweicloud_ccm_certificate_import" "test" {
  name              = "certificate_name"
  certificate       = file("your_directory/xxx_ca.crt")
  certificate_chain = file("your_directory/xxx_ca_chain.crt")
  private_key       = file("your_directory/xxx_server.key")

  target {
    project = ["la-south-2"]
    service = "ELB"
  }

  target {
    service = "CDN"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to import SSL certificate.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the imported SSL certificate.
  The valid length is limited from `3` to `63`. Only letters, digits, hyphens (-), underscores (_), and periods (.) are allowed.

  Changing this parameter will create a new resource.

* `certificate` - (Required, String, ForceNew) Specifies the content of the SSL certificate, PEM format.
  The certificate content can include intermediate certificates and root certificates.
  If the certificate chain is passed in using the field `certificate_chain`, then the field `certificate` only takes
  the certificate content itself. Using the escape character `\n` or `\r\n` to replace carriage return and line feed.

  Changing this parameter will create a new resource.

* `private_key` - (Required, String, ForceNew) Specifies the private encrypted key of the SSL certificate, PEM format.
  The private key protected by password cannot be uploaded. The carriage return character must be replaced with the
  escape character `\n` or `\r\n`.

  Changing this parameter will create a new resource.

* `certificate_chain` - (Optional, String, ForceNew) Specifies the certificate chain of the SSL certificate.
  The certificate chain can also be passed in through `certificate`. Using the escape character `\n` or `\r\n` to
  replace carriage return and line feed characters.

  Changing this parameter will create a new resource.

* `target` - (Optional, List) Specifies the service information that needs to push the SSL certificate.
The [target](#block_target) structure is documented below.

<a name="block_target"></a>
The `target` block supports:

* `service` - (Required, String) Specifies the service to which the certificate is pushed. The options include `CDN`,`WAF`
  and `ELB`.

* `project` - (Optional, List) Specifies the projects where the service you want to push a certificate to.
  The same certificate can be pushed repeatedly to the same WAF or ELB service in the same `project`, but the CDN service
  can only be pushed once.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also the certificate ID).

* `domain` - The domain name bound to the certificate.

* `domain_count` - The number of domain names that can be bound to the certificate.

* `push_support` - Whether the certificate support push.

* `not_before` - The time when the certificate takes effect. If no valid value is obtained, this parameter is left blank.

* `not_after` - The time when the certificate becomes invalid. If no valid value is obtained, this parameter is left blank.

* `status` - The certificate status. Valid values are as follows:
  + `PAID` - The certificate has been paid, and needs to be applied for from the CA.
  + `ISSUED` - The certificate has been issued.
  + `CHECKING` - The certificate application is under review.
  + `CANCELCHECKING` - The application of certificate cancellation is under review.
  + `UNPASSED` - The certificate application fails.
  + `EXPIRED` - The certificate has expired.
  + `REVOKING` - The application of certificate revocation is under review.
  + `REVOKED` - The certificate has been revoked.
  + `UPLOAD` - The certificate is being hosted.
  + `SUPPLEMENTCHECKING` - The application for the new domain name of the multi-domain certificate is under review.
  + `CANCELSUPPLEMENTING` - The cancellation on additional domain names to be added is being reviewed.

* `authentifications` - The domain ownership verification information.
The [authentifications](#authentifications_struct) structure is documented below.

<a name="authentifications_struct"></a>
The `authentifications` block supports:

* `record_name` - The name of a domain ownership verification value.

* `record_type` - The type of the domain name verification value.

* `record_value` - The domain verification value.

* `domain` - The domain name mapping to the verification value

## Import

The CCM certificate import resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ccm_certificate_import.test <id>
```
