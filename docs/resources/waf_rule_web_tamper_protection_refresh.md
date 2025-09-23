---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_rule_web_tamper_protection_refresh"
description: |-
  Manages a resource to update the cache of the web tamper protection rule within HuaweiCloud.
---

# huaweicloud_waf_rule_web_tamper_protection_refresh

Manages a resource to update the cache of the web tamper protection rule within HuaweiCloud.

-> All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be used.

-> The current resource is a one-time resource, and destroying this resource will not change the current status.

## Example Usage

```hcl
variable "policy_id" {}
variable "rule_id" {}

resource "huaweicloud_waf_rule_web_tamper_protection_refresh" "test" {
  policy_id = var.policy_id
  rule_id   = var.rule_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `policy_id` - (Required, String, NonUpdatable) Specifies the ID of the policy.

* `rule_id` - (Required, String, NonUpdatable) Specifies the ID of the web tamper protection rule.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID to which the web
  tamper protection rule belongs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
