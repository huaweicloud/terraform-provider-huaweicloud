---
subcategory: "Cloud Application Engine (CAE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cae_environment"
description: |-
  Manages an environment resource within HuaweiCloud.
---

# huaweicloud_cae_environment

Manages an environment resource within HuaweiCloud.

## Example Usage

```hcl
variable "environment_name" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "security_group_id" {}
variable "swr_organization_name" {}

resource "huaweicloud_cae_environment" "test" {
  name = var.environment_name

  annotations = {
    type              = "exclusive"
    vpc_id            = var.vpc_id
    subnet_id         = var.subnet_id
    security_group_id = var.security_group_id
    group_name        = var.swr_organization_name
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the environment is located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the environment.  
  The valid length is limited from `3` to `30`, only lowercase letters, digits and hyphens (-) are allowed.  
  The name must start with a lowercase letter and end with a lowercase letter or a digit.  
  Changing this creates a new resource.

* `annotations` - (Required, Map, ForceNew) Specifies the additional attributes of the environment.  
  Changing this creates a new resource.  
  The required keys are as follows:
  + **vpc_id**: The VPC ID bound to the environment.
  + **subnet_id**: The ID of the VPC subnet bound to the environment.
  + **group_name**: The SWR organization name bound to the environment.

  The optional keys are as follows:
  + **type**: The environment type. Currently, only **exclusive** is supported.
  + **security_group_id**: The ID of the security group bound to the environment.  
    If omitted, the CAE service will automatically create it.

  -> Deleting the resource does not delete the security group that the service automatically created.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the ID of the enterprise project to which the
  environment belongs.
  Changing this creates a new resource.  
  This parameter is only valid for enterprise users, if omitted, default enterprise project will be used.  

* `max_retries` - (Optional, Int) Specifies the maximum retry number in the **create** or **delete** operation.  
  Defaults to `0`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, in UUID format.

* `status` - The status of the environment.

* `created_at` - The creation time of the environment, in RFC3339 format.

* `updated_at` - The latest update time of the environment, in RFC3339 format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `delete` - Default is 20 minutes.

## Import

The environment can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_cae_environment.test <id>
```
