---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_rule_known_attack_source"
description: |-
  Manages a WAF rule known attack source resource within HuaweiCloud.
---

# huaweicloud_waf_rule_known_attack_source

Manages a WAF rule known attack source resource within HuaweiCloud.

-> **NOTE:** All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be
used. The known attack source rule resource can be used in Cloud Mode and Dedicated Mode.

## Example Usage

```hcl
variable policy_id {}
variable enterprise_project_id {}

resource "huaweicloud_waf_rule_known_attack_source" "test" {
  policy_id             = var.policy_id
  enterprise_project_id = var.enterprise_project_id
  block_type            = "long_ip_block"
  block_time            = 500
  description           = "test description"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `policy_id` - (Required, String, ForceNew) Specifies the policy ID.

  Changing this parameter will create a new resource.

* `block_type` - (Required, String, ForceNew) Specifies the type of WAF known attack source rule.

  Changing this parameter will create a new resource.

  Valid values are as follows:
  + **long_ip_block**: Long-term IP address blocking.
  + **long_cookie_block**: Long-term Cookie blocking.
  + **long_params_block**: Long-term Params blocking.
  + **short_ip_block**: Short-term IP address blocking.
  + **short_cookie_block**: Short-term Cookie blocking.
  + **short_params_block**: Short-term Params blocking.

* `block_time` - (Required, Int) Specifies the blocking time in seconds.
  + If the prefix of `block_type` is **long**, the value ranges from `301` to `1,800`.
  + If the prefix of `block_type` is **short**, the value ranges from `1` to `300`.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID of WAF known attack
  source rule.

  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description of WAF known attack source rule.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

There are two ways to import WAF rule known attack source state.

* Using `policy_id` and `rule_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_waf_rule_known_attack_source.test <policy_id>/<rule_id>
```

* Using `policy_id`, `rule_id` and `enterprise_project_id`, separated by slashes, e.g.

```bash
$ terraform import huaweicloud_waf_rule_known_attack_source.test <policy_id>/<rule_id>/<enterprise_project_id>
```
