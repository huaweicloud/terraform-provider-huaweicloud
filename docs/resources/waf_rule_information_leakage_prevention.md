---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_rule_information_leakage_prevention"
description: |-
  Manages a WAF rule information leakage prevention resource within HuaweiCloud.
---

# huaweicloud_waf_rule_information_leakage_prevention

Manages a WAF rule information leakage prevention resource within HuaweiCloud.

-> **NOTE:** All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be
used. The information leakage prevention rule resource can be used in Cloud Mode and Dedicated Mode.

## Example Usage

```hcl
variable policy_id {}
variable enterprise_project_id {}

resource "huaweicloud_waf_rule_information_leakage_prevention" "test" {
  policy_id             = var.policy_id
  enterprise_project_id = var.enterprise_project_id
  path                  = "/test/path"
  type                  = "sensitive"
  contents              = ["phone", "id_card"]
  protective_action     = "block"
  description           = "test description"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `policy_id` - (Required, String, ForceNew) Specifies the policy ID.

  Changing this parameter will create a new resource.

* `path` - (Required, String) Specifies the path to which the rule applies. The field value needs to be prefixed
  with a slash, such as `/xxx/xxx` or `/xxx/xxx*`. The asterisk (*) refers to using the path name as a prefix.

* `type` - (Required, String) Specifies the type of WAF information leakage prevention rule. The value can be **code**
  for response code or **sensitive** for sensitive information.

* `contents` - (Required, List) Specifies the rule contents.
  + If `type` is set to **code**, valid values are **400**, **401**, **402**, **403**, **404**, **405**, **500**,
  **501**, **502**, **503**, **504** and **507**.
  + If `type` is set to **sensitive**, valid values are **phone**, **id_card** and **email**.

* `protective_action` - (Required, String) Specifies the protective action. Valid values are **log** and **block**.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID of WAF information leakage
  prevention rule. For enterprise users, if omitted, default enterprise project will be used.

  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the rule description.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The rule status.

## Import

There are two ways to import WAF rule information leakage prevention state.

* Using `policy_id` and `rule_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_waf_rule_information_leakage_prevention.test <policy_id>/<rule_id>
```

* Using `policy_id`, `rule_id` and `enterprise_project_id`, separated by slashes, e.g.

```bash
$ terraform import huaweicloud_waf_rule_information_leakage_prevention.test <policy_id>/<rule_id>/<enterprise_project_id>
```
