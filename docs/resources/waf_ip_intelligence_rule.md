---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_ip_intelligence_rule"
description: |-
  Manages a WAF IP intelligence rule resource within HuaweiCloud.
---

# huaweicloud_waf_ip_intelligence_rule

Manages a WAF IP intelligence rule resource within HuaweiCloud.

-> **NOTE:** All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be
  used.

-> For IP intelligence rule configuration, you need to know:
  <br/>1. In cloud mode, the IP intelligence rule is available only in the professional and enterprise editions.
  <br/>2. In dedicated mode, only dedicated instances released in September 2022 and later support IP intelligence
  rule.
  <br/>3. ELB-mode WAF does not support IP intelligence rule.

## Example Usage

```hcl
variable "policy_id" {}
variable "type" {}
variable "tags" {
  type = list(string)
}
variable "rule_name" {}

resource "huaweicloud_waf_ip_intelligence_rule" "test" {
  policy_id = var.policy_id
  type      = var.type
  tags      = var.tags
  name      = var.rule_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this will create a new resource.

* `policy_id` - (Required, String, NonUpdatable) Specifies the policy ID to which the IP intelligence rule belongs.

* `type` - (Required, String) Specifies the IP intelligence rule type.
  Currently, the value only can be **idc**.

* `tags` - (Required, List) Specifies the IP reputation library source.
  The valid value are **Dr.Peng**, **AliCloud**, **VNET**, **Microsoft**, **Google**, **Tencent**, **Amazon**, **HW**
  and **MeiTuan**.

* `name` - (Optional, String) Specifies the IP intelligence rule name.

* `action` - (Optional, List) Specifies the protective action configuration.
  The [action](#rule_action) structure is documented below.

* `policyname` - (Optional, String) Specifies the policy name.
  
* `description` - (Optional, String) Specifies the IP intelligence rule description.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.

<a name="rule_action"></a>
The `action` block supports:

* `category` - (Required, String) Specifies the protective action type.
  The value can be **log** **block** or **pass**.

## Attribute Reference

The following attributes are exported:

* `id` - The resource ID.

* `status` - The IP intelligence rule status.
  The value can be `0` (indicates disabled) or `1` (indicates enabled).

* `policyid` - The policy ID.

* `timestamp` - The creation time of the IP intelligence rule, in millisecond.

## Import

* The resource can be imported using `policy_id` and `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_waf_ip_intelligence_rule.test <policy_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `policyname`, `enterprise_project_id`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_waf_ip_intelligence_rule" "test" {
  ...
  
  lifecycle {
    ignore_changes = [
      policyname,
      enterprise_project_id,
    ]
  }
}
```
