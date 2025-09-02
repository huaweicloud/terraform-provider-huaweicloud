---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_api_batch_action"
description: |-
   Use this resource to batch operate APIs within HuaweiCloud.
---

# huaweicloud_apig_api_batch_action

Use this resource to batch operate APIs within HuaweiCloud.

-> This resource is only a one-time action resource for performing an operation with the API list. Deleting this
   resource will not clear the corresponding request record, but will only remove the resource information from the
   tfstate file.

## Example Usage

### Batch publish APIs to the specified environment

```hcl
variable "instance_id" {}
variable "env_id" {}
variable "api_ids" {
  type = list(string)
}

resource "huaweicloud_apig_api_batch_action" "publish" {
  instance_id = var.instance_id
  action      = "online"
  env_id      = var.env_id
  apis        = var.api_ids
}
```

### Batch unpublish APIs in the specified environment

```hcl
variable "instance_id" {}
variable "env_id" {}
variable "api_ids" {
  type = list(string)
}

resource "huaweicloud_apig_api_batch_action" "unpublish" {
  instance_id = var.instance_id
  action      = "offline"
  env_id      = var.env_id
  apis        = var.api_ids
}
```

### Batch publish APIs by group to the specified environment

```hcl
variable "instance_id" {}
variable "env_id" {}
variable "group_id" {}

resource "huaweicloud_apig_api_batch_action" "publish_group" {
  instance_id = var.instance_id
  action      = "online"
  env_id      = var.env_id
  group_id    = var.group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the APIG instance to which the API belongs is
located.  
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the dedicated instance to which the APIs
belong.  
  This parameter cannot be updated after creation.

* `action` - (Required, String, NonUpdatable) Specifies the action to perform on the APIs.  
  The valid values are as follows:
  + **online**: Publish (Release) the list of APIs to the specified environment
  + **offline**: Cancel the publish (release) the list of the APIs in the specified environment.

* `env_id` - (Required, String, NonUpdatable) Specifies the ID of the environment where the action will be performed.

* `apis` - (Optional, List, NonUpdatable) Specifies the list of API IDs to perform the action on.  
  Maximum `1,000` APIs per operation.

-> Only one of group_id or apis can be provided; they are mutually exclusive.

* `group_id` - (Optional, String, NonUpdatable) Specifies the ID of the API group.  
  If specified, all APIs in this group will be processed.

-> Only one of group_id or apis can be provided; they are mutually exclusive.

* `remark` - (Optional, String, NonUpdatable) Specifies the remark for the batch operation.  
  The description contain a maximum of `255` characters.
  Chinese characters must be in **UTF-8** or **Unicode** format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
