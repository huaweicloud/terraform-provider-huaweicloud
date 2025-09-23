---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_certificate_batch_domains_associate"
description: |-
  Use this resource to batch bind domain names to an SSL certificate under the specified dedicated instance
  within HuaweiCloud.
---

# huaweicloud_apig_certificate_batch_domains_associate

Use this resource to batch bind domain names to an SSL certificate under the specified dedicated instance
within HuaweiCloud.

-> 1. If this resource was imported and no changes were deployed before deletion (change(s) must be triggered to
   apply the `verify_enabled_domain_names` or `verify_disabled_domain_names` configured in the script), terraform
   will delete all bound domains for current configured certificate. Otherwise, terraform will only delete the bound
   domains(s) managed by the last change.
   <br/>2. For instance with custom inbound ports, the same domain name is bound to a certificate at the same time.
   Enabling or disabling client verification takes effect for different ports of the same domain name.

## Example Usage

```hcl
variable "certificate_id" {}
variable "instance_id" {}
variable "verify_enabled_domain_names" {
  type = list(string)
}
variable "verify_disabled_domain_names" {
  type = list(string)
}

resource "huaweicloud_apig_certificate_batch_domains_associate" "test" {
  certificate_id = var.certificate_id
  instance_id    = var.instance_id

  verify_enabled_domain_names  = var.verify_enabled_domain_names
  verify_disabled_domain_names = var.verify_disabled_domain_names
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the SSL certificate and domains are located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the instance to which the certificate and
  domains belong.

* `certificate_id` - (Required, String, NonUpdatable) Specifies the ID of the SSL certificate.

* `verify_enabled_domain_names` - (Optional, List) Specifies the domain list to be enabled client certificate
  verification.  
  This parameter is valid only when the SSL certificate enabled trusted root CA.

* `verify_disabled_domain_names` - (Optional, List) Specifies the domain list to be disabled client certificate
  verification.

-> The `verify_enabled_domain_names` and `verify_disabled_domain_names` parameters must be specified at least one.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, the format is `<instance_id>/<certificate_id>`.

## Import

This resource can be imported using its `id` (consists of `instance_id` and `certificate_id`, separated by
a slash (/)), e.g.

```bash
$ terraform import huaweicloud_apig_certificate_batch_domains_associate.test <instance_id>/<certificate_id>
```
