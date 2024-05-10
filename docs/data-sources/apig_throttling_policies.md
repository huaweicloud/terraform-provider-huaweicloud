---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_api_associated_throttling_policies"
description: |-
  Use this data source to get the list of the throttling policies under the APIG instance within HuaweiCloud.
---

# huaweicloud_apig_throttling_policies

Use this data source to get the list of the throttling policies under the APIG instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_apig_throttling_policies" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the throttling policies belong.

* `name` - (Optional, String) Specifies the name of the throttling policy. Fuzzy search is supported.

* `policy_id` - (Optional, String) Specifies the ID of the throttling policy.

* `type` - (Optional, String) The type of the throttling policy.
  The valid values are as follows:
  + **API-based**: Limiting the maximum number of times a single API bound to the policy can be called within the
    specified period.
  + **API-shared**: Limiting the maximum number of times all APIs bound to the policy can be called within the specified
    period.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `policies` - All throttling policies that match the filter parameters.
  The [policies](#throttling_policies) structure is documented below.

<a name="throttling_policies"></a>
The `policies` block supports:

* `id` - The ID of the throttling policy.

* `name` - The name of the throttling policy.

* `period` - The period of time for limiting the number of API calls.

* `period_unit` - The time unit for limiting the number of API calls.
  The valid values are **SECOND**, **MINUTE**, **HOUR** and **DAY**.

* `max_api_requests` - The maximum number of times an API can be accessed within a specified period.

* `max_app_requests` - The maximum number of times the API can be accessed by an app within the same period.

* `max_ip_requests` - The maximum number of times the API can be accessed by an IP address within the same period.

* `max_user_requests` - The maximum number of times the API can be accessed by a user within the same period.

* `type` - The type of the throttling policy.

* `user_throttles` - The array of one or more special throttling policies for IAM user limit.
  The [user_throttles](#special_throttling_policies) structure is documented below.

* `app_throttles` - The array of one or more special throttling policies for APP limit.
  The [app_throttles](#special_throttling_policies) structure is documented below.

* `bind_num` - The number of APIs bound to the throttling policy.

* `description` - The description of throttling policy.

* `created_at` - The creation time of the throttling policy, in RFC3339 format.

<a name="special_throttling_policies"></a>
The `user_throttles` and `app_throttles` blocks support:

* `id` - The ID of the special user/application throttling policy.

* `max_api_requests` - The maximum number of times an API can be accessed within a specified period.

* `throttling_object_id` - The object ID which the special user/application throttling policy belongs.

* `throttling_object_name` - The object name which the special user/application throttling policy belongs.
