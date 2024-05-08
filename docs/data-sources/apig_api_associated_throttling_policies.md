---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_api_associated_throttling_policies"
description: |-
  Use this data source to query the throttling policies associated with the specified API within HuaweiCloud.
---

# huaweicloud_apig_api_associated_throttling_policies

Use this data source to query the throttling policies associated with the specified API within HuaweiCloud.

## Example Usage

### Query the contents of all throttling policies bound to the current API

```hcl
variable "instance_id" {}
variable "associated_api_id" {}

data "huaweicloud_apig_api_associated_throttling_policies" "test" {
  instance_id = var.instance_id
  api_id      = var.associated_api_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the associated throttling policies.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the throttling policies belong.

* `api_id` - (Required, String) Specifies the ID of the API bound to the throttling policy.

* `policy_id` - (Optional, String) Specifies the ID of the throttling policy.

* `name` - (Optional, String) Specifies the name of the throttling policy.

* `type` - (Optional, String) Specifies the type of the throttling policy.  
  The valid values are as follows:
  + **API-based**: limiting the maximum number of times a single API bound to the policy can be called within the
    specified period.
  + **API-shared**: limiting the maximum number of times all APIs bound to the policy can be called within the specified
    period.

* `env_name` - (Optional, String) Specifies the name of the environment where the API is published.

* `period_unit` - (Optional, String) Specifies the time unit for limiting the number of API calls.  
  The valid values are **SECOND**, **MINUTE**, **HOUR** and **DAY**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `policies` - All throttling policies that match the filter parameters.
  The [policies](#api_associated_throttling_policies) structure is documented below.

<a name="api_associated_throttling_policies"></a>
The `policies` block supports:

* `id` - The ID of the throttling policy.

* `name` - The name of the throttling policy.

* `type` - The type of the throttling policy.

* `period_unit` - The time unit for limiting the number of API calls.

* `period` - The period of time for limiting the number of API calls.

* `max_api_requests` - The maximum number of times an API can be accessed within a specified period.

* `max_app_requests` - The maximum number of times the API can be accessed by an app within the same period.

* `max_ip_requests` - The maximum number of times the API can be accessed by an IP address within the same period.

* `max_user_requests` - The maximum number of times the API can be accessed by a user within the same period.

* `env_name` - The name of the environment where the API is published.

* `description` - The description of the throttling policy.

* `user_throttles` - The array of one or more special throttling policies for IAM user limit.
  The [user_throttles](#throttling_policy_rule_detail_attr) structure is documented below.

* `app_throttles` - The array of one or more special throttling policies for APP limit.
  The [app_throttles](#throttling_policy_rule_detail_attr) structure is documented below.

* `bind_id` - The bind ID.

* `bind_time` - The time that the throttling policy is bound to the API, in RFC3339 format.

* `created_at` - The creation time of the throttling policy, in RFC3339 format.

<a name="throttling_policy_rule_detail_attr"></a>
The `user_throttles` and `app_throttles` blocks support:

* `max_api_requests` - The maximum number of times an API can be accessed within a specified period.

* `throttling_object_id` - The object ID which the special user/application throttling policy belongs.

* `throttling_object_name` - The object name which the special user/application throttling policy belongs.

* `id` - The ID of the special user/application throttling policy.
