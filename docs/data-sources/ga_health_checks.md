---
subcategory: "Global Accelerator (GA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ga_health_checks"
description: |-
  Use this data source to get the list of health checks.
---

# huaweicloud_ga_health_checks

Use this data source to get the list of health checks.

## Example Usage

```hcl
variable "health_check_id" {}

data "huaweicloud_ga_health_checks" "test" {
  health_check_id = var.health_check_id
}
```

## Argument Reference

The following arguments are supported:

* `health_check_id` - (Optional, String) Specifies the ID of the health check.

* `endpoint_group_id` - (Optional, String) Specifies the ID of the endpoint group to which the health check belongs.

* `status` - (Optional, String) Specifies the status of the health check.
  The valid values are as follows:
  + **ACTIVE**: The status of the health check is normal operation.
  + **ERROR**: The status of the health check is error.

* `protocol` - (Optional, String) Specifies the front end protocol of the health check used.
  Currently only supported **TCP**.

* `enabled` - (Optional, String) Specifies whether health check is enabled.
  The value can be **true** and **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `health_checks` - The list of the health checks.
  The [health_checks](#ga_health_checks) structure is documented below.

<a name="ga_health_checks"></a>
The `health_checks` block supports:

* `id` - The ID of the health check.

* `endpoint_group_id` - The ID of the endpoint group to which the health check belongs.

* `status` - The status of the health check.

* `protocol` - The front end protocol of the health check used.

* `port` - The port of the health check.

* `interval` - The time interval of the health check. The unit is seconds, the valid value ranges from `1` to `60`.

* `timeout` - The timeout of the health check. The unit is seconds, The valid value ranges from `1` to `60`.

* `max_retries` - The max retries of the health check. The valid value ranges from `1` to `10`.

* `enabled` - Whether health check is enabled.

* `created_at` - The creation time of the health check.

* `updated_at` - The latest update time of the health check.

* `frozen_info` - The frozen details of cloud services or resources.
  The [frozen_info](#health_checks_frozen_info) structure is documented below.

<a name="health_checks_frozen_info"></a>
The `frozen_info` block supports:

* `status` - The status of a cloud service or resource.

* `effect` - The status of the resource after being forzen.

* `scene` - The service scenario.
