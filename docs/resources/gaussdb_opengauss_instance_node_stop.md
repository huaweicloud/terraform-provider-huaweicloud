---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_instance_node_stop"
description: |-
  Manages a GaussDB OpenGauss instance node stop resource within HuaweiCloud.
---

# huaweicloud_gaussdb_opengauss_instance_node_stop

Manages a GaussDB OpenGauss instance node stop resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}

resource "huaweicloud_gaussdb_opengauss_instance_node_stop" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the GaussDB OpenGauss instance. Changing this parameter
  will create a new resource.

* `node_id` - (Required, String, ForceNew) Specifies the ID of the GaussDB OpenGauss instance node that needs to be stopped.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals to the `node_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
