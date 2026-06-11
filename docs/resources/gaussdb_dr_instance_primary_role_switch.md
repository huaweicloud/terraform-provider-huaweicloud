---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_dr_instance_primary_role_switch"
description: |-
  Manages a GaussDB DR instance primary role switch resource within HuaweiCloud.
---

# huaweicloud_gaussdb_dr_instance_primary_role_switch

Manages a GaussDB DR instance primary role switch resource within HuaweiCloud.

-> This resource is a one-time action resource for switching the primary role between primary and DR instances.
   Deleting this resource will not revert the operation, but will only remove the resource information from the
   tfstate file.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_gaussdb_dr_instance_primary_role_switch" "test" {
  instance_id   = var.instance_id
  disaster_type = "stream"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to operate the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String, NonUpdatable) Specifies the instance ID to switch the primary role.

* `disaster_type` - (Required, String, NonUpdatable) Specifies the disaster recovery type.
  The value can be **stream** (streaming disaster recovery).

* `post_process_config` - (Optional, String, NonUpdatable) Specifies whether to support automatic recovery when
  switchover fails. This field is only supported for database engine version V2.0-8.200 and later with Quorum streaming
  disaster recovery. Other disaster recovery scenarios do not support this feature.
  The valid values are:
  + **AUTO**: Automatic recovery when switchover fails.
  + **MANUAL**: No automatic recovery when switchover fails.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals the `instance_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
