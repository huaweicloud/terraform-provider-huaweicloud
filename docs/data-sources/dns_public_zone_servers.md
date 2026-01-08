---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_public_zone_servers"
description: |-
  Use this data source to query the DNS server addresses of the public zone within HuaweiCloud.
---

# huaweicloud_dns_public_zone_servers

Use this data source to query the DNS server addresses of the public zone within HuaweiCloud.

## Example Usage

```hcl
variable "domain_name" {}

data "huaweicloud_dns_public_zone_servers" "test" {
  domain_name = var.domain_name
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, String) Specifies the domain name of the public zone.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `all_hw_dns` - Whether all servers are HuaweiCloud DNS servers.

* `include_hw_dns` - Whether HuaweiCloud DNS servers are included.

* `dns_servers` - The list of DNS server addresses.

* `expected_dns_servers` - The list of expected DNS server addresses.
