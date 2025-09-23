---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_rule_data_masking"
description: |-
  Manages a WAF Data Masking Rule resource within HuaweiCloud.
---

# huaweicloud_waf_rule_data_masking

Manages a WAF Data Masking Rule resource within HuaweiCloud.

-> **NOTE:** All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be
used. The data masking rule resource can be used in Cloud Mode and Dedicated Mode.

## Example Usage

```hcl
variable "enterprise_project_id" {}
variable "policy_id" {}

resource "huaweicloud_waf_rule_data_masking" "test" {
  policy_id             = var.policy_id
  enterprise_project_id = var.enterprise_project_id
  path                  = "/login"
  field                 = "params"
  subfield              = "password"
  description           = "test description"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the WAF Data Masking rule resource.
  If omitted, the provider-level region will be used. Changing this setting will create a new rule.

* `policy_id` - (Required, String, ForceNew) Specifies the WAF policy ID. Changing this creates a new rule.

* `path` - (Required, String) Specifies the URL to which the data masking rule applies (exact match by default).

* `field` - (Required, String) Specifies the position where the masked field stored. Valid values are:
  + **params**: The field in the parameter.
  + **header**: The field in the header.
  + **form**: The field in the form.
  + **cookie**: The field in the cookie.

* `subfield` - (Required, String) Specifies the name of the masked field, e.g.: **password**.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID of WAF data masking rule.
  For enterprise users, if omitted, default enterprise project will be used.
  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description of WAF data masking rule.

* `status` - (Optional, Int) Specifies the status of WAF web tamper protection rule.
  Valid values are as follows:
  + `0`: Disabled.
  + `1`: Enabled.

  The default value is `1`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The rule ID in UUID format.

## Import

There are two ways to import WAF rule data masking state.

* Using `policy_id` and `rule_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_waf_rule_data_masking.test <policy_id>/<rule_id>
```

* Using `policy_id`, `rule_id` and `enterprise_project_id`, separated by slashes, e.g.

```bash
$ terraform import huaweicloud_waf_rule_data_masking.test <policy_id>/<rule_id>/<enterprise_project_id>
```
