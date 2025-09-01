---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmasterv2_alert_rule_template_detail"
description: |-
  Use this data source to query a specific alert rule template (v2) detail.
---

# huaweicloud_secmasterv2_alert_rule_template_detail

Use this data source to query a specific alert rule template (v2) detail.

## Example Usage

```hcl
data "huaweicloud_secmasterv2_alert_rule_template_detail" "test" {
  workspace_id = "your_workspace_id"
  template_id  = "your_template_id"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `template_id` - (Required, String) Specifies the alert rule template ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rule_template_id` - The template ID.

* `template_name` - The template name.

* `accumulated_times` - The cumulative number of times.

* `cu_quota_amount` - The CU quota amount.

* `description` - The alert rule template description.

* `environment` - The environment type.
  The valid values are as follows:
  + **PROD**: Production environment.
  + **TEST**: Test environment.

* `job_mode` - The job mode.
  The valid values are as follows:
  + **STREAMING**: Streaming processing.
  + **BATCH**: Batch processing.
  + **SEARCH**: Search.

* `job_mode_setting` - The job mode setting.
  The [job_mode_setting](#job_mode_setting_attr) structure is documented below.

* `job_output_setting` - The job output setting.
  The [job_output_setting](#job_output_setting_attr) structure is documented below.

* `process_error` - The process error.

* `process_status` - The process status.
  The value can be **COMPLETED**, **CREATING**, **UPDATING**, **ENABLING**, **DISABLING**, **DELETING**,
  **CREATE_FAILED**, **UPDATE_FAILED**, **ENABLE_FAILED**, **DISABLE_FAILED**, **DELETE_FAILED** or **RECOVERING**.

* `query_type` - The query type.
  The value can be **SQL** or **CBSL**.

* `script` - The script.

* `status` - The status.
  The value can be **ENABLED** or **DISABLED**.

* `table_name` - The table name.

* `triggers` - The triggers information.
  The [triggers](#triggers_attr) structure is documented below.

* `create_by` - The ID of the creator.

* `create_time` - The creation time, in milliseconds timestamp.

* `update_by` - The ID of the updater.

* `update_time` - The update time, in milliseconds timestamp.

<a name="job_mode_setting_attr"></a>
The `job_mode_setting` block supports:

* `batch_overtime_strategy_interval` - The batch processing strategy overtime interval.

* `batch_overtime_strategy_unit` - The batch processing strategy overtime time unit.
  The value can be **MINUTE**, **HOUR**, **DAY** or **MONTH**.

* `search_delay_interval` - The search delay interval.

* `search_delay_unit` - The search delay time unit.
  The value can be **MINUTE**, **HOUR**, **DAY** or **MONTH**.

* `search_frequency_interval` - The search frequency interval.

* `search_frequency_unit` - The search frequency interval time unit.
  The value can be **MINUTE**, **HOUR**, **DAY** or **MONTH**.

* `search_overtime_interval` - The search overtime interval.

* `search_overtime_unit` - The search overtime interval time unit.
  The value can be **MINUTE**, **HOUR**, **DAY** or **MONTH**.

* `search_period_interval` - The search period interval.

* `search_period_unit` - The search period interval time unit.
  The value can be **MINUTE**, **HOUR**, **DAY** or **MONTH**.

* `search_table_id` - The search table ID.

* `search_table_name` - The search table name.

* `streaming_checkpoint_ttl_interval` - The streaming checkpoint TTL interval.

* `streaming_checkpoint_ttl_unit` - The streaming checkpoint TTL interval time unit.
  The value can be **MINUTE**, **HOUR**, **DAY** or **MONTH**.

* `streaming_startup_mode` - The startup mode.including  () and  (REFRESH_NEW).
  The valid values are as follows:
  + **UPGRADE**: Upgrade startup.
  + **REFRESH_NEW**: Refresh new startup.

* `streaming_state_ttl_unit` - The streaming state TTL time unit.
  The value can be **MINUTE**, **HOUR**, **DAY** or **MONTH**.

<a name="job_output_setting_attr"></a>
The `job_output_setting` block supports:

* `alert_custom_properties` - The alert custom properties.

* `alert_description` - The alert description.

* `alert_grouping` - The alert grouping flag.

* `alert_mapping` - The alert mapping.

* `alert_name` - The alert name.

* `alert_remediation` - The alert remediation suggestion.

* `alert_severity` - The alert severity.
  The value can be **TIPS**, **LOW**, **MEDIUM**, **HIGH** or **FATAL**.

* `alert_suppression` - The alert suppression flag.

* `alert_type` - The alert type.

* `entity_extraction` - The entity extraction.

* `field_mapping` - The field mapping.

<a name="triggers_attr"></a>
The `triggers` block supports:

* `accumulated_times` - The cumulative number of times.

* `expression` - The expression.

* `job_id` - The job ID.

* `mode` - The mode.

* `operator` - The operator type.
  The valid values are as follows:
  + **EQ**: equal,
  + **NE**: not equal,
  + **GT**: greater than,
  + **LT**: less than.

* `severity` - The severity.
