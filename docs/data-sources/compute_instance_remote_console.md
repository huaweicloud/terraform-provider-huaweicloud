---
subcategory: "Elastic Cloud Server (ECS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_instance_remote_console"
description: ""
---

# huaweicloud_compute_instance_remote_console

Use this data source to get an available HuaweiCloud ECS compute instance remote console.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_compute_instance_remote_console" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the instances.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ECS ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `protocol` - The protocol of the ECS compute instance remote console.

* `type` - The type of ECS compute instance remote console.

* `url` - The url of ECS compute instance remote console.
