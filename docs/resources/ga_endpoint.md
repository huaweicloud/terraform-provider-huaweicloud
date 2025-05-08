---
subcategory: "Global Accelerator (GA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ga_endpoint"
description: |-
  Manages a GA endpoint resource within HuaweiCloud.
---

# huaweicloud_ga_endpoint

Manages a GA endpoint resource within HuaweiCloud.

## Example Usage

```hcl
variable "endpoint_group_id" {}
variable "resource_id" {}
variable "ip_address" {}

resource "huaweicloud_ga_endpoint" "test" {
  endpoint_group_id = var.endpoint_group_id
  resource_id       = var.resource_id
  ip_address        = var.ip_address
}
```

## Argument Reference

The following arguments are supported:

* `endpoint_group_id` - (Required, String, ForceNew) Specifies the ID of the endpoint group
  to which the endpoint belongs.

  Changing this parameter will create a new resource.

* `resource_id` - (Required, String, ForceNew) Specifies the endpoint ID, for example, EIP ID.

  Changing this parameter will create a new resource.

* `ip_address` - (Required, String, ForceNew) Specifies the IP address of the endpoint.

  Changing this parameter will create a new resource.

* `resource_type` - (Optional, String, ForceNew) Specifies the endpoint type.
  The value can be **EIP**. Defaults to **EIP**.

  Changing this parameter will create a new resource.

* `weight` - (Optional, Int) Specifies the weight of the endpoint based on which the listener distributes traffic.
  The value ranges from `0` to `100`. Defaults to `1`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - Indicates the provisioning status. The value can be one of the following:
  + **ACTIVE**: The resource is running.
  + **PENDING**: The status is to be determined.
  + **ERROR**: Failed to create the resource.
  + **DELETING**: The resource is being deleted.

* `health_state` - Indicates the health check result of the endpoint. The value can be one of the following:
  + **INITIAL**: Initial.
  + **HEALTHY**: Healthy.
  + **UNHEALTHY**: Unhealthy.
  + **NO_MONITOR**: Not monitored.

* `created_at` - Indicates when the endpoint was created.

* `updated_at` - Indicates when the endpoint was updated.

* `frozen_info` - The frozen details of cloud services or resources.
  The [frozen_info](#endpoint_frozen_info) structure is documented below.

<a name="endpoint_frozen_info"></a>
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

The endpoint can be imported using `endpoint_group_id`, `id`, separated by slashes, e.g.

```bash
$ terraform import huaweicloud_ga_endpoint.test <endpoint_group_id>/<id>
```
