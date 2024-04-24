---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_rules_geolocation_access_control"
description: |-
  Use this data source to get a list of geolocation access control rules.
---

# huaweicloud_waf_rules_geolocation_access_control

Use this data source to get a list of geolocation access control rules.

## Example Usage

```hcl
variable "policy_id" {}
variable "rule_id" {}

data "huaweicloud_waf_rules_geolocation_access_control" "test" {
  policy_id = var.policy_id
  rule_id   = var.rule_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `policy_id` - (Required, String) Specifies the ID of the policy to which the the geolocation access control rules belong.

* `rule_id` - (Optional, String) Specifies the ID of the geolocation access control rule.

* `name` - (Optional, String) Specifies the name of the geolocation access control rule.

* `status` - (Optional, String) Specifies the status of the geolocation access control rule.
  The valid values are as follows:
  + **0**: The geolocation access control rule is disabled.
  + **1**: The geolocation access control rule is active.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the protection policies belong.
  If omitted, will query the geolocation access control rules under the default enterprise project for enterprise users.

* `action` - (Optional, String) Specifies the protective action of the geolocation access control rule.
  The valid values are as follows:
  + **0**: Intercept the request.
  + **1**: Release the request.
  + **2**: Record the request only.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list of the geolocation access control rules.

  The [rules](#rules_struct) structure is documented below.

<a name="rules_struct"></a>
The `rules` block supports:

* `id` - The ID of the geolocation access control rule.

* `policy_id` - The ID of the policy to which the geolocation access control rule belongs.

* `name` - The name of the geolocation access control rule.

* `status` - The status of the geolocation access control rule.

* `geolocation` - The locations that can be configured in the geolocation access control rule.

* `action` - The protective action of the geolocation access control rule.

* `created_at` - The creation time of the geolocation access control rule.
