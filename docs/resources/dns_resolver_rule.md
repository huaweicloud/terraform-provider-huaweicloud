---
subcategory: "Domain Name Service (DNS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dns_resolver_rule"
description: |-
  Manages a DNS resolver rule resource within HuaweiCloud.
---

# huaweicloud_dns_resolver_rule

Manages a DNS resolver rule resource within HuaweiCloud.

## Example Usage

```hcl
variable "endpoint_id" {}
variable "resolver_rule_name" {}
variable "domain_name" {}
variable "ip_address_list" {
  type = list(string)
}

resource "huaweicloud_dns_resolver_rule" "test" {
  endpoint_id = var.endpoint_id
  name        = var.resolver_rule_name
  domain_name = var.domain_name

  dynamic "ip_addresses" {
    for_each = var.ip_address_list
    content {
      ip = ip_addresses.value
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the DNS resolver rule.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `endpoint_id` - (Required, String, ForceNew) Specifies the ID of the DNS endpoint to which the resolver rule belongs.
  Changing this parameter will create a new resource.
  
* `name` - (Required, String) Specifies the resolver rule name.  
  The valid length is limited from `1` to `64`, only Chinese and English characters, digits, underscores (_), hyphens (-)
  and dots (.) are allowed.

* `domain_name` - (Required, String, ForceNew) Specifies the domain name.  
  The maximum length of the domain name is `254` characters.  
  The domain name consists of multiple strings separated by dots (.), and the maximum length of a single string is `63`
  characters. Only Chinese and English characters, digits, and hyphens (-) allowed, and it cannot start or end with a
  hyphen. Changing this parameter will create a new resource.

* `ip_addresses` - (Required, List) Specifies the IP address list of the DNS resolver rule.  
  The [ip_address](#resolver_rule_ip_addresses) structure is documented below.

<a name="resolver_rule_ip_addresses"></a>
The `ip_address` block supports:

* `ip` - (Optional, String) Specifies the IP of the IP address.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also the resolver rule ID.

* `status` - The status of the resolver rule.

* `rule_type` - The rule type of the resolver rule.

* `created_at` - The creation time of the resolver rule.

* `updated_at` - The latest update time of the resolver rule.

* `vpcs` - The list of the VPCs to which the resolver rule is bound.
  The [vpcs](#resolver_rule_associated_vpcs) structure is documented below.

-> The newly added VPCs needs to wait until the next time `terraform refresh` command is executed before that value
   can be refreshed.

<a name="resolver_rule_associated_vpcs"></a>
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
$ terraform import huaweicloud_dns_resolver_rule.test <id>
```
