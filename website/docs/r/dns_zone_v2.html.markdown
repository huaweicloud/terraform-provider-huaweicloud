---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_zone_v2"
sidebar_current: "docs-huaweicloud-resource-dns-zone-v2"
description: |-
  Manages a DNS zone in the HuaweiCloud DNS Service
---

# huaweicloud\_dns\_zone_v2

Manages a DNS zone in the HuaweiCloud DNS Service.

## Example Usage

### Create a public DNS zone

```hcl
resource "huaweicloud_dns_zone_v2" "my_public_zone" {
  name        = "example.com."
  email       = "jdoe@example.com"
  description = "An example zone"
  ttl         = 3000
  zone_type   = "public"
}
```

### Create a private DNS zone

```hcl
resource "huaweicloud_dns_zone_v2" "my_private_zone" {
  name        = "1.example.com."
  email       = "jdoe@example.com"
  description = "An example zone"
  ttl         = 3000
  zone_type   = "private"
  router = [
    {
      router_region = "cn-north-1"
      router_id     = "2c1fe4bd-ebad-44ca-ae9d-e94e63847b75"
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Compute client.
    Keypairs are associated with accounts, but a Compute client is needed to
    create one. If omitted, the `region` argument of the provider is used.
    Changing this creates a new DNS zone. Changing this creates a new DNS zone.

* `name` - (Required) The name of the zone. Note the `.` at the end of the name.
  Changing this creates a new DNS zone.

* `email` - (Optional) The email contact for the zone record.

* `zone_type` - (Optional) The type of zone. Can either be `public` or `private`.
  Changing this creates a new DNS zone.

* `router` - (Optional) Router configuration block which is required if zone_type is private.
  The router structure is documented below.

* `ttl` - (Optional) The time to live (TTL) of the zone.

* `description` - (Optional) A description of the zone.

* `value_specs` - (Optional) Map of additional options. Changing this creates a
  new DNS zone.

The `router` block supports:

* `router_id` - (Required) The router UUID.

* `router_region` - (Required) The region of the router.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `email` - See Argument Reference above.
* `zone_type` - See Argument Reference above.
* `ttl` - See Argument Reference above.
* `description` - See Argument Reference above.
* `masters` - An array of master DNS servers.
* `value_specs` - See Argument Reference above.

## Import

This resource can be imported by specifying the zone ID:

```
$ terraform import huaweicloud_dns_zone_v2.zone_1 <zone_id>
```
