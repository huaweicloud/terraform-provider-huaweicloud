---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_rule_blacklist"
description: |-
  Manages a WAF blacklist and whitelist rule resource within HuaweiCloud.
---

# huaweicloud_waf_rule_blacklist

Manages a WAF blacklist and whitelist rule resource within HuaweiCloud.

-> **NOTE:** All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be
used. The blacklist and whitelist rule resource can be used in Cloud Mode and Dedicated Mode.

## Example Usage

### WAF rule blacklist and whitelist with IP address

```hcl
variable "policy_id" {}

resource "huaweicloud_waf_rule_blacklist" "rule" {
  policy_id   = var.policy_id
  ip_address  = "192.168.0.0/24"
  action      = 0
  name        = "test_name"
  description = "test description"
}
```

### WAF rule blacklist and whitelist with address group

```hcl
variable "policy_id" {}
variable "address_group_id" {}
variable "enterprise_project_id" {}

resource "huaweicloud_waf_rule_blacklist" "rule" {
  policy_id             = var.policy_id
  address_group_id      = var.address_group_id
  enterprise_project_id = var.enterprise_project_id
  action                = 1
  name                  = "test_name"
  description           = "test description"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the WAF blacklist and whitelist rule resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `policy_id` - (Required, String, ForceNew) Specifies the WAF policy ID. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the rule name. The value can contain a maximum of `64` characters.
  Only letters, digits, hyphens (-), underscores (_) and periods (.) are allowed.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID of WAF rule blacklist
  and whitelist. For enterprise users, if omitted, default enterprise project will be used.
  Changing this parameter will create a new resource.

* `ip_address` - (Optional, String) Specifies the IP address or range. For example, **192.168.0.125** or **192.168.0.0/24**.
  This parameter is required when `address_group_id` is not specified. The parameter `address_group_id` and `ip_address`
  can not be configured together.

* `address_group_id` - (Optional, String) Specifies the WAF address group ID.
  This parameter is required when `ip_address` is not specified. The parameter `address_group_id` and `ip_address`
  can not be configured together.

* `description` - (Optional, String) Specifies the rule description of the WAF address group.

* `action` - (Optional, Int) Specifies the protective action. Defaults to `0`. The value can be:
  + `0`: block the request.
  + `1`: allow the request.
  + `2`: log the request only.

* `status` - (Optional, Int) Specifies the status of WAF blacklist and whitelist rule.
  Valid values are as follows:
  + **0**: Disabled.
  + **1**: Enabled.

  Defaults to `1`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `address_group_name` - The name of the IP address group.

* `address_group_size` - The number of IP addresses or IP address ranges in the IP address group.

## Import

There are two ways to import WAF rule blacklist state.

* Using `policy_id` and `rule_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_waf_rule_blacklist.test <policy_id>/<rule_id>
```

* Using `policy_id`, `rule_id` and `enterprise_project_id`, separated by slashes, e.g.

```bash
$ terraform import huaweicloud_waf_rule_blacklist.test <policy_id>/<rule_id>/<enterprise_project_id>
```
