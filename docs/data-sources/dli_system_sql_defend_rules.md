---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_system_sql_defend_rules"
description: |-
  Use this data source to query the system SQL defend rules of the DLI service within HuaweiCloud.
---

# huaweicloud_dli_system_sql_defend_rules

Use this data source to query the system SQL defend rules of the DLI service within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_dli_system_sql_defend_rules" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the system SQL defend rules are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list of system SQL defend rules that matched filter parameters.  
  The [rules](#system_sql_defend_rules_attr) structure is documented below.

<a name="system_sql_defend_rules_attr"></a>
The `rules` block supports:

* `id` - The ID of the rule type.
  - **static_0001**
  - **static_0002**
  - **static_0003**
  - **static_0004**
  - **static_0005**
  - **static_0006**
  - **static_0007**
  - **dynamic_0001**
  - **dynamic_0002**
  - **running_0002**
  - **running_0003**
  - **running_0004**

* `category` - The category of the rule.
  - **static**
  - **dynamic**
  - **running**

* `engines` - The list of supported engines that matched filter parameters.

* `actions` - The list of executable actions that matched filter parameters.

* `no_limit` - Whether the rule has a limit value.

* `description` - The description of the rule.

* `param` - The list of rule parameters that matched filter parameters.  
  The [param](#system_sql_defend_rules_param) structure is documented below.

<a name="system_sql_defend_rules_param"></a>
The `param` block supports:

* `default_value` - The default value of the threshold.

* `min` - The minimum value of the threshold.

* `max` - The maximum value of the threshold.

* `description` - The description of the parameter.
