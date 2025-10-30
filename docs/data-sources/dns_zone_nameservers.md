---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_zone_nameservers"
description: |-
  Use this data source to get the list of name servers for a public zone.
---

# huaweicloud_dns_zone_nameservers

Use this data source to get the list of name servers for a public zone.

## Example Usage

```hcl
variable "zone_id" {}

data "huaweicloud_dns_zone_nameservers" "test" {
  zone_id = var.zone_id
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Specifies the ID of the public zone to which the name servers belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `nameservers` - The list of name servers of the public zone.
  The [nameservers](#dns_zone_nameservers_attr) structure is documented below.

<a name="dns_zone_nameservers_attr"></a>
The `nameservers` block supports:

* `hostname` - The host name of the name server.

* `priority` - The priority of the name server.
