---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_business_metric"
description: ""
---

# huaweicloud_dataarts_architecture_business_metric

Manages a DataArts Architecture business metric resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "biz_catalog_id" {}
variable "owner" {}
variable "owner_department" {}

resource "huaweicloud_dataarts_architecture_business_metric" "test" {
  workspace_id     = var.workspace_id
  biz_catalog_id   = var.biz_catalog_id
  owner            = var.owner
  owner_department = var.owner_department
  name             = "test-name"
  name_alias       = "alias-name"
  time_filters     = "Monthly,Yearly"
  interval_type    = "HOUR"
  destination      = "test destination"
  definition       = "test definition"
  expression       = "a+b+c"
  description      = "test description"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the ID of DataArts Studio workspace.
  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the indicator name. This value can contain only letters, digits, parentheses,
  commas, spaces, and special characters `+#-_/[]`. The maximum length is 500 characters.

* `biz_catalog_id` - (Required, String) Specifies the process architecture ID.

* `time_filters` - (Required, String) Specifies the statistical frequency. Use commas to separate multiple fields.

* `interval_type` - (Required, String) Specifies the refresh frequency. Value values are: **REAL_TIME**, **HOUR**,
  **HALF_DAY**, **DAY**, **WEEK**, **DOUBLE_WEEK**, **MONTH**, **QUART**, **HALF_YEAR** and **YEAR**.

* `owner` - (Required, String) Specifies the name of person responsible for the indicator. The owner must exist in the
  system. The maximum length is 100 characters.

* `owner_department` - (Required, String) Specifies the indicator management department name. The maximum length is 600
  characters.

* `destination` - (Required, String) Specifies the purpose of setting. The maximum length is 1000 characters.

* `definition` - (Required, String) Specifies the indicator definition. The maximum length is 1000 characters.

* `expression` - (Required, String) Specifies the calculation formula. The maximum length is 1000 characters.

* `description` - (Optional, String) Specifies the description. The maximum length is 600 characters.

* `apply_scenario` - (Optional, String) Specifies the application scenarios. The maximum length is 255 characters.

* `technical_metric` - (Optional, Int) Specifies the related technical indicators.

* `measure` - (Optional, String) Specifies the measurement object. The maximum length is 255 characters.

* `dimensions` - (Optional, String) Specifies the statistical dimension. The maximum length is 1000 characters.

* `general_filters` - (Optional, String) Specifies the statistical caliber and modifiers. The maximum length is 1000 characters.

* `data_origin` - (Optional, String) Specifies the data sources. The value needs to be a descriptive value. The maximum
  length is 1000 characters.

* `unit` - (Optional, String) Specifies the unit of measurement. The value of this field needs to be a quantifier, for example,
  **percentage**, **hour** or **minute**. The maximum length is 50 characters.

* `name_alias` - (Optional, String) Specifies the indicator alias. The maximum length is 500 characters.

* `code` - (Optional, String) Specifies the indicator encoding. If this field is not specified, an indicator code will
  be automatically generated. The value of this field does not support changing to empty.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status. Valid values are: **DRAFT**, **PUBLISH_DEVELOPING**, **PUBLISHED**, **OFFLINE_DEVELOPING**,
  **OFFLINE** and **REJECT**.

* `biz_catalog_path` - The process architecture path.

* `created_by` - The creator.

* `updated_by` - The editor.

* `created_at` - The creation time.

* `updated_at` - The update time.

* `technical_metric_type` - The related technical indicator type.

* `technical_metric_name` - The related technical indicator name.

* `l1` - The subject domain grouping Chinese name.

* `l2` - The subject field Chinese name.

* `l3` - The business object Chinese name.

* `biz_metric` - The business indicator synchronization status. Valid values are: **NO_NEED**, **CREATE_SUCCESS**,
  **CREATE_FAILED**, **UPDATE_SUCCESS**, **UPDATE_FAILED**, **SUMMARY_SUCCESS**, **SUMMARY_FAILED**, **RUNNING** and **OFFLINE**.

* `summary_status` - The synchronize statistics status. Valid values are: **NO_NEED**, **CREATE_SUCCESS**,
  **CREATE_FAILED**, **UPDATE_SUCCESS**, **UPDATE_FAILED**, **SUMMARY_SUCCESS**, **SUMMARY_FAILED**, **RUNNING** and **OFFLINE**.

## Import

The DataArts Architecture business metric resource can be imported using the `workspace_id` and `id`, separated by a
slash, e.g.

```bash
$ terraform import huaweicloud_dataarts_architecture_business_metric.test <workspace_id>/<id>
```
