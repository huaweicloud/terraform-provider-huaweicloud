---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_incidents"
description: |-
  Use this data source to get the list of SecMaster incidents.
---

# huaweicloud_secmaster_incidents

Use this data source to get the list of SecMaster incidents.

## Example Usage

```hcl
variable "workspace_id" {}
variable "from_date" {}
variable "to_date" {}

data "huaweicloud_secmaster_incidents" "test" {
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

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the incident belongs.

* `from_date` - (Optional, String) Specifies the search start time. For example: **2023-04-18T13:00:00.000+08:00**.

* `to_date` - (Optional, String) Specifies the search end time. For example: **2023-04-18T13:00:00.000+08:00**.

* `condition` - (Optional, List) Specifies the search condition expression.
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
  + About `incident_type.incident_type` expression, e.g. **["incident_type.incident_type", "=", "xxx"]**.
  + About `first_occurrence_time` expression, e.g. **["first_observed_time", "<=", "2024-08-23T20:09:26Z+0800"]**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `incidents` - The incident list.
  The [incidents](#incidents_struct) structure is documented below.

<a name="incidents_struct"></a>
The `incidents` block supports:

* `id` - The incident ID.

* `name` - The incident name.

* `description` - The incident description.

* `type` - The incident type configuration.
  The [type](#incidents_type_struct) structure is documented below.

* `level` - The incident level.

* `status` - The incident status.

* `data_source` - The data source configuration.
  The [data_source](#incidents_data_source_struct) structure is documented below.

* `first_occurrence_time` - The first occurrence time of the incident.

* `owner` - The user name of the owner.

* `last_occurrence_time` - The last occurrence time of the incident.

* `planned_closure_time` - The planned closure time of the incident.

* `verification_status` - The verification status.

* `stage` - The stage of the incident.

* `debugging_data` - Whether it's a debugging data.

* `labels` - The labels.

* `close_reason` - The close reason.

* `close_comment` - The close comment.

* `creator` - The name creator name.

* `created_at` - The creation time.

* `updated_at` - The update time.

* `version` - The version of the data source of an incident.

* `domain_id` - The ID of the account (domain_id) to whom the data is delivered and hosted.

* `project_id` - The ID of project where the account to whom the data is delivered and hosted belongs to.

* `region_id` - The ID of the region where the account to whom the data is delivered and hosted belongs to.

* `workspace_id` - The ID of the current workspace.

* `arrive_time` - The data receiving time.

* `count` - The times of the incident occurrences.

* `ipdrr_phase` - The handling phase No.

<a name="incidents_data_source_struct"></a>
The `data_source` block supports:

* `product_feature` - The product feature.

* `product_name` - The product name.

* `source_type` - The source type.

<a name="incidents_type_struct"></a>
The `type` block supports:

* `category` - The category.

* `incident_type` - The incident type.
