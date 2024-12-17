---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_metric_rule"
description: |-
  Manages an LTS log metric rule resource within HuaweiCloud.
---

# huaweicloud_lts_metric_rule

Manages an LTS log metric rule resource within HuaweiCloud.

## Example Usage

```hcl
variable "metric_rule_name" {}
variable "log_group_id" {}
variable "log_stream_id" {}
variable "sinks_metric_name" {}
variable "prometheus_instance_name" {}
variable "prometheus_instance_id" {}
variable "aggregator_field" {}
variable "log_filters" {
  type = list(object({
    type = string
    filters = list(object({
      type  = string
      key   = string
      value = string
    }))
  }))
}

resource "huaweicloud_lts_metric_rule" "test" {
  name          = var.metric_rule_name
  status        = "enable"
  log_group_id  = var.log_group_id
  log_stream_id = var.log_stream_id

  sampler {
    type  = "random"
    ratio = "0.5"
  }

  sinks {
    type        = "aom"
    metric_name = var.sinks_metric_name
    name        = var.prometheus_instance_name
    instance_id = var.prometheus_instance_id
  }

  aggregator {
    type  = "count"
    field = var.aggregator_field
  }

  window_size = "PT1M"

  filter {
    type = "and"
  
    dynamic "filters" {
      for_each = var.log_filters
      content {
        type = filters.value.type
  
        dynamic "filters" {
          for_each = filters.value.filters
          content {
            type  = filters.value.type
            key   = filters.value.key
            value = filters.value.value
          }
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String) Specifies the name of the log metric rule. The name must be unique.  
  The name contain a maximum of `256` characters, and only English letters, digits, and hyphens (-) are allowed.

* `status` - (Required, String) Specifies the status of the log metric rule.  
  The valid values are as follows:
  + **enable**
  + **disable**

* `log_group_id` - (Required, String) Specifies the log group ID to which the log metric rule belongs.

* `log_stream_id` - (Required, String) Specifies the log stream ID to which the log metric rule belongs.

* `sampler` - (Required, List) Specifies the sampling configuration of the log.
  The [sampler](#metric_rule_sampler) structure is documented below.

* `sinks` - (Required, List) Specifies the storage location of the generated metrics.
  The [sinks](#metric_rule_sinks) structure is documented below.

* `aggregator` - (Required, List) Specifies the configuration of log statistics mode.
  The [aggregator](#metric_rule_aggregator) structure is documented below.

* `window_size` - (Required, String) Specifies the interval time for processing data windows.  
  The valid values are as follows:
  + **PT5S**: Indicates `5` seconds.
  + **PT1M**: Indicates `1` minute.
  + **PT5M**: Indicates `5` minute.

* `report` - (Optional, Bool) Specifies whether to report data to sinks, defaults to **false**.

* `filter` - (Optional, List) Specifies the configuration of log filtering rule.
  The [filter](#metric_rule_filter) structure is documented below.

* `description` - (Optional, String) Specifies the description of the log metric rule.

<a name="metric_rule_sampler"></a>
The `sampler` block supports:

* `type` - (Required, String) Specifies the type of the log sampling.  
  The valid values are as follows:
  + **random**: Indicates random sampling of the logs, processing only part of the data.
  + **none**: Indicates random sampling is disabled and all data is processed.

* `ratio` - (Required, String) Specifies the sampling rate of the log.
  + If `sampler.type` is set to **random**, the valid value ranges from `0.1` to `1`.
  + If `sampler.type` is set to **none**, the value is set to `1`.

<a name="metric_rule_sinks"></a>
The `sinks` block supports:

* `type` - (Required, String) Specifies the type of the stored object.  
  The valid values are as follows:
  + **aom**

* `metric_name` - (Required, String) Specifies the name of the generated log metric. The name must be unique.
  The name only English letters, digits, hyphens (-) and colon(:) are allowed, and must start with an English letter.

* `name` - (Optional, String) Specifies the name of the AOM Prometheus common instance.  
  This parameter is required and available only when the `sinks.type` parameter is set to **aom**.

* `instance_id` - (Optional, String) Specifies the ID of the AOM Prometheus common instance.  
  This parameter is required and available only when the `sinks.type` parameter is set to **aom**.

<a name="metric_rule_aggregator"></a>
The `aggregator` block supports:

* `type` - (Required, String) Specifies the type of the log statistics.  
  The valid values are as follows:
  + **count**: Indicates the number of the logs.
  + **countKeyword**: Indicates the number of times the keyword appears.
  + **max**: Indicates the maximum value ​​of the specified field.
  + **min**: Indicates the minimum value ​​of the specified field.
  + **avg**: Indicates the average value ​​of the specified field.
  + **sum**: Indicates the sum value ​​of the specified field.
  + **p50**: Indicates the value below which `50%` of the data falls.
  + **p70**: Indicates the value below which `75%` of the data falls.
  + **p90**: Indicates the value below which `90%`of the data falls.
  + **p95**: Indicates the value below which `95%` of the data falls.
  + **p99**: Indicates the value below which `99%` of the data falls.

* `field` - (Required, String) Specifies the field of the log statistics.

* `group_by` - (Optional, List) Specifies the list of the group fields of the log statistics.

* `keyword` - (Optional, String) Specifies the keyword of the log statistics. The keyword is case sensitive.  
  This parameter is required and available only when the `aggregator.type` parameter is set to **countKeyword**.

<a name="metric_rule_filter"></a>
The `filter` block supports:

* `type` - (Optional, String) Specifies the filter type of the log.  
  The parameter must be used together with `filter.filters`.  
  The valid values are as follows:
  + **or**
  + **and**

* `filters` - (Optional, List) Specifies the list of log filtering rule groups.
  The [filters](#metric_rule_filter_groups) structure is documented below.

<a name="metric_rule_filter_groups"></a>
The `filters` block supports:

* `type` - (Optional, String) Specifies the filter type of the log.  
  The parameter must be used together with `filter.filters.filters`.  
  The valid values are as follows:
  + **or**
  + **and**

* `filters` - (Optional, List) Specifies the list of the log filter rule associations.
  The [filters](#metric_rule_associated_filters) structure is documented below.

<a name="metric_rule_associated_filters"></a>
The `filters` block supports:

* `key` - (Required, String) Specifies the filter field of the log.

* `type` - (Required, String) Specifies the filter conditions of the log.  
  The valid values are as follows:
  + **contains**: Applicable to `string` data type.
  + **notContains**: Applicable to `string` data type.
  + **fieldExist**: Applicable to `string` data type.
  + **fieldNotExist**: Applicable to `string`, `float` and `long` data types.
  + **equal**: Applicable to `string`, `float` and `long` data types.
  + **notEqual**: Applicable to `string`, `float` and `long` data types.
  + **gt**: Applicable to `float` and `long` data types.
  + **gte**: Applicable to `float` and `long` data types.
  + **lt**: Applicable to `float` and `long` data types.
  + **lte**: Applicable to `float` and `long` data types.
  + **range**: Applicable to `float` and `long` data types.
  + **outRange**: Applicable to `float` and `long` data types.

* `value` - (Optional, String) Specifies the value corresponding to the log filter field.
  This parameter is required and available only when the `filters.filters.filters.type` parameter is set to **contains**,
  **notContains**, **equal**, **notEqual**, **gt**, **gte**, **lt** or **lte**.

* `lower` - (Optional, String) Specifies the minimum value corresponding to the log filter field.  
  This parameter is required and available only when the `filters.filters.filters.type` parameter is set to **range**
  or **outRange**.

* `upper` - (Optional, String) Specifies the maximum value corresponding to the log filter field.  
  This parameter is required and available only when the `filters.filters.filters.type` parameter is set to **range**
  or **outRange**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also log metric rule ID.

* `created_at` - The creation time of the log metric rule, in RFC3339 format.

## Import

The log metric rule resource can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_lts_metric_rule.test <id>
```
