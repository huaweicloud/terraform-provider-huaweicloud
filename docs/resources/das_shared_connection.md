---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_shared_connection"
description: |-
  Manages DAS shared connection resource within HuaweiCloud.
---

# huaweicloud_das_shared_connection

Manages DAS shared connection resource within HuaweiCloud.

## Example Usage

```hcl
variable "connection_id" {}
variable "user_id" {}
variable "user_name" {}

resource "huaweicloud_das_shared_connection" "test" {
  connection_id = var.connection_id
  user_id       = var.user_id
  user_name     = var.user_name
  expired_at    = "2006-01-02T15:04:05+08:00"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the shared connection is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `connection_id` - (Required, String, NonUpdatable) Specifies the ID of the connection to which the shared connection belongs.

* `user_id` - (Required, String, NonUpdatable) Specifies the user ID of the shared connection.

* `user_name` - (Required, String, NonUpdatable) Specifies the user name of the shared connection.

* `expired_at` - (Optional, String, NonUpdatable) Specifies the expiration time of the shared connection, in RFC3339 format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in the format of `<connection_id>/<user_id>`.

* `shared_at` - The creation time of the shared connection, in RFC3339 format.

## Import

The shared connection can be imported using `connection_id` and `user_id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_das_shared_connection.test <connection_id>/<user_id>
```
