---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_inspection_reports"
description: |-
  Use this data source to query DAS inspection reports within HuaweiCloud.
---

# huaweicloud_das_inspection_reports

Use this data source to query DAS inspection reports within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_das_inspection_reports" "test" {
  start_time     = "2025-05-01T00:00:00+08:00"
  end_time       = "2025-05-02T00:00:00+08:00"
  datastore_type = "MySQL"
}
```

### Filter by health rank

```hcl
data "huaweicloud_das_inspection_reports" "test" {
  start_time     = "2025-05-01T00:00:00+08:00"
  end_time       = "2025-05-02T00:00:00+08:00"
  datastore_type = "MySQL"
  health_rank    = "healthy"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the inspection reports are located.  
  If omitted, the provider-level region will be used.

* `start_time` - (Required, String) Specifies the start time of the inspection report.  

* `end_time` - (Required, String) Specifies the end time of the inspection report.  

* `datastore_type` - (Required, String) Specifies the database type.  
  The valid values are as follows:
  + **MySQL**
  + **TaurusDB**
  + **GaussDB**
  + **MariaDB**

* `health_rank` - (Optional, String) Specifies the health rank of the inspection report.  
  The valid values are as follows:
  + **healthy**
  + **sub_healthy**
  + **dangerous**
  + **high_risk**

* `sort_field` - (Optional, String) Specifies the field used for sorting.  
  The valid value is **create_at**, which means the generation time of the inspection report.

* `asc` - (Optional, Bool) Specifies whether to sort in ascending order.  
  Defaults to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `reports` - The list of inspection reports that matched filter parameters.  
  The [reports](#inspection_reports_attr) structure is documented below.

<a name="inspection_reports_attr"></a>
The `reports` block supports:

* `task_id` - The ID of the inspection report.

* `instance_id` - The ID of the instance.

* `instance_name` - The name of the instance.

* `cpu` - The CPU size.

* `mem` - The memory size in GB.

* `disk_size` - The disk size in GB.

* `created_time` - The generation time of the inspection report, in RFC3339 format.

* `start_time` - The start time of the diagnosis, in RFC3339 format.

* `end_time` - The end time of the diagnosis, in RFC3339 format.

* `health_rank` - The health rank of the instance.

* `score` - The score of the inspection.

* `lost_points_details` - The list of lost points details.  
  The [lost_points_details](#inspection_reports_lost_points_details_attr) structure is documented below.

<a name="inspection_reports_lost_points_details_attr"></a>
The `lost_points_details` block supports:

* `risk_level` - The risk level.

* `metric` - The metric name.

* `metric_value` - The value of the metric.

* `deducted_points` - The deducted points.

* `deducted_condition` - The deducted condition.

* `deducted_formula` - The deducted formula.

* `suggestions` - The optimization suggestions.
