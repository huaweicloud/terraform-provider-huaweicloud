---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_instance_restart"
description: |-
  Manages a GaussDB MySQL instance restart resource within HuaweiCloud.
---

# huaweicloud_gaussdb_mysql_instance_restart

Manages a GaussDB MySQL instance restart resource within HuaweiCloud.

## Example Usage

### Restart GaussDB MySQL instance

```hcl
variable "instance_id" {}

resource "huaweicloud_gaussdb_mysql_instance_restart" "test" {
  instance_id = var.instance_id
}
```

### Restart GaussDB MySQL instance node

```hcl
variable "instance_id" {}
variable "node_id" {}

resource "huaweicloud_gaussdb_mysql_instance_restart" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the GaussDB MySQL instance. Changing this parameter
  will create a new resource.

* `node_id` - (Optional, String, ForceNew) Specifies the node ID of the GaussDB MySQL instance. Changing this parameter
  will create a new resource.

  -> **NOTE:** If `node_id` is not specified, then the GaussDB MySQL instance will be rebooted.
    <br/> If `node_id` is specified, then the node of the GaussDB MySQL instance will be rebooted.

* `delay` - (Optional, Bool, ForceNew) Specifies whether the instance/node will be rebooted with a delay. Value options:
  + **true**: The instance/node will be rebooted during the specified maintenance window.
  + **false(default)**: The instance/node will be rebooted immediately.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. If an instance is rebooted, the value is `instance_id`. If a node is rebooted, the value is `node_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
