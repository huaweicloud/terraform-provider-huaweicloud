---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_metric_data_add"
description: |-
  Manages a CES metric data add resource within HuaweiCloud.
---

# huaweicloud_ces_metric_data_add

Manages a CES metric data add resource within HuaweiCloud.

-> This resource is only a one-time action resource for operating the API.
Deleting this resource will not clear the corresponding request record,
but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "namespace" {}
variable "metric_name" {}
variable "name" {}
variable "value" {}
variable "ttl" {}
variable "collect_time" {}
variable "unit" {}
variable "type" {}

resource "huaweicloud_ces_metric_data_add" "test" {
  metric {
    namespace   = var.namespace
    metric_name = var.metric_name

    dimensions {
      name  = var.name
      value = var.value
    }
  }

  ttl          = var.ttl
  collect_time = var.collect_time
  value        = var.value
  unit         = var.unit
  type         = var.type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `metric` - (Required, List, NonUpdatable) Specifies the custom CES monitoring metric data.

  The [metric](#Metric) structure is documented below.

* `ttl` - (Required, Int, NonUpdatable) Specifies the monitoring metric data retention period.
  The unit is second. The range of values is from **1** to **604800**.

* `collect_time` - (Required, String, NonUpdatable) Specifies the collect time.
  The time is in UTC. The format is **yyyy-MM-dd HH:mm:ss**.

* `value` - (Required, Float, NonUpdatable) Specifies the value of the monitoring metric data.
  The value can be an integer or a floating point number.

* `unit` - (Optional, String, NonUpdatable) Specifies the unit of the monitoring metric data.

* `type` - (Optional, String, NonUpdatable) Specifies the type of the monitoring metric data.
  The valid value can be **int** or **float**.

<a name="Metric"></a>
The `metric` block supports:

* `namespace` - (Required, String, NonUpdatable) Specifies the customized namespace.
  The namespace must be in the **service.item** format and contain `3` to `32` characters.
  **service** and **item** each must start with a letter and contain only letters, digits, and underscores (_).
  In addition, **service** cannot start with **SYS**, **AGT**, or **SRE**. The namespace cannot be **SERVICE.BMS**,
  because this namespace has been used by the system.

* `dimensions` - (Required, List, NonUpdatable) Specifies the metric dimension.
  A maximum of four dimensions are supported.
  
  The [dimensions](#Dimensions) structure is documented below.

* `metric_name` - (Required, String, NonUpdatable) Specifies the metric ID.

<a name="Dimensions"></a>
The `dimensions` block supports:

* `name` - (Required, String, NonUpdatable) Specifies the dimension.

* `value` - (Required, String, NonUpdatable) Specifies the dimension value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
