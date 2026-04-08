---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_sql_defend_rules"
description: |-
  Use this data source to query the SQL defend rules of the DLI service within HuaweiCloud.
---

# huaweicloud_dli_sql_defend_rules

Use this data source to query the SQL defend rules of the DLI service within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_dli_sql_defend_rules" "test" {}
```

### Filter by rule name

```hcl
variable "rule_name" {}

data "huaweicloud_dli_sql_defend_rules" "test" {
  rule_name = var.rule_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the SQL defend rules are located.
  If omitted, the provider-level region will be used.

* `queue_name` - (Optional, String) Specifies the queue name used to filter SQL defend rules.

* `rule_name` - (Optional, String) Specifies the rule name used to filter SQL defend rules.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list of SQL defend rules that matched filter parameters.  
  The [rules](#sql_defend_rules_attr) structure is documented below.

<a name="sql_defend_rules_attr"></a>
The `rules` block supports:

* `name` - The name of the rule.

* `uuid` - The UUID of the rule.

* `id` - The ID of the rule.
  + **static_0001**
  + **static_0002**
  + **static_0003**
  + **static_0004**
  + **static_0005**
  + **static_0006**
  + **static_0007**
  + **dynamic_0001**
  + **dynamic_0002**
  + **running_0002**
  + **running_0003**
  + **running_0004**

* `category` - The category of the rule.
  + **static**
  + **dynamic**
  + **running**

* `description` - The description of the rule.

* `engine_rules` - The engine rules of the rule.

* `project_id` - The project ID of the rule.

* `queue_names` - The list of queue names that matched filter parameters.

* `sys_desc` - The system description of the rule.

* `created_at` - The creation time of the rule.

* `updated_at` - The update time of the rule.
