---
subcategory: "Database Security Service (DBSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dbss_audit_sql_injection_rules"
description: |-
  Use this data source to get a list of SQL injection rules.
---

# huaweicloud_dbss_audit_sql_injection_rules

Use this data source to get a list of SQL injection rules.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dbss_audit_sql_injection_rules" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the audit instance ID to which the SQL injection rules belong.

* `risk_levels` - (Optional, String) Specifies the risk level of the SQL injection rule.
  The valid values are as follows:
  + **HIGH**
  + **MEDIUM**
  + **LOW**
  + **NO_RISK**

  This parameter can also be set to multiple enumerated values and separated by a comma (,). e.g. **HIGH,LOW**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list of the SQL injection rules.

  The [rules](#rules_struct) structure is documented below.

<a name="rules_struct"></a>
The `rules` block supports:

* `id` - The ID of the SQL injection rule.

* `name` - The name of the SQL injection rule.

* `type` - The type of the SQL injection rule.

* `status` - The status of the SQL injection rule.

* `risk_level` - The risk level of the SQL injection rule.

* `rank` - The rank of the SQL injection rule.

* `feature` - The SQL command characteristics of the SQL injection rule.

* `regex` - The regular expression content of the SQL injection rule.
