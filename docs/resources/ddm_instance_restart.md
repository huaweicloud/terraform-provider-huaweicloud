---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ddm_instance_restart"
description: |-
  Manages a DDM instance restart resource within HuaweiCloud.
---

# huaweicloud_ddm_instance_restart

Manages a DDM instance restart resource within HuaweiCloud.

## Example Usage

```hcl
variable instance_id {}

resource "huaweicloud_ddm_instance_restart" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of a DDM instance.
  Changing this creates a new resource.

* `type` - (Optional, String, ForceNew) Specifies the restart type. Value options:
  + **soft**: Only the process is restarted.
  + **hard**: The instance VM is forcibly restarted.

## Attribute Reference

* `id` - The resource ID. The value is the instance ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
