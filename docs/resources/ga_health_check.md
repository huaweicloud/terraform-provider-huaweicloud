---
subcategory: "Global Accelerator (GA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ga_health_check"
description: |-
  Manages a GA health check resource within HuaweiCloud.
---

# huaweicloud_ga_health_check

Manages a GA health check resource within HuaweiCloud.

## Example Usage

```hcl
variable "endpoint_group_id" {}

resource "huaweicloud_ga_health_check" "test" {
  endpoint_group_id = var.endpoint_group_id
  enabled           = true
  interval          = 10
  max_retries       = 5
  port              = 8001
  timeout           = 10
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `endpoint_group_id` - (Required, String, ForceNew) Specifies the endpoint group ID.

  Changing this parameter will create a new resource.

* `enabled` - (Required, Bool) Specifies whether to enable health check.

* `interval` - (Required, Int) Specifies the health check interval, in seconds.  
  The valid value is range from `1` to `60`.

* `max_retries` - (Required, Int) Specifies the maximum number of retries.  
  The valid value is range from `1` to `10`.  
  Specifies the number of consecutive health checks when the health check result of an endpoint changes
  from **HEALTHY** to **UNHEALTHY**, or from **UNHEALTHY** to **HEALTHY**.

* `port` - (Required, Int) Specifies the port used for the health check.  
  The valid value is range from `1` to `65,535`.

* `timeout` - (Required, Int) Specifies the timeout duration of the health check, in seconds.  
  The valid value is range from `1` to `60`.
  It is recommended that you set a value less than that of parameter **interval**.

* `protocol` - (Optional, String) Specifies the health check protocol.
  Only **TCP** supported for now. Defaults to **TCP**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - Indicates the provisioning status. The value can be one of the following:
  + **ACTIVE**: The resource is running.
  + **PENDING**: The status is to be determined.
  + **ERROR**: Failed to create the resource.
  + **DELETING**: The resource is being deleted.

* `created_at` - Indicates when the health check was configured.

* `updated_at` - Indicates when the health check was updated.

* `frozen_info` - The frozen details of cloud services or resources.
  The [frozen_info](#health_check_frozen_info) structure is documented below.

<a name="health_check_frozen_info"></a>
The `frozen_info` block supports:

* `status` - The status of a cloud service or resource.
  The valid values are as follows:
  + `0`: unfrozen/normal (The cloud service will recover after being unfrozen.)
  + `1`: frozen (Resources and data will be retained, but the cloud service cannot be used.)
  + `2`: deleted/terminated (Both resources and data will be cleared.)

* `effect` - The status of the resource after being forzen.
  The valid values are as follows:
  + `1` (default): The resource is frozen and can be released.
  + `2`: The resource is frozen and cannot be released.
  + `3`: The resource is frozen and cannot be renewed.

* `scene` - The service scenario.
  The valid values are as follows:
  + **ARREAR**: The cloud service is in arrears, including expiration of yearly/monthly resources and fee deduction
    failure of pay-per-use resources.
  + **POLICE**: The cloud service is frozen for public security.
  + **ILLEGAL**: The cloud service is frozen due to violation of laws and regulations.
  + **VERIFY**: The cloud service is frozen because the user fails to pass the real-name authentication.
  + **PARTNER**: A partner freezes their customer's resources.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The healthcheck can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ga_health_check.test <id>
```
