---
subcategory: "Cloud Bastion Host (CBH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbh_change_instance_type"
description: |-
  Manages a resource to change single node CBH instance type within HuaweiCloud.
---

# huaweicloud_cbh_change_instance_type

Manages a resource to change single node CBH instance type within HuaweiCloud.

-> This resource is only a one-time action resource to change a CBH instance type. Deleting this resource
will not clear the corresponding change instance type, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "server_id" {}
variable "availability_zone" {}

resource "huaweicloud_cbh_change_instance_type" "test" {
  server_id         = var.server_id
  availability_zone = var.availability_zone
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to change a CBH instance type.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `server_id` - (Required, String, NonUpdatable) Specifies the ID of the single node CBH instance to change type.
  
  -> This field only supports passing the single node CBH instance ID. Passing the primary/standby instance ID may result
  in a timeout exception.

* `availability_zone` - (Optional, String, NonUpdatable) Specifies the availability zone of the single node CBH instance
  to change type.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (same as `server_id`).

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
