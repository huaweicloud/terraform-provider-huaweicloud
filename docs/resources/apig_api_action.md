---
subcategory: "API Gateway"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_api_action"
description: |-
  Use this resource to operate API within HuaweiCloud.
---

# huaweicloud_apig_api_action

Use this resource to operate API within HuaweiCloud.

-> This resource is only a one-time action resource for operating API. Deleting this resource will not clear
   the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Publish the API to the specified environment

```hcl
variable "instance_id" {}
variable "api_id" {}
variable "env_id" {}

resource "huaweicloud_apig_api_action" "test" {
  instance_id = var.instance_id
  api_id      = var.api_id
  env_id      = var.env_id
  action      = "online"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the APIG instance to which the API belongs is
located.  
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the dedicated instance to which the API belongs.

* `api_id` - (Required, String, NonUpdatable) Specifies the ID of the API to be published.

* `env_id` - (Required, String, NonUpdatable) Specifies the ID of the environment to which the API will be published.

* `action` - (Required, String, NonUpdatable) Specifies the operation on the API will be performed.  
  The valid values are as follows:
  + **online**: Publish (Release) the API to the specified environment.
  + **offline**: Cancel the publish (release) of the API in the specified environment.

* `remark` - (Optional, String, NonUpdatable) Specifies the description of the publish action.  
  The description contain a maximum of `255` characters.
  Chinese characters must be in **UTF-8** or **Unicode** format.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `publish_id` - The ID of the publish record.

* `api_name` - The name of the API.

* `publish_time` - The time when the API was published, in UTC format.

* `version_id` - The version ID of the online API.
