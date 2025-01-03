---
subcategory: "Global Accelerator (GA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ga_access_log"
description: |-
  Manages an access log resource within HuaweiCloud.
---

# huaweicloud_ga_access_log

Manages an access log resource within HuaweiCloud.

-> Currently, the GA flow log interconnects with LTS only in the **cn-north-4** region.

## Example Usage

```hcl
variable "resource_type" {}
variable "resource_id" {}
variable "log_group_id" {}
variable "log_stream_id" {}

resource "huaweicloud_ga_access_log" "test" {
  resource_type = var.resource_type
  resource_id   = var.resource_id
  log_group_id  = var.log_group_id
  log_stream_id = var.log_stream_id
}
```

## Argument Reference

The following arguments are supported:

* `resource_type` - (Required, String, ForceNew) Specifies the type of the resource to which the access log belongs.
  Currently, only **LISTENER** is supported.
  Changing this parameter will create a new resource.

* `resource_id` - (Required, String, ForceNew) Specifies the ID of the resource to which the access log belongs.
  Changing this parameter will create a new resource.

* `log_group_id` - (Required, String) Specifies the ID of the log group to which the access log belongs.

* `log_stream_id` - (Required, String) Specifies the ID of the log stream to which the access log belongs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of the access log.
  The valid values are as follows:
  + **ACTIVE**: The resource is running.
  + **PENDING**: The status is to be determined.
  + **ERROR**: Failed to create the resource.
  + **DELETING**: The resource is being deleted.

* `created_at` - The creation time of the access log, in RFC3339 format.

* `updated_at` - The latest update time of the access log, in RFC3339 format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

The resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ga_access_log.test <id>
```
