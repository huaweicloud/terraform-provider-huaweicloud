---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ddm_instance_upgrade"
description: |-
  Manages a DDM instance upgrade resource within HuaweiCloud.
---

# huaweicloud_ddm_instance_upgrade

Manages a DDM instance upgrade resource within HuaweiCloud.

-> **NOTE:** Deleting instance upgrade is not supported. If you destroy a resource of instance upgrade, the resource is
only removed from the state, but it remains in the cloud. And the instance doesn't return to the state before upgrade.

## Example Usage

```hcl
variable instance_id {}
variable target_version {}

resource "huaweicloud_ddm_instance_upgrade" "test" {
  instance_id    = var.instance_id
  target_version = var.target_version
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of a DDM instance.
  Changing this creates a new resource.

* `target_version` - (Required, String, ForceNew) Specifies the target version.

## Attribute Reference

* `id` - The resource ID. The value is the instance ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
