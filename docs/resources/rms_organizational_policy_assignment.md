---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_organizational_policy_assignment"
description: ""
---

# huaweicloud_rms_organizational_policy_assignment

Using this resource to assign the organizational policy HuaweiCloud resources.

## Example Usage

### Assign a built-in policy

```hcl
variable "policy_assignment_name" {}

data "huaweicloud_organizations_organization" "test" {}

data "huaweicloud_rms_policy_definitions" "test" {
  name = "access-keys-rotated"
}

resource "huaweicloud_rms_organizational_policy_assignment" "test" {
  organization_id      = data.huaweicloud_organizations_organization.test.id
  name                 = var.policy_assignment_name
  description          = "The maximum number of days without rotation. Default 90."
  policy_definition_id = try(data.huaweicloud_rms_policy_definitions.test.definitions[0].id, "")
  period               = "TwentyFour_Hours"

  parameters = {
    maxAccessKeyAge = "\"90\""
  }
}
```

### Assign a custom policy

```hcl
variable "function_name" {}
variable "policy_assignment_name" {}

data "huaweicloud_organizations_organization" "test" {}

resource "huaweicloud_fgs_function" "test" {
  name                  = var.function_name
  code_type             = "inline"
  handler               = "index.handler"
  runtime               = "Node.js10.16"
  functiongraph_version = "v2"
  app                   = "default"
  enterprise_project_id = "0"
  memory_size           = 128
  timeout               = 3
}

resource "huaweicloud_rms_organizational_policy_assignment" "test" {
  organization_id = data.huaweicloud_organizations_organization.test.id
  name            = var.policy_assignment_name
  description     = "This is a custom policy assignment."
  function_urn    = "${huaweicloud_fgs_function.test.urn}:${huaweicloud_fgs_function.test.version}"
  period          = "TwentyFour_Hours"

  parameters = {
    string_test = "\"string_value\""
    array_test  = "[\"array_element\"]"
    object_test = jsonencode({"terraform_version": "1.xx.x"})
  }
}
```

## Argument Reference

The following arguments are supported:

* `organization_id` - (Required, String, ForceNew) Specifies the ID of the organization.  
  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the organizational policy assignment.  
  The valid length is limited from `1` to `60`, only letters, digits, hyphens (-) and underscores (_) are allowed.  
  Change this parameter will create a new resource.

* `excluded_accounts` - (Optional, List) Specifies the excluded accounts of the organizational policy assignment.

* `description` - (Optional, String) Specifies the description of the organizational policy assignment,
  which contain maximum of `512` characters.

* `policy_definition_id` - (Optional, String, ForceNew) Specifies the ID of the built-in policy definition.  
  This parameter and `function_urn` are alternative.  
  Changing this parameter will create a new resource.

* `function_urn` - (Optional, String, ForceNew) Specifies the function URN used to create the custom policy.  
  Changing this parameter will create a new resource.

* `period` - (Optional, String) Specifies the period of the organizational policy assignment.  
  The valid values are as follows:
  + **One_Hour**
  + **Three_Hours**
  + **Six_Hours**
  + **Twelve_Hours**
  + **TwentyFour_Hours**

  Most one of `period` and `policy_filter` can be configured.

* `policy_filter` - (Optional, List) Specifies the configuration used to filter resources.  
  The [policy_filter](#rms_policy_filter) structure is documented below.

* `parameters` - (Optional, Map) Specifies the rule definition of the organizational policy assignment.

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

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the organizational policy assignment.

* `owner_id` - Indicates the creator of the organizational policy assignment.

* `organization_policy_assignment_urn` - Indicates the unique identifier of the organizational policy assignment.

* `created_at` - The creation time of the organizational policy assignment.

* `updated_at` - The latest update time of the organizational policy assignment.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.

* `update` - Default is 30 minutes.

* `delete` - Default is 30 minutes.

## Import

The organizational policy assignment can be imported using the `organization_id` and `id`separated by a slash, e.g.

```bash
$ terraform import huaweicloud_rms_organizational_policy_assignment.test <organization_id>/<id>
```
