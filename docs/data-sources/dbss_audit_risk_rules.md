---
subcategory: "Database Security Service (DBSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dbss_audit_risk_rules"
description: |-
  Use this data source to get a list of risk rules.
---

# huaweicloud_dbss_audit_risk_rules

Use this data source to get a list of risk rules.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dbss_audit_risk_rules" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the audit instance ID to which the risk rules belong.

* `name` - (Optional, String) Specifies the name of the risk rule. Supports fuzzy search.

* `risk_level` - (Optional, String) Specifies  the risk level of the risk rule.
  The valid values are as follows:
  + **LOW**
  + **MEDIUM**
  + **HIGH**
  + **NO_RISK**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list of the risk rules.

  The [rules](#rules_struct) structure is documented below.

<a name="rules_struct"></a>
The `rules` block supports:

* `id` - The ID of the risk rule.

* `name` - The name of the risk rule.

* `type` - The type of the risk rule.

* `status` - The status of the risk rule.

* `feature` - The risk characteristics of the risk rule.

* `rank` - The priority of the risk rule.

* `risk_level` - The risk level of the risk rule.
