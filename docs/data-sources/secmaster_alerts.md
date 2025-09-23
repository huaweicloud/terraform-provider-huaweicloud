---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_alerts"
description: |-
  Use this data source to get the list of SecMaster alerts.
---

# huaweicloud_secmaster_alerts

Use this data source to get the list of SecMaster alerts.

## Example Usage

```hcl
variable "workspace_id" {}
variable "from_date" {}
variable "to_date" {}

data "huaweicloud_secmaster_alerts" "test" {
  workspace_id = var.workspace_id
  from_date    = var.from_date
  to_date      = var.to_date

  condition {
    conditions {
      name = "severity"
      data = [ "severity", "=", "Tips" ]
    }

    logics = ["severity"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the alert belongs.

* `from_date` - (Optional, String) Specifies the search start time. For example: **2023-04-18T13:00:00.000+08:00**.

* `to_date` - (Optional, String) Specifies the search end time. For example: **2023-04-18T13:00:00.000+08:00**.

* `condition` - (Optional, List) Specifies the search condition expressions.
  The [condition](#condition) structure is documented below.

<a name="condition"></a>
The `condition` block supports:

* `conditions` - (Optional, List) Specifies the condition expression list.
  The [conditions](#condition_conditions) structure is documented below.

* `logics` - (Optional, List) Specifies the expression logic.

<a name="condition_conditions"></a>
The `conditions` block supports:

* `name` - (Optional, String) Specifies the expression name.

* `data` - (Optional, List) Specifies the expression content.
  + About `status` expression, e.g. **["handle_status", "!=", "Closed"]**.
  + About `name` expression, e.g. **["title", "contains", "test"]**.
  + About `level` expression, e.g. **["severity", "in", "Tips,Low"]**.
  + About `created_at` expression, e.g. **["create_time", ">=", "2024-08-15T19:18:38Z+0800"]**.
  + About `type.alert_type` expression, e.g. **["alert_type.alert_type", "=", "xxx"]**.
  + About `first_occurrence_time` expression, e.g. **["first_observed_time", "<=", "2024-08-23T20:09:26Z+0800"]**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `alerts` - The alert list.
  The [alerts](#alerts_struct) structure is documented below.

<a name="alerts_struct"></a>
The `alerts` block supports:

* `id` - The alert ID.

* `name` - The alert name.

* `description` - The alert description.

* `type` - The alert type configuration.
  The [type](#alerts_type_struct) structure is documented below.

* `level` - The alert level.

* `status` - The alert status.

* `data_source` - The data source configuration.
  The [data_source](#alerts_data_source_struct) structure is documented below.

* `first_occurrence_time` - The first occurrence time of the alert.

* `owner` - The owner name.

* `last_occurrence_time` - The last occurrence time of the alert.

* `planned_closure_time` - The planned closure time of the alert.

* `verification_status` - The verification status.

* `stage` - The stage of the alert.

* `debugging_data` - Whether it's debugging data.

* `labels` - The labels.

* `close_reason` - The close reason.

* `close_comment` - The close comment.

* `creator` - The creator name.

* `created_at` - The creation time.

* `updated_at` - The update time.

* `version` - The version of the data source of an alert.

* `domain_id` - The ID of the account (domain_id) to whom the data is delivered and hosted.

* `project_id` - The ID of project where the account to whom the data is delivered and hosted belongs to.

* `region_id` - The ID of the region where the account to whom the data is delivered and hosted belongs to.

* `workspace_id` - The ID of the current workspace.

* `arrive_time` - The data receiving time.

* `count` - The times of the alert occurrences.

* `data_class_id` - The data class ID.

* `environment` - The coordinates of the environment where the alert was generated.
  The [environment](#alerts_environment) structure is documented below.

* `network_list` - The network information.
  The [network_list](#alerts_network_list) structure is documented below.

* `user_info` - The user information.
  The [user_info](#alerts_user_info) structure is documented below.

* `resource_list` - The affected resources.
  The [resource_list](#alerts_resource_list) structure is documented below.

* `process` - The process information.
  The [process](#alerts_process) structure is documented below.

* `malware` - The malware information.
  The [malware](#alerts_malware) structure is documented below.

* `remediation` - The remedy measure.
  The [remediation](#alerts_remediation) structure is documented below.

* `file_info` - The file information.
  The [file_info](#alerts_file_info) structure is documented below.

<a name="alerts_data_source_struct"></a>
The `data_source` block supports:

* `product_feature` - The product feature.

* `product_name` - The product name.

* `source_type` - The source type.

<a name="alerts_type_struct"></a>
The `type` block supports:

* `category` - The category.

* `alert_type` - The alert type.

<a name="alerts_environment"></a>
The `environment` block supports:

* `vendor_type` - The environment provider. The value can be **HWCP**, **HWC**, **AWS**, **Azure**, or **GCP**.

* `cross_workspace_id` - The ID of the source workspace for the data delivery.

* `domain_id` - The domain ID.

* `project_id` - The project ID. The default value is empty for global services.

* `region_id` - The region ID. **global** is returned for global services.

<a name="alerts_network_list"></a>
The `network_list` block supports:

* `protocol` - The protocol.

* `src_ip` - The source IP address.

* `src_port` - The source port.

* `dest_ip` - The destination IP address.

* `dest_port` - The destination port.

* `src_geo` - The geographical location of the source IP address.
  The [src_geo](#alerts_network_list_geo) structure is documented below.

* `dest_geo` - The geographical location of the destination IP address.
  The [dest_geo](#alerts_network_list_geo) structure is documented below.

* `direction` - The direction. The value can be **IN** or **OUT**.

* `src_domain` - The source domain name.

* `dest_domain` - The destination domain name.

<a name="alerts_network_list_geo"></a>
The `dest_geo` block supports:

* `longitude` - The longitude of geographical location.

* `latitude` - The latitude of the geographical location.

* `city_code` - The city code.

* `country_code` - The country code.

<a name="alerts_user_info"></a>
The `user_info` block supports:

* `user_id` - The user ID.

* `user_name` - The user name.

<a name="alerts_resource_list"></a>
The `resource_list` block supports:

* `domain_id` - The ID of the account to which the resource belongs.

* `ep_id` - The enterprise project ID.

* `ep_name` - The enterprise project name.

* `id` - The resource ID.

* `name` - The resource name.

* `project_id` - The ID of the account to which the resource belongs.

* `provider` - The cloud service name, which is the same as the provider field in the RMS service.

* `region_id` - The region ID.

* `tags` - The resource tags.

* `type` - The resource type.

<a name="alerts_process"></a>
The `process` block supports:

* `process_child_cmdline` - The subprocess command line.

* `process_child_name` - The subprocess name.

* `process_child_path` - The subprocess execution file path.

* `process_child_pid` - The subprocess ID.

* `process_child_uid` - The subprocess user ID.

* `process_cmdline` - The process command line.

* `process_launche_time` - The process start time.

* `process_name` - The process name.

* `process_parent_cmdline` - The parent process command line.

* `process_parent_name` - The parent process name.

* `process_parent_path` - The parent process execution file path.

* `process_parent_pid` - The parent process ID.

* `process_parent_uid` - The parent process user ID.

* `process_path` - The process execution file path.

* `process_pid` - The process ID.

* `process_terminate_time` - The process end time.

* `process_uid` - The process user ID.

<a name="alerts_malware"></a>
The `malware` block supports:

* `malware_class` - The malware category.

* `malware_family` - The malicious family.

<a name="alerts_remediation"></a>
The `remediation` block supports:

* `recommendation` - The recommended solution.

* `url` - The link to the general fix information for the incident.

<a name="alerts_file_info"></a>
The `file_info` block supports:

* `file_path` - The file path.

* `file_content` - The file content.

* `file_new_path` - The file new path.

* `file_hash` - The file hash value.

* `file_md5` - The file MD5 value.

* `file_sha256` - The file SHA256 value.

* `file_attr` - The file attribute.
