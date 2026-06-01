---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_lock_blocking_switch"
description: |-
  Use this resource to enable or disable the lock blocking switch within HuaweiCloud.
---

# huaweicloud_das_lock_blocking_switch

Use this resource to enable or disable the lock blocking switch within HuaweiCloud.

-> This resource is a one-time action resource for switching the lock blocking. Deleting this resource will not clear
   the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Disable lock blocking switch

```hcl
variable "sqlserver_instance_id" {}

resource "huaweicloud_das_lock_blocking_switch" "test" {
  instance_id = var.sqlserver_instance_id
  switch_on   = false
}
```

### Enable lock blocking switch with retention hours

```hcl
variable "sqlserver_instance_id" {}

resource "huaweicloud_das_lock_blocking_switch" "test" {
  instance_id     = var.sqlserver_instance_id
  switch_on       = true
  retention_hours = 200
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the lock blocking switch is located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the instance.

  -> The `instance_id` only supports **SQLServer** instance.

* `switch_on` - (Required, Bool, NonUpdatable) Whether to enable the lock blocking switch.  
  The valid values are as follows:
  + **true**: Enable the lock blocking switch.
  + **false**: Disable the lock blocking switch.

* `retention_hours` - (Optional, Int, NonUpdatable) Specifies the retention hours of the lock blocking data.  
  The valid value is range from `168` to `720`. Defaults to `168`.

  -> The `retention_hours` is only valid when `switch_on` is set to **true**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
