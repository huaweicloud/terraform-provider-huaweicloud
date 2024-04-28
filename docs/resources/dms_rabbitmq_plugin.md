---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rabbitmq_plugin"
description: ""
---

# huaweicloud_dms_rabbitmq_plugin

Manage DMS RabbitMQ plugin resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_dms_rabbitmq_plugin" "test" {
  instance_id = var.instance_id
  name        = "rabbitmq_sharding"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the RabbitMQ instance.
  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the plugin.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attribute are exported:

* `id` - The resource ID.

* `enable` - Indicates whether the plugin is enabled.

* `running` - Indicates whether the plugin is running.

* `version` - Indicates the version of the plugin.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 50 minutes.
* `delete` - Default is 50 minutes.

## Import

The RabbitMQ plugin can be imported using the RabbitMQ `instance_id` and `name` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dms_rabbitmq_plugin.test <instance_id>/<name>
```
