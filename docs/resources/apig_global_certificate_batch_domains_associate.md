---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_global_certificate_batch_domains_associate"
description: |-
  Use this resource to batch bind domain names to the specified global SSL certificate within HuaweiCloud.
---

# huaweicloud_apig_global_certificate_batch_domains_associate

Use this resource to batch bind domain names to the specified global SSL certificate within HuaweiCloud.

-> 1. If this resource was imported and no changes were deployed before deletion (change must be triggered to
   apply the `verify_disabled_domain_names` configured in the script), terraform will delete all bound domain(s) for
   current configured certificate. Otherwise, terraform will only delete the bound domain(s) managed by the last change.
   <br/>2. For instance with custom inbound ports, the same domain name is bound to a certificate at the same time.
   <br/>3. If you need to enable client certificate verification when binding domain(s), please use the
   `huaweicloud_apig_certificate_batch_domains_associate` resource.

## Example Usage

```hcl
variable "certificate_id" {}
variable "verify_disabled_domain_names" {
  type = list(string)
}

resource "huaweicloud_apig_global_certificate_batch_domains_associate" "test" {
  certificate_id = var.certificate_id

  verify_disabled_domain_names = var.verify_disabled_domain_names
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the global SSL certificate and domains
  are located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `certificate_id` - (Required, String, NonUpdatable) Specifies the ID of the global SSL certificate.

* `verify_disabled_domain_names` - (Required, List) Specifies the domain list to be disabled client certificate
  verification.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `domains` - The domain list associated with the global SSL certificate.  
  The [object](#global_certificate_batch_domains_associate_domains_attr) structure is documented below.
  
<a name="global_certificate_batch_domains_associate_domains_attr"></a>
The `domains` block supports:

* `id` - The ID of the associated domain.

* `url_domain` - The associated domain name.

* `instance_id` - The ID of the dedicated instance to which the domain belongs.

* `status` - The CNAME resolution status of the domain name.

* `min_ssl_version` - The minimum SSL protocol version of the domain.

* `api_group_id` - The ID of the API group to which the domain belongs.

* `api_group_name` - The name of the API group to which the domain belongs.

## Import

This resource can be imported using its `id`, e.g.

```bash
$ terraform import huaweicloud_apig_global_certificate_batch_domains_associate.test <certificate_id>
```
