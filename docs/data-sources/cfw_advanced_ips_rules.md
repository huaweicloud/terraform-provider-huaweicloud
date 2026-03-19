---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_advanced_ips_rules"
description: |-
  Use this data source to get the list of CFW advanced ips rules.
---

# huaweicloud_cfw_advanced_ips_rules

Use this data source to get the list of CFW advanced ips rules.

## Example Usage

```hcl
variable "object_id" {}

data "huaweicloud_cfw_advanced_ips_rules" "test" {
  object_id = var.object_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `object_id` - (Required, String) Specifies the protected object ID. This ID is used to distinguish
  between Internet boundary protection and VPC boundary protection after the cloud firewall is created.
  You can get this value from data source `huaweicloud_cfw_firewalls`.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  For enterprise users, if omitted, all enterprise project will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `advanced_ips_rules` - The list of advanced ips rules.

  The [advanced_ips_rules](#data_advanced_ips_rules_struct) structure is documented below.

<a name="data_advanced_ips_rules_struct"></a>
The `advanced_ips_rules` block supports:

* `action` - The action of the advanced ips rule. The value can be:
  + `0`: Record only
  + `1`: Block session
  + `2`: Block IP

* `ips_rule_id` - The ID of the advanced ips rule.

* `ips_rule_type` - The type of the advanced ips rule. The value can be:
  + `0`: Sensitive directory scan
  + `1`: Reverse shell

* `param` - The parameter of the advanced ips rule.

* `status` - The status of the advanced ips rule. The value can be:
  + `0`: Closed
  + `1`: Opened
