---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_indicators"
description: |-
  Use this data source to get the list of SecMaster indicators.
---

# huaweicloud_secmaster_indicators

Use this data source to get the list of SecMaster indicators.

## Example Usage

```hcl
variable "workspace_id" {}
variable "from_date" {}
variable "to_date" {}

data "huaweicloud_secmaster_indicators" "test" {
  workspace_id = var.workspace_id
  from_date    = var.from_date
  to_date      = var.to_date

  condition {
    conditions {
      name = "name"
      data = ["name", "=", "test"]
    }

    conditions {
      name = "status"
      data = ["status", "=", "Open"]
    }

    logics = ["name", "and", "status"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the indicator belongs.

* `ids` - (Optional, List) Specifies the indicator IDs.

* `data_class_id` - (Optional, String) Specifies the data class ID.

* `from_date` - (Optional, String) Specifies the search start time. For example: **2023-04-18T13:00:00.000+08:00**.

* `to_date` - (Optional, String) Specifies the search end time. For example: **2023-04-18T13:00:00.000+08:00**.

* `condition` - (Optional, List) Specifies the search condition expressions.
  The [condition](#condition) structure is documented below.

<a name="condition"></a>
The `condition` block supports:

* `conditions` - (Optional, List) Specifies the condition expression list.
  The [conditions](#condition_conditions) structure is documented below.

* `logics` - (Optional, List) Specifies the expression logic.
  For example, **["conditions.name1", "and", "conditions.name2"]**.

<a name="condition_conditions"></a>
The `conditions` block supports:

* `name` - (Optional, String) Specifies the expression name.

* `data` - (Optional, List) Specifies the expression content.
  + About `threat_degree` expression, e.g. **["verdict", "=", "Gray"]**.
  + About `created_at` expression, e.g. **["create_time", ">=", "2024-08-15T19:18:38Z+0800"]**.
  + About `updated_at` expression, e.g. **["update_time", ">=", "2024-08-15T19:18:38Z+0800"]**.
  + About `type.indicator_type` expression, e.g. **["indicator_type.indicator_type", "=", "IPv6"]**.
  + About `first_occurrence_time` expression, e.g. **["first_report_time", ">=", "2024-08-20T14:52:06Z+0800"]**.
  + About `last_occurrence_time` expression, e.g. **["last_report_time", ">=", "2024-08-20T14:52:06Z+0800"]**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `indicators` - The list of indicators.

  The [indicators](#indicators_struct) structure is documented below.

<a name="indicators_struct"></a>
The `indicators` block supports:

* `id` - The indicator ID.

* `name` - The indicator name.

* `threat_degree` - The threat degree.

* `type` - The indicator type.
  The [type](#indicators_type_struct) structure is documented below.

* `status` - The status.

* `confidence` - The confidence. The value range is `80` to `100`.

* `first_occurrence_time` - The first occurred time of the indicator.

* `last_occurrence_time` - The last occurred time of the indicator.

* `created_at` - The creation time.

* `updated_at` - The update time.

* `granularity` - The confidentiality level.
  + **1**: First discovery.
  + **2**: Self-produced data.
  + **3**: Purchase required.
  + **4**: Direct query from the external network.

* `data_class_id` - The data class ID.

* `value` - The value. Such as **ip**, **url**, **domain** etc.

* `data_source` - The data source configuration.
  The [data_source](#indicators_data_source_struct) structure is documented below.

* `environment` - The coordinates of the environment where the indicator was generated.
  The [environment](#indicators_environment_struct) structure is documented below.

* `revoked` - Whether the indicator is discard.

* `is_deleted` - Whether the indicator is deleted.

* `workspace_id` - The workspace ID.

* `project_id` - The project ID.

<a name="indicators_type_struct"></a>
The `indicator_type` block supports:

* `indicator_type` - The indicator type.

* `id` - The indicator type ID.

<a name="indicators_data_source_struct"></a>
The `data_source` block supports:

* `source_type` - The data source type.

* `domain_id` - The domain ID.

* `project_id` - The project ID.

* `region_id` - The region ID.

<a name="indicators_environment_struct"></a>
The `environment` block supports:

* `vendor_type` - The environment provider. Such as **HWCP**, **HWC**, **AWS** or **Azure**.

* `domain_id` - The domain ID.

* `project_id` - The project ID.

* `region_id` - The region ID.
