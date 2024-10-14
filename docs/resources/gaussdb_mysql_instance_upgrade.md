---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_instance_upgrade"
description: |-
  Manages a GaussDB MySQL instance upgrade resource within HuaweiCloud.
---

# huaweicloud_gaussdb_mysql_instance_upgrade

Manages a GaussDB MySQL instance upgrade resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_gaussdb_mysql_instance_upgrade" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the GaussDB MySQL instance. Changing this parameter
  will create a new resource.

* `delay` - (Optional, Bool, ForceNew) Specifies whether the instance will be upgraded with a delay. Value options:
  + **true**: The instance will be upgraded during the specified maintenance window.
  + **false(default)**: The instance will be upgraded immediately.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is `instance_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
