---
subcategory: "Database Security Service (DBSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dbss_audit_data_masking_rules"
description: |-
  Use this data source to get a list of privacy data masking rules.
---

# huaweicloud_dbss_audit_data_masking_rules

Use this data source to get a list of privacy data masking rules.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dbss_audit_data_masking_rules" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the audit instance ID to which the privacy data masking rules belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list of the privacy data masking rules.

  The [rules](#rules_struct) structure is documented below.

<a name="rules_struct"></a>
The `rules` block supports:

* `id` - The ID of the privacy data masking rule.

* `name` - The name of the privacy data masking rule.

* `type` - The type of the privacy data masking rule.

* `status` - The status of the privacy data masking rule.

* `regex` - The regular expression of the privacy data masking rule.

* `mask_value` - The privacy data display substitution value.

* `operate_time` - The operation time of the privacy data masking rule, in UTC format.
