---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_full_dead_lock_switch"
description: |-
  Use this resource to enable or disable the full dead lock switch within HuaweiCloud.
---

# huaweicloud_das_full_dead_lock_switch

Use this resource to enable or disable the full dead lock switch within HuaweiCloud.

-> This resource is a one-time action resource for switching the full dead lock. Deleting this resource will not clear
   the corresponding request record, but will only remove the resource information from the tfstate file.

-> This resource only supports to switch the full dead lock switch of **MySQL** instances.

## Example Usage

### Disable full dead lock switch

```hcl
variable "mysql_instance_id" {}

resource "huaweicloud_das_full_dead_lock_switch" "test" {
  instance_id = var.mysql_instance_id
  switch_on   = false
}
```

### Enable full dead lock switch

```hcl
variable "mysql_instance_id" {}

resource "huaweicloud_das_full_dead_lock_switch" "test" {
  instance_id     = var.mysql_instance_id
  switch_on       = true
  retention_hours = 200
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the full dead lock switch is located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the instance.

* `switch_on` - (Required, Bool, NonUpdatable) Whether to enable the full dead lock switch.  
  The valid values are as follows:
  + **true**: Enable the full dead lock switch.
  + **false**: Disable the full dead lock switch.

* `retention_hours` - (Optional, Int, NonUpdatable) Specifies the retention hours of the full dead lock data.  
  The valid value is range from `168` to `720`. Defaults to `168`.

  -> The `retention_hours` is only valid when `switch_on` is set to **true**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
