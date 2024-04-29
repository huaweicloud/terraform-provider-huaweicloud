---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rabbitmq_plugins"
description: ""
---

# huaweicloud_dms_rabbitmq_plugins

Use this data source to get the list of DMS RabbitMQ plugins.

## Example Usage

```hcl
var "instance_id" {}

data "huaweicloud_dms_rabbitmq_plugins" "test" {
  instance_id  = var.instance_id
  name         = "rabbitmq_mqtt"
  enable       = true
  running      = true
  version      = "3.8.35"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RabbitMQ instance.

* `name` - (Optional, String) Specifies the name of the plugin.

* `enable` - (Optional, Bool) Specifies whether the plugin is enabled. Defaults to **false**.

* `running` - (Optional, Bool) Specifies whether the plugin is running. Defaults to **false**.

* `version` - (Optional, String) Specifies the version of the plugin.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `plugins` - The list of the RabbitMQ plugins.
  The [plugins](#DMS_rabbitmq_plugins) structure is documented below.

<a name="DMS_rabbitmq_plugins"></a>
The `plugins` block supports:

* `name` - Indicates the name of the plugin.

* `enable` - Indicates whether the plugin is enabled.

* `running` - Indicates whether the plugin is running.

* `version` - Indicates the version of the plugin.
