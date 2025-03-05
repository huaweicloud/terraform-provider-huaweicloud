---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_l7rules"
description: |-
  Use this data source to get the list of ELB L7 rules.
---

# huaweicloud_elb_l7rules

Use this data source to get the list of ELB L7 rules.

## Example Usage

```hcl
variable "compare_type" {}

data "huaweicloud_elb_l7rules" "test" {
  compare_type = var.compare_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `l7policy_id` - (Required, String) Specifies the forwarding policy ID.

* `l7rule_id` - (Optional, String) Specifies the forwarding rule ID.

* `compare_type` - (Optional, String) Specifies how requests are matched with the domain names or URL. Values options:
  **EQUAL_TO**, **REGEX**, **STARTS_WITH**.

* `value` - (Optional, String) Specifies the value of the match content.

* `type` - (Optional, String) Specifies the match type. Value options: **HOST_NAME**, **PATH**, **METHOD**, **HEADER**,
  **QUERY_STRING**, **SOURCE_IP**, **COOKIE**.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `l7rules` - Lists the L7 rules.
  The [l7rules](#Elb_l7rules) structure is documented below.

<a name="Elb_l7rules"></a>
The `l7rules` block supports:

* `id` - The ID of the forwarding policy.

* `type` - The type of the forwarding rule.

* `value` - The value of the match item.

* `compare_type` - How the requests are matched with the domain name or URL.

* `conditions` - The matching conditions of the forwarding rule. The [conditions](#Elb_l7rules_conditions) structure is
  documented below.

* `created_at` - The time when the forwarding rule was created.

* `updated_at` - The time when the forwarding rule was updated.

<a name="Elb_l7rules_conditions"></a>
The `conditions` block supports:

* `key` - The key of match item.

* `value` - The value of the match item.
