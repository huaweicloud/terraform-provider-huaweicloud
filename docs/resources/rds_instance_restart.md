---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_instance_restart"
description: |-
  Manages an RDS instance restart resource within HuaweiCloud.
---

# huaweicloud_rds_instance_restart

Manages an RDS instance restart resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_instance_restart" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS MySQL instance.

* `restart_server` - (Optional, Bool, NonUpdatable) Specifies whether to restart the VM. This parameter is supported with
  the SQL Server DB engine only.

* `forcible` - (Optional, Bool, NonUpdatable) Specifies whether to forcibly restart the instance. This parameter is
  supported with the SQL Server DB engine only. Forcible restart will forcibly interrupt uncommitted transactions.

* `delay` - (Optional, Bool, NonUpdatable) Specifies whether to restart the instance during maintenance window.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals to the `instance_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
