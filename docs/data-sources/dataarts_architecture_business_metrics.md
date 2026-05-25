---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_business_metrics"
description: |-
  Use this data source to query DataArts Architecture business metrics within HuaweiCloud.
---

# huaweicloud_dataarts_architecture_business_metrics

Use this data source to query DataArts Architecture business metrics within HuaweiCloud.

## Example Usage

### Query all business metrics under a specified workspace

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_architecture_business_metrics" "test" {
  workspace_id = var.workspace_id
}
```

### Query all business metrics under a specified workspace and using name filter

```hcl
variable "workspace_id" {}
variable "metric_name" {}

data "huaweicloud_dataarts_architecture_business_metrics" "test" {
  workspace_id = var.workspace_id
  name         = var.metric_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the business metrics are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the business metrics belong.

* `name` - (Optional, String) Specifies the name or code of the business metric to be fuzzy queried.

* `create_by` - (Optional, String) Specifies the creator of the business metric to be queried.

* `owner` - (Optional, String) Specifies the owner of the business metric to be queried.

* `status` - (Optional, String) Specifies the publishing status of the business metric to be queried.  
  The valid values are as follows:
  + **DRAFT**
  + **PUBLISH_DEVELOPING**
  + **PUBLISHED**
  + **OFFLINE_DEVELOPING**
  + **OFFLINE**
  + **REJECT**

* `biz_catalog_id` - (Optional, String) Specifies the process architecture ID to which the business metric belongs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `metrics` - The list of business metrics that matched filter parameters.  
  The [metrics](#dataarts_architecture_business_metrics_attr) structure is documented below.

<a name="dataarts_architecture_business_metrics_attr"></a>
The `metrics` block supports:

* `id` - The ID of the business metric.

* `name` - The name of the business metric.

* `biz_catalog_id` - The process architecture ID.

* `time_filters` - The statistical frequency.

* `interval_type` - The refresh frequency.

* `owner` - The name of person responsible for the indicator.

* `owner_department` - The indicator management department name.

* `destination` - The purpose of setting.

* `definition` - The indicator definition.

* `expression` - The calculation formula.

* `description` - The description.

* `apply_scenario` - The application scenarios.

* `technical_metric` - The related technical indicators.

* `measure` - The measurement object.

* `dimensions` - The statistical dimension.

* `general_filters` - The statistical caliber and modifiers.

* `data_origin` - The data sources.

* `unit` - The unit of measurement.

* `name_alias` - The indicator alias.

* `code` - The indicator encoding.

* `status` - The status of the business metric.

* `biz_catalog_path` - The process architecture path.

* `created_by` - The creator of the business metric.

* `updated_by` - The editor of the business metric.

* `created_at` - The creation time of the metric, in RFC3339 format.

* `updated_at` - The latest update time of the metric, in RFC3339 format.

* `technical_metric_type` - The related technical indicator type.

* `technical_metric_name` - The related technical indicator name.

* `l1` - The subject domain grouping Chinese name.

* `l2` - The subject field Chinese name.

* `l3` - The business object Chinese name.

* `biz_metric` - The business indicator synchronization status.

* `summary_status` - The synchronize statistics status.
