---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_domain_certificate_associate"
description: |-
  Use this resource to associate SSL certificates to the specified domain within HuaweiCloud.
---

# huaweicloud_apig_domain_certificate_associate

Use this resource to associate SSL certificates to the specified domain within HuaweiCloud.

-> Currently, one certificate can only be associated with one domain. For instance with custom inbound ports,
   the same domain name is bound to a certificate at the same time. Enabling or disabling client verification
   takes effect for different ports of the same domain name.

## Example Usage

```hcl
variable "instance_id" {}
variable "group_id" {}
variable "domain_id" {}
variable "certificate_ids" {}

resource "huaweicloud_apig_domain_certificate_associate" "test" {
  instance_id     = var.instance_id
  group_id        = var.group_id
  domain_id       = var.domain_id
  certificate_ids = var.certificate_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the domain and certificates are located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the dedicated instance to which the domain belongs.

* `group_id` - (Required, String, NonUpdatable) Specifies the ID of the API group to which the domain belongs.

* `domain_id` - (Required, String, NonUpdatable) Specifies the ID of the domain.

* `certificate_ids` - (Required, List, NonUpdatable) Specifies the list of the certificate IDs to associate with the domain.

* `verified_client_certificate_enabled` - (Optional, Bool) Specifies whether to enable verified client certificate.  
  Default to **false**.  
  This parameter only can be set **true** only when the SSL certificate contains a trusted root CA.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID,the format is `<instance_id>/<group_id>/<domain_id>`.

## Import

This resource can be imported using its `id` (consists of `instance_id`, `group_id` and `domain_id`, separated by the
slashes (/)), e.g.

```shell
$ terraform import huaweicloud_apig_domain_certificate_associate.test <instance_id>/<group_id>/<domain_id>
```
