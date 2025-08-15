---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_zone"
description: |-
  Manages a DNS zone resource within HuaweiCloud.
---

# huaweicloud_dns_zone

Manages a DNS zone resource within HuaweiCloud.

## Example Usage

### Create a public DNS zone

```hcl
variable "zone_name" {}
variable "email" {}
variable "description" {}

resource "huaweicloud_dns_zone" "test" {
  name        = var.zone_name
  email       = var.email
  zone_type   = "public"
  ttl         = 3000
  description = var.description
}
```

### Create a private DNS zone

```hcl
variable "zone_name" {}
variable "email" {}
variable "description" {}
variable "router_id" {}

resource "huaweicloud_dns_zone" "test" {
  name        = var.zone_name
  email       = var.email
  zone_type   = "private"
  ttl         = 3000
  description = var.description

  router {
    router_id = var.router_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the zone. Note the `.` at the end of the name.  
  Changing this parameter will create a new resource.

* `email` - (Optional, String) Specifies the email address of the administrator managing the zone.

* `zone_type` - (Optional, String, ForceNew) Specifies the type of zone, defaults to **public**.  
  Changing this parameter will create a new resource.  
  The valid values are as follows:
  + **public**
  + **private**

* `ttl` - (Optional, Int) Specifies the time to live (TTL) of the zone, defaults to `300`.  
  The valid value is range from `1` to `2,147,483,647`.
  
* `description` - (Optional, String) Specifies the description of the zone.  
  A maximum of `255` characters are allowed.

* `router` - (Optional, List) Specifies the list of the router of the zone.
Router configuration block which is required if zone_type is private.
  The [router](#zone_router) structure is documented below.

  -> Before changing this parameter, make sure the zone status is enabled.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID of the zone.  
  Changing this parameter will create a new resource.  
  This parameter is only valid for enterprise users, if omitted, default enterprise project will be used.

* `status` - (Optional, String) Specifies the status of the zone, defaults to **ENABLE**.  
  The valid values are as follows:
  + **ENABLE**
  + **DISABLE**

  -> This is a one-time action.

* `proxy_pattern` - (Optional, String, ForceNew) Specifies the recursive resolution proxy mode for subdomains of
  the private zone.  
  Defaults to **AUTHORITY**.  
  Changing this parameter will create a new resource.  
  The valid values are as follows:
  + **AUTHORITY**: The recursive resolution proxy is disabled for the private zone.
  + **RECURSIVE**: The recursive resolution proxy is enabled for the private zone.
  
  -> 1. This parameter ia available only when the `zone_type` parameter is set to **private**.
     <br>2. If this parameter is set to **RECURSIVE**, but you query subdomains that are not configured in the zone namespace,
     the DNS will recursively resolve the subdomains on the Internet and use the result from authoritative DNS servers.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the zone.

<a name="zone_router"></a>
The `router` block supports:

* `router_id` - (Required, String) Specifies the ID of the associated VPC.

* `router_region` - (Optional, String) Specifies the region of the VPC.

* `dnssec` - (Optional, String) Specifies whether to enable DNSSEC for a public zone.
  The valid values are as follows:
  + **ENABLE**
  + **DISABLE**

  -> Before changing this parameter, make sure the zone status is enabled.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` -  The resource ID, also the zone ID.

* `dnssec_infos` - Indicates the DNSSEC infos.
  The [dnssec_infos](#attrblock--dnssec_infos) structure is documented below.

* `masters` - The list of the masters of the DNS server.

<a name="attrblock--dnssec_infos"></a>
The `dnssec_infos` block supports:

* `digest` - Indicates the digest.

* `digest_algorithm` - Indicates the digest algorithm.

* `digest_type` - Indicates the digest type.

* `ds_record` - Indicates the DS record.

* `flag` - Indicates the flag.

* `key_tag` - Indicates the key tag.

* `ksk_public_key` - Indicates the public key.

* `signature` - Indicates the signature algorithm.

* `signature_type` - Indicates the signature type.

* `created_at` - Indicates the creation time. Format is **yyyy-MM-dd'T'HH:mm:ss.SSS**.

* `updated_at` - Indicates the update time. Format is **yyyy-MM-dd'T'HH:mm:ss.SSS**.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

This resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dns_zone.test <id>
```
