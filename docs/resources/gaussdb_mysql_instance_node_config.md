---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_instance_node_config"
description: |-
  Manages a GaussDB MySQL instance node config resource within HuaweiCloud.
---

# huaweicloud_gaussdb_mysql_instance_node_config

Manages a GaussDB MySQL instance node config resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}

resource "huaweicloud_gaussdb_mysql_instance_node_config" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
  name        = "test_name"
  priority    = 3
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the GaussDB MySQL instance. Changing this parameter
  will create a new resource.

* `node_id` - (Required, String, ForceNew) Specifies the ID of the GaussDB MySQL instance node. Changing this parameter
  will create a new resource.

* `name` - (Optional, String) Specifies the node name. The value can contain **4** to **128** characters. Only digits,
  letters, hyphens (-), and underscores (_) are allowed.

* `priority` - (Optional, Int) Specifies the fail-over priority. The value can be **-1** or any number ranging from **1**
  to **16**. If the value is a positive number, a smaller value indicates a higher priority. This priority determines the
  order in which read replicas are promoted when recovering from a primary node failure. Read replicas with the same
  priority have the same probability of being promoted to the new primary node. If the value is **-1**, the read replica
  does not participate in a fail-over. After the priority of a read replica is set to **-1**, ensure that a single-AZ
  instance still has at least one read replica or that the remaining nodes of a cross-AZ instance are in different AZs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is `node_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.

## Import

The GaussDB MySQL instance node config can be imported using the `instance_id` and `node_id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_gaussdb_mysql_instance_node_config.test <instance_id>/<node_id>
```
