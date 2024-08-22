---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_resolver_rule"
description: ""
---

# huaweicloud_dns_resolver_rule

Manages a DNS resolver rule resource within HuaweiCloud.

## Example Usage

```hcl
variable subnet_id {}
variable ip {}
variable domain_name {}

resource "huaweicloud_dns_endpoint" "test" {
  name      = "test"
  direction = "inbound"

  ip_addresses {
    subnet_id = var.subnet_id
    ip        = var.ip
  }
  ip_addresses {
    subnet_id = var.subnet_id
  }
}

resource "huaweicloud_dns_resolver_rule" "test" {
  name        = "test"
  domain_name = var.domain_name
  endpoint_id = huaweicloud_dns_endpoint.test.id
  ip_addresses {
    ip = huaweicloud_dns_endpoint.test.ip_addresses[0].ip
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the DNS resolver rule.
  If omitted, the provider-level region will be used. Changing this parameter will create a new DNS resolver rule.

* `name` - (Required, String) Specifies the resolver rule name.

* `domain_name` - (Required, String, ForceNew) Specifies the domain name. Changing this parameter will create
  a new DNS resolver rule.

* `endpoint_id` - (Required, String, ForceNew) Specifies the DNS endpoint id. Changing this parameter will create
  a new DNS resolver rule.

* `ip_addresses` - (Required, List) Specifies the IP address list of the DNS resolver rule.
  The [ip_address](#Address) structure is documented below.

<a name="Address"></a>
The `ip_address` block supports:

* `ip` - (Optional, String) Specifies the IP of the IP address.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of the resolver rule.

* `rule_type` - The rule type of the resolver rule.

* `created_at` - The created time.

* `updated_at` - The last updated time.

* `vpcs` - The VPC list of the DNS resolver rule.
  The [vpcs](#Dns_vpcs) structure is documented below.

<a name="Dns_vpcs"></a>
The `vpcs` block supports:

* `vpc_id` - The VPC ID.

* `vpc_region` - The region of the VPC.

* `status` - The status of the VPC.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The resolver rule can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dns_resolver_rule.test ff8080828a94313a018bf50d67110c86
```
