---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_instance_upgrade"
description: |-
  Manages a GaussDB OpenGauss instance upgrade resource within HuaweiCloud.
---

# huaweicloud_gaussdb_opengauss_instance_upgrade

Manages a GaussDB OpenGauss instance upgrade resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_gaussdb_opengauss_instance_upgrade" "test" {
  instance_id  = var.instance_id
  upgrade_type = "inplace"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the GaussDB OpenGauss instance.

  Changing this parameter will create a new resource.

* `upgrade_type` - (Required, String, ForceNew) Specifies the instance upgrade type. Value options:
  + **inplace**: In-place upgrade
  + **grey**: Gray upgrade
  + **hotfix**: Hot patch update

  Changing this parameter will create a new resource.

* `upgrade_action` - (Optional, String, ForceNew) Specifies the instance upgrade action. Value options:
  + **upgradeAutoCommit**: Auto-commit
  + **upgrade**: Rolling upgrade
  + **commit**: Commit
  + **rollback**: Rollback

  -> **NOTE:** If `upgrade_type` is set to **inplace**, this parameter is optional.
  <br>If `upgrade_type` is set to **grey**, this parameter can be set to **upgradeAutoCommit**, **upgrade**, **commit**,
    or **rollback**.
  <br>If `upgrade_type` is set to **hotfix**, this parameter can be set to **upgradeAutoCommit** or **rollback**.

  Changing this parameter will create a new resource.

* `target_version` - (Optional, String, ForceNew) Specifies the target version that the instance will be upgraded to.
  In a hot patch update, multiple hot patch versions can be configured.

  Changing this parameter will create a new resource.

* `upgrade_shard_num` - (Optional, Int, ForceNew) Specifies the number of shards to be upgraded in a gray upgrade for
  a distributed instance. It is mandatory when `upgrade_action` is set to **upgrade**. The value cannot be greater than
  the number of shards that have not been upgraded.

  Changing this parameter will create a new resource.

* `upgrade_az` - (Optional, String, ForceNew) Specifies the AZ to be upgraded in a gray upgrade. It is mandatory when
  `upgrade_action` is set to **upgrade**. Multiple AZs can be upgraded at the same time. Use commas(,) to separate AZs.
  You cannot enter an AZ that does not belong to the instance.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals to the `instance_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
