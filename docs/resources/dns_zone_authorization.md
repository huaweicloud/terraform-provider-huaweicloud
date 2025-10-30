---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_zone_authorization"
description: |-
  Using this resource to request a sub-domain authorization request to the owner of the main-domain within HuaweiCloud.
---

# huaweicloud_dns_zone_authorization

Using this resource to request a sub-domain authorization request to the owner of the main-domain within HuaweiCloud.

-> This resource is used to request a sub-domain authorization when creating a sub-domain prompts this following error:
   `domain conflicts with other tenants, you need to add TXT authorization verification`.
   Deleting this resource will not clear the corresponding request record, but will only remove the resource information
   from the tfstate file.

## Example Usage

```hcl
variable "sub_domain_name" {}

resource "huaweicloud_dns_zone_authorization" "test" {
  zone_name = var.sub_domain_name
}
```

## Argument Reference

The following arguments are supported:

* `zone_name` - (Required, String, NonUpdatable) Specifies the name of the sub-domain to be authorized.

  -> The main-domain to which this sub-domain belongs must belong to another user, and the main-domain must use the
     HuaweiCloud DNS server addresses.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the sub-domain authorization.

* `second_level_zone_name` - The second-level domain name to which the sub-domain belongs.

* `record` - The TXT record information.  
  The [record](#dns_zone_txt_record) structure is documented below.

* `status` - The authorization status.
  + **CREATED**: Authorization has been created.
  + **VERIFIED**: Authorization has been verified.

* `created_at` - The creation time of the authorization, in RFC3339 format.

* `updated_at` - The latest update time of the authorization, in RFC3339 format.

<a name="dns_zone_txt_record"></a>
The `record` block supports:

* `host` - The host record of the TXT record.

* `value` - The record value of the TXT record.
