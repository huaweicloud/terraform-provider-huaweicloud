---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_plugin_associate"
description: ""
---

# huaweicloud_apig_plugin_associate

Use this resource to bind the APIs to the plugin within HuaweiCloud.

~> Before binding the API(s), please make sure all APIs have been published, otherwise you will receive a service error.

-> If this resource was imported and no changes were deployed before deletion (a change must be triggered to apply the
   `api_ids` configured in the script), terraform will delete all bound APIs for current configured plugin in
   specified publish environment. Otherwise, terraform will only delete the bound API(s) managed by the last change.

## Example Usage

```hcl
variable "instance_id" {}
variable "plugin_id" {}
variable "publish_environment_id" {}
variable "bind_api_ids" {
  type = list(string)
}

resource "huaweicloud_apig_plugin_associate" "test" {
  instance_id = var.instance_id
  plugin_id   = var.plugin_id
  env_id      = var.publish_environment_id
  api_ids     = var.bind_api_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the plugin and the APIs are located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the dedicated instance to which the APIs and the
  plugin belong.  
  Changing this will create a new resource.

* `plugin_id` - (Required, String, ForceNew) Specifies the plugin ID for APIs binding.  
  Changing this will create a new resource.

* `env_id` - (Required, String, ForceNew) The environment ID where the API was published.
  Changing this will create a new resource.

* `api_ids` - (Required, List) Specifies the API IDs bound by the plugin.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource ID. The format is `<instance_id>/<plugin_id>/<env_id>`.

## Import

Associate resources can be imported using their related dedicated instance ID of plugin (`instance_id`), `plugin_id` and
`env_id`, separated by slashes, e.g.

```bash
$ terraform import huaweicloud_apig_plugin_associate.test <instance_id>/<plugin_id>/<env_id>
```
