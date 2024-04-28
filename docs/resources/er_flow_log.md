---
subcategory: "Enterprise Router (ER)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_er_flow_log"
description: ""
---

# huaweicloud_er_flow_log

Manages a flow log resource under the ER instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "log_group_id" {}
variable "log_stream_id" {}
variable "resource_id" {}
variable "flow_log_name" {}

resource "huaweicloud_er_flow_log" "test" {
  instance_id    = var.instance_id
  log_store_type = "LTS"
  log_group_id   = var.log_group_id
  log_stream_id  = var.log_stream_id
  resource_type  = "attachment"
  resource_id    = var.resource_id
  name           = var.flow_log_name
  description    = "Flow log created by terraform"
  enabled        = false
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the ER instance to which the flow log belongs.
  Changing this creates a new resource.

* `log_store_type` - (Required, String, ForceNew) Specifies the storage type of flow log. The valid value is **LTS**.  
  Changing this creates a new resource.

* `log_group_id` - (Required, String, ForceNew) Specifies the LTS log group ID.
  Changing this creates a new resource.

* `log_stream_id` - (Required, String, ForceNew) Specifies the LTS log stream ID.
  Changing this creates a new resource.

* `resource_type` - (Required, String, ForceNew) Specifies the resource type to which the logs to be collected.
  The valid value is **attachment**.
  Changing this creates a new resource.

* `resource_id` - (Required, String, ForceNew) Specifies the resource ID to which the logs to be collected.
  Changing this creates a new resource.

* `name` - (Required, String) Specifies the name of the flow log.

* `description` - (Optional, String) Specifies the description of the flow log.

* `enabled` - (Optional, Bool) Specifies whether to enable the flow log function. The default value is **true**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `state` - The current status of the flow log.

* `created_at` - The creation time of the flow log.

* `updated_at` - The latest update time of the flow log.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.
* `delete` - Default is 2 minutes.

## Import

The flow log can be imported using the related `instance_id` and their `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_er_flow_log.test <instance_id>/<id>
```
