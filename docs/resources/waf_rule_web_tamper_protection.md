---
subcategory: "Web Application Firewall (WAF)"
---

# huaweicloud_waf_rule_web_tamper_protection

Manages a WAF web tamper protection rule resource within HuaweiCloud.

-> **NOTE:** All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be
used. The web tamper protection rule resource can be used in Cloud Mode, Dedicated Mode and ELB Mode.

## Example Usage

```hcl
variable enterprise_project_id {}
variable policy_id {}

resource "huaweicloud_waf_rule_web_tamper_protection" "rule_1" {
  policy_id             = var.policy_id
  domain                = "www.your-domain.com"
  path                  = "/payment"
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the WAF web tamper protection rules resource. If
  omitted, the provider-level region will be used. Changing this creates a new rule.

* `policy_id` - (Required, String, ForceNew) Specifies the WAF policy ID. Changing this creates a new rule.

* `domain` - (Required, String, ForceNew) Specifies the domain name. Changing this creates a new rule.

* `path` - (Required, String, ForceNew) Specifies the URL protected by the web tamper protection rule, excluding a
  domain name. Changing this creates a new rule.

* `enterprise_project_id` - (Optional, String, ForceNew) The enterprise project ID of WAF tamper protection rule.
  Changing this parameter will create a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The rule ID in UUID format.

## Import

Web Tamper Protection Rules can be imported using the policy ID and rule ID separated by a slash, e.g.:

```sh
terraform import huaweicloud_waf_rule_web_tamper_protection.rule_1 840c6dfdd5604c1781eea033eae2004f/c6dbc13bb7e74788ae53ecc9254b3ea8
```
