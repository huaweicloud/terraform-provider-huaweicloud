---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_resolver_rule_associate"
description: |-
  Manages a DNS resolver rule associate resource within HuaweiCloud.
---

# huaweicloud_dns_resolver_rule_associate

Manages a DNS resolver rule associate resource within HuaweiCloud.

## Example Usage

```hcl
variable "resolver_rule_id" {}
variable "vpc_id" {}

resource "huaweicloud_dns_resolver_rule_associate" "test" {
  resolver_rule_id = var.resolver_rule_id
  vpc_id           = var.vpc_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the DNS resolver rule associate.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `resolver_rule_id` - (Required, String, ForceNew) Specifies the DNS resolver rule ID.  
  Changing this parameter will create a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID to associate.  
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, consisting of the `resolver_rule_id` and the `vpc_id`, separated by a slash.

* `status` - The status of the resolver rule association.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The DNS resolver rule associate resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dns_resolver_rule_associate.test <id>
```
