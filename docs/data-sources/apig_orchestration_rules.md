---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_orchestration_rules"
description: |-
  Use this data source to get the list of orchestration rules under specified instance within HuaweiCloud.
---

# huaweicloud_apig_orchestration_rules

Use this data source to get the list of orchestration rules under specified instance within HuaweiCloud.

## Example Usage

### Query all orchestration rules under specified APIG instance

```hcl
variable "instance_id" {}

data "huaweicloud_apig_orchestration_rules" "test" {
  instance_id = var.instance_id
}
```

### Query specified orchestration rule using its ID

```hcl
variable "instance_id" {}
variable "orchestration_rule_id" {}

data "huaweicloud_apig_orchestration_rules" "test" {
  instance_id = var.instance_id
  rule_id     = var.orchestration_rule_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the associated signatures.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the orchestration rules belong.

* `rule_id` - (Optional, String) Specifies the ID of the orchestration rule to be queried.

* `name` - (Optional, String) Specifies the name of the orchestration rule to be queried, fuzzy matching is supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - All orchestration rules that match the filter parameters.
  The [rules](#orchestration_rules) structure is documented below.

<a name="orchestration_rules"></a>
The `rules` block supports:

* `id` - The ID of the orchestration rule.

* `name` - The name of the orchestration rule.

* `strategy` - The type of the orchestration rule.  
  The values are as follows:
  + **list**: Maps the values ​​in the list to new values.
  + **range**: Maps the values ​​in the range to new values.
  + **hash**: The value of the request header is directly mapped to the new request header after hash calculation.
  + **hash_range**: Use the request parameter to generate a hash value, and then use the hash value to perform range
    arrangement.
  + **none_value**: Value returned when the request parameter is empty.
  + **default**: When the request parameters exist but no orchestration rule can match them, the orchestration
    mapping value of the default rule is returned.
  + **head_n**: Try to intercept the first N characters of the string as the new value.
  + **tail_n**: Try to intercept the last N characters of the string as the new value.

* `is_preprocessing` - Whether rule is a preprocessing rule.

* `mapped_param` - The parameter configuration after orchestration, in JSON format.

* `created_at` - The creation time of the orchestration rule, in RFC3339 format.

* `updated_at` - The latest update time of the orchestration rule, in RFC3339 format.
