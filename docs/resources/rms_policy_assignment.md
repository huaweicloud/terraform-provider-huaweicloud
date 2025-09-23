---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_policy_assignment"
description: ""
---

# huaweicloud_rms_policy_assignment

Using this resource to assign the policy and evaluate HuaweiCloud resources.

## Example Usage

### Assign a built-in policy to check a specified instance by a flavor

```hcl
variable "policy_assignment_name" {}
variable "region_name" {}
variable "ecs_instance_id" {}
variable "compliant_flavor" {}

data "huaweicloud_rms_policy_definitions" "test" {
  name = "allowed-ecs-flavors"
}

resource "huaweicloud_rms_policy_assignment" "test" {
  name                 = var.policy_assignment_name
  description          = "An ECS is noncompliant if its flavor is not in the specified flavor list (filter by resource ID)."
  policy_definition_id = try(data.huaweicloud_rms_policy_definitions.test.definitions[0].id, "")
  status               = "Enabled"

  policy_filter {
    region            = var.region_name
    resource_provider = "ecs"
    resource_type     = "cloudservers"
    resource_id       = var.ecs_instance_id
  }

  parameters = {
    listOfAllowedFlavors = "[\"${var.compliant_flavor}\"]"
  }
}
```

### Assign a built-in policy to periodically check whether an OBS bucket is tracked by CTS

```hcl
variable "policy_assignment_name" {}
variable "bucket_name" {}

data "huaweicloud_rms_policy_definitions" "test" {
  name = "cts-obs-bucket-track"
}

resource "huaweicloud_rms_policy_assignment" "test" {
  name                 = var.policy_assignment_name
  description          = "An account is noncompliant if none of its CTS trackers track specified OBS buckets."
  period               = "Six_Hours"
  policy_definition_id = try(data.huaweicloud_rms_policy_definitions.test.definitions[0].id, "")
  status               = "Enabled"

  parameters = {
    trackBucket = "\"${var.bucket_name}\""
  }
}
```

### Assign a custom policy

```hcl
variable "policy_assignment_name" {}
variable "function_urn" {}
variable "function_version" {}
variable "rms_admin_trust_agency" {}

resource "huaweicloud_rms_policy_assignment" "test" {
  name        = var.policy_assignment_name
  description = "The ECS instances that do not conform to the custom function logic are considered non-compliant."
  status      = "Enabled"

  custom_policy {
    function_urn = "${var.function_urn}:${var.function_version}"
    auth_type    = "agency"
    auth_value   = {
      agency_name = "\"${var.rms_admin_trust_agency}\""
    }
  }

  parameters = {
    string_example = "\"string_value\""
    array_example  = "[\"array_element\"]"
    object_example = "{\"terraform_version\":\"1.xx.x\"}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Specifies the name of the policy assignment.  
  The valid length is limited from `1` to `64`.  
  Change this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description of the policy assignment, which contain maximum of
  `512` characters.

* `policy_definition_id` - (Optional, String) Specifies the ID of the built-in policy definition.  
  This parameter and `custom_policy` are alternative.

* `period` - (Optional, String) Specifies the period of the policy assignment.  
  The valid values are as follows:
  + **One_Hour**
  + **Three_Hours**
  + **Six_Hours**
  + **Twelve_Hours**
  + **TwentyFour_Hours**

  Most one of `period` and `policy_filter` can be configured.

* `policy_filter` - (Optional, List) Specifies the configuration used to filter resources.  
  The [object](#rms_policy_filter) structure is documented below.

-> If the `period` is configured, it means that the evaluation is performed periodically.
  If the `policy_filter` is configured, it means that the evaluation is performed on the specified resources through
  the filter. If neither parameter is configured, it means that the evaluation is performed on all resources under the
  account.

* `custom_policy` - (Optional, List) Specifies the configuration of the custom policy.  
  The [object](#rms_custom_policy) structure is documented below.

* `parameters` - (Optional, Map) Specifies the rule definition of the policy assignment.

* `status` - (Optional, String) Specifies the expect status of the policy.
  The valid values are **Enabled** and **Disabled**.

* `tags` - (Optional, Map)  Specifies the key/value pairs to associate with the policy assignment.

<a name="rms_policy_filter"></a>
The `policy_filter` block supports:

* `region` - (Optional, String) Specifies the name of the region to which the filtered resources belong.

* `resource_provider` - (Optional, String) Specifies the service name to which the filtered resources belong.

* `resource_type` - (Optional, String) Specifies the resource type of the filtered resources.

* `resource_id` - (Optional, String) Specifies the resource ID used to filter a specified resource.

* `tag_key` - (Optional, String) Specifies the tag name used to filter resources.  
  This parameter and `resource_id` are alternative.

* `tag_value` - (Optional, String) Specifies the tag value used to filter resources.  
  Required if `tag_key` is set.

<a name="rms_custom_policy"></a>
The `custom_policy` block supports:

* `function_urn` - (Required, String) Specifies the function URN used to create the custom policy.

* `auth_type` - (Required, String) Specifies the authorization type of the custom policy.

* `auth_value` - (Optional, Map) Specifies the authorization value of the custom policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the policy assignment.

* `type` - The type of the policy assignment.  
  The valid values are as follows:
  + **builtin**
  + **custom**

* `created_at` - The creation time of the policy assignment.

* `updated_at` - The latest update time of the policy assignment.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.

## Import

Policy assignments can be imported using their `id`, e.g.

```bash
$ terraform import huaweicloud_rms_policy_assignment.test 63f48e3762ce955981ab7e25
```
