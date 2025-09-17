---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_api_batch_plugins_associate"
description: |-
  Use this resource to bind the plugins to the API within HuaweiCloud.
---

# huaweicloud_apig_api_batch_plugins_associate

Use this resource to bind the plugins to the API within HuaweiCloud.

~> Before binding the plugins, please make sure the API has been published, otherwise you will receive a service error.

-> If this resource was imported and no changes were deployed before deletion (a increase change must be triggered to
   apply the `plugin_ids` configured in the script), terraform will delete all bound plugins for current configured API
   in specified publish environment. Otherwise, terraform will only delete the bound plugins managed by the last change.

## Example Usage

```hcl
variable "instance_id" {}
variable "published_api_id" {}
variable "published_env_id" {}
variable "plugin_ids" {
  type = list(string)
}

resource "huaweicloud_apig_api_batch_plugins_associate" "test" {
  instance_id = var.instance_id
  api_id      = var.published_api_id
  env_id      = var.published_env_id
  plugin_ids  = var.plugin_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the API and plugins are located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the dedicated instance to which the API and
  plugins belong.

* `api_id` - (Required, String, NonUpdatable) Specifies the ID of the API to be bound with plugins.

* `env_id` - (Required, String, NonUpdatable) Specifies the ID of the environment where the API is published.

* `plugin_ids` - (Required, List) Specifies the list of plugin IDs to be bound to the API.

  -> Only one plugin of each type can be bound to an API.
     <br>For example, when two CORS plugins are bound to one API, the later-bound plugin will overwrite the previous
     plugin of the same type, causing the list to always prompt changes.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, the format is `{instance_id}/{api_id}/{env_id}`.

## Import

All associated plugins for the specified API can be imported using the related instance ID, application ID and the
published environment ID, separated by the slashes, e.g.

```bash
terraform import huaweicloud_apig_api_batch_plugins_associate.test <instance_id>/<api_id>/<env_id>
```
