---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ddm_instance_rollback"
description: |-
  Manages a DDM instance rollback resource within HuaweiCloud.
---

# huaweicloud_ddm_instance_rollback

Manages a DDM instance rollback resource within HuaweiCloud.

-> **NOTE:** Deleting instance rollback is not supported. If you destroy a resource of instance rollback, the resource is
only removed from the state, but it remains in the cloud. And the instance doesn't return to the state before rollback.

## Example Usage

```hcl
variable instance_id {}

resource "huaweicloud_ddm_instance_rollback" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of a DDM instance.
  Changing this creates a new resource.

## Attribute Reference

* `id` - The resource ID. The value is the instance ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
