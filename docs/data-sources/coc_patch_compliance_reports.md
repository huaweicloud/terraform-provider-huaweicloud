---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_patch_compliance_reports"
description: |-
  Use this data source to get the list of COC patch compliance reports.
---

# huaweicloud_coc_patch_compliance_reports

Use this data source to get the list of COC patch compliance reports.

## Example Usage

```hcl
data "huaweicloud_coc_patch_compliance_reports" "test" {}
```

## Argument Reference

The following arguments are supported:

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `name` - (Optional, String) Specifies the name.

* `instance_id` - (Optional, String) Specifies the ECS instance ID.

* `ip` - (Optional, String) Specifies the internal network IP address.

* `eip` - (Optional, String) Specifies the elastic IP address.

* `operating_system` - (Optional, String) Specifies the OS.
  Values can be **HuaweiCloudEulerOS**, **CentOS** and **EulerOS**.

* `region` - (Optional, String) Specifies the region.

* `group` - (Optional, String) Specifies the group.

* `compliant_status` - (Optional, String) Specifies the compliance status.
  Values can be **non_compliant** and **compliant**.

* `order_id` - (Optional, String) Specifies the service ticket ID.

* `sort_dir` - (Optional, String) Specifies the sorting order.
  Values can be as follows:
  + **asc**: The query results are displayed in ascending order.
  + **desc**: The query results are displayed in the descending order.

* `sort_key` - (Optional, String) Specifies the sorting field.
  Values can be **report_time**.

* `report_scene` - (Optional, String) Specifies the report scenario.
  Values can be **CCE** and **ECS**.

* `cce_info_id` - (Optional, String) Specifies the CCE cluster information ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instance_compliant` - Indicates the node compliance report.

  The [instance_compliant](#instance_compliant_struct) structure is documented below.

<a name="instance_compliant_struct"></a>
The `instance_compliant` block supports:

* `compliant_summary` - Indicates the compliance patch information.

  The [compliant_summary](#instance_compliant_compliant_summary_struct) structure is documented below.

* `non_compliant_summary` - Indicates the non-compliant patch information.

  The [non_compliant_summary](#instance_compliant_non_compliant_summary_struct) structure is documented below.

* `execution_summary` - Indicates the execution information.

  The [execution_summary](#instance_compliant_execution_summary_struct) structure is documented below.

* `id` - Indicates the ID.

* `enterprise_project_id` - Indicates the enterprise project ID.

* `name` - Indicates the node name.

* `instance_id` - Indicates the node ID.

* `node_id` - Indicates the CCE cluster node ID.

* `ip` - Indicates the node IP address.

* `eip` - Indicates the elastic IP address.

* `region` - Indicates the region.

* `group` - Indicates the group.

* `report_scene` - Indicates the reporting scenario, CCE or ECS.

* `cce_info_id` - Indicates the CCE cluster information ID.

* `status` - Indicates the compliance status.

* `baseline_id` - Indicates the baseline ID.

* `baseline_name` - Indicates the baseline name.

* `rule_type` - Indicates the baseline rule type.

* `operating_system` - Indicates the OS.

<a name="instance_compliant_compliant_summary_struct"></a>
The `compliant_summary` block supports:

* `compliant_count` - Indicates the number of compliant patches.

* `severity_summary` - Indicates the compliance summary.

  The [severity_summary](#compliant_summary_severity_summary_struct) structure is documented below.

<a name="compliant_summary_severity_summary_struct"></a>
The `severity_summary` block supports:

* `critical_count` - Indicates the number of major compliance reports.

* `high_count` - Indicates the number of high compliance reports.

* `informational_count` - Indicates the number of informational compliance reports.

* `low_count` - Indicates the number of low compliance reports.

* `medium_count` - Indicates the number of medium compliance reports.

* `unspecified_count` - Indicates the number of unspecified compliance reports.

<a name="instance_compliant_non_compliant_summary_struct"></a>
The `non_compliant_summary` block supports:

* `non_compliant_count` - Indicates the number of non-compliant patches.

* `severity_summary` - Indicates the compliance summary.

  The [severity_summary](#non_compliant_summary_severity_summary_struct) structure is documented below.

<a name="non_compliant_summary_severity_summary_struct"></a>
The `severity_summary` block supports:

* `critical_count` - Indicates the number of major compliance reports.

* `high_count` - Indicates the number of high compliance reports.

* `informational_count` - Indicates the number of informational compliance reports.

* `low_count` - Indicates the number of low compliance reports.

* `medium_count` - Indicates the number of medium compliance reports.

* `unspecified_count` - Indicates the number of unspecified compliance reports.

<a name="instance_compliant_execution_summary_struct"></a>
The `execution_summary` block supports:

* `order_id` - Indicates the service ticket ID.

* `job_id` - Indicates the script execution ID.

* `report_time` - Indicates the reporting time.
