---
subcategory: "Global Accelerator (GA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ga_endpoints"
description: |-
  Use this data source to get the list of endpoints.
---

# huaweicloud_ga_endpoints

Use this data source to get the list of endpoints.

## Example Usage

```hcl
variable "endpoint_group_id" {}
variable "endpoint_id" {}

data "huaweicloud_ga_endpoints" "test" {
  endpoint_group_id = var.endpoint_group_id
  endpoint_id       = var.endpoint_id
}
```

## Argument Reference

The following arguments are supported:

* `endpoint_group_id` - (Required, String) Specifies the ID of the endpoint group to which the endpoint belongs.

* `endpoint_id` - (Optional, String) Specifies the ID of the endpoint.

* `status` - (Optional, String) Specifies the status of the endpoint.
  The valid values are as follows:
  + **ACTIVE**: The status of the endpoint is normal operation.
  + **ERROR**: The status of the endpoint is error.

* `resource_id` - (Optional, String) Specifies the ID of the backend resource corresponding to the endpoint.

* `resource_type` - (Optional, String) Specifies the type of the backend resource corresponding to the endpoint.
  Currently only supported **EIP**.

* `health_state` - (Optional, String) Specifies the health status of the endpoint.
  The valid values are as follows:
  + **INITIAL**: The endpoint status is initializing.
  + **HEALTHY**: The endpoint status is normal.
  + **UNHEALTHY**: The endpoint status is abnormal.
  + **NO_MONITOR**: The endpoint status is not monitored.

* `ip_address` - (Optional, String) Specifies the IP address of the backend resource corresponding to the endpoint.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `endpoints` - The list of the endpoints.
  The [endpoints](#ga_endpoints) structure is documented below.

<a name="ga_endpoints"></a>
The `endpoints` block supports:

* `id` - The ID of the endpoint.

* `endpoint_group_id` - The ID of the endpoint group to which the endpoint belongs.

* `resource_id` - The ID of the backend resource corresponding to the endpoint.

* `resource_type` - The type of the backend resource corresponding to the endpoint.

* `status` - The status of the endpoint.

* `weight` - The weight of traffic distribution to the endpoint. The valid value ranges from `0` to `100`.

* `health_state` - The health status of the endpoint.

* `ip_address` - The IP address of the backend resource corresponding to the endpoint.

* `created_at` - The creation time of the endpoint.

* `updated_at` - The latest update time of the endpoint.

* `frozen_info` - The frozen details of cloud services or resources.
  The [frozen_info](#endpoints_frozen_info) structure is documented below.

<a name="endpoints_frozen_info"></a>
The `frozen_info` block supports:

* `status` - The status of a cloud service or resource.

* `effect` - The status of the resource after being forzen.

* `scene` - The service scenario.
