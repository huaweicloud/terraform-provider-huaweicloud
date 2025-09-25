---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_host_configurations"
description: |-
  Use this data source to get the list of CES host configurations.
---

# huaweicloud_ces_host_configurations

Use this data source to get the list of CES host configurations.

## Example Usage

```hcl
variable "type" {}
variable "from" {}
variable "to" {}
variable "dim_0" {}

data "huaweicloud_ces_host_configurations" "test" {
  namespace = "SYS.ECS"
  type      = var.type
  from      = var.from
  to        = var.to
  dim_0     = var.dim_0
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `namespace` - (Required, String) Specifies the namespace of the queried service.
  For details, see [namespace](https://support.huaweicloud.com/intl/en-us/api-ces/ces_03_0059.html#ces_03_0059).

* `type` - (Required, String) Specifies the event type. It can contain only letters, underscores (_), and hyphens (-).

* `from` - (Required, Int) Specifies the start time of the query. UTC timestamp in milliseconds.

* `to` - (Required, Int) Specifies the end time of the query. UTC timestamp in milliseconds.

* `dim_0` - (Required, String) Specifies the first monitoring dimension.
  The format of the metric dimension is **key,value**. For example **dim.0=instance_id,i-12345**.
  For details, see [namespace](https://support.huaweicloud.com/intl/en-us/api-ces/ces_03_0059.html#ces_03_0059).

* `dim_1` - (Optional, String) Specifies the second monitoring dimension.
  The format of the metric dimension is **key,value**. For example **dim.1=instance_id,i-12345**.
  For details, see [namespace](https://support.huaweicloud.com/intl/en-us/api-ces/ces_03_0059.html#ces_03_0059).

* `dim_2` - (Optional, String) Specifies the third monitoring dimension.
  The format of the metric dimension is **key,value**. For example **dim.2=instance_id,i-12345**.
  For details, see [namespace](https://support.huaweicloud.com/intl/en-us/api-ces/ces_03_0059.html#ces_03_0059).

* `dim_3` - (Optional, String) Specifies the fourth monitoring dimension.
  The format of the metric dimension is **key,value**. For example **dim.3=instance_id,i-12345**.
  For details, see [namespace](https://support.huaweicloud.com/intl/en-us/api-ces/ces_03_0059.html#ces_03_0059).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `datapoints` - Indicates the list of configuration information.

  The [datapoints](#datapoints_struct) structure is documented below.

<a name="datapoints_struct"></a>
The `datapoints` block supports:

* `type` - Indicates the type of event.

* `timestamp` - Indicates the time the incident was reported.

* `value` - Indicates the host configuration information.
