---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_resolver_rule_associate"
description: ""
---

# huaweicloud_dns_resolver_rule_associate

Manages a DNS resolver rule associate resource within HuaweiCloud.

## Example Usage

```hcl
variable vpc_id {}
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

resource "huaweicloud_dns_resolver_rule_associate" "test" {
  resolver_rule_id = huaweicloud_dns_resolver_rule.test.id
  vpc_id           = var.vpc_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the DNS resolver rule associate.
  If omitted, the provider-level region will be used. Changing this parameter will create a new DNS resolver rule associate.

* `resolver_rule_id` - (Required, String, ForceNew) Specifies the DNS resolver rule ID.
  Changing this parameter will create a new DNS resolver rule associate.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID.
  Changing this parameter will create a new DNS resolver rule associate.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the `resolver_rule_id` and `vpc_id` separated by a slash.

* `status` - The status of the resolver rule associate.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

DNS resolver rule associate can be imported using the `resolver_rule_id` and `vpc_id` separated by a slash e.g.

```bash
$ terraform import huaweicloud_dns_resolver_rule_associate.test ff8080828b0e8c29018bfb599512069d/46fa7c9d-d047-47d9-b5b7-c8d0c0fccc08
```
