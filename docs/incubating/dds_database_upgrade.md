---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_database_upgrade"
description: |-
  Manages a reosurce to upgrade database patch within HuaweiCloud.
---

# huaweicloud_dds_database_upgrade

Manages a reosurce to upgrade database patch within HuaweiCloud.

-> 1. This resource is a one-time action resource. Deleting this resource will not clear the corresponding request record,
   but will only remove the resource information from the tf state file.
   <br/>2. This resource is not available to frozen or abnormal instances.
   <br/>3. This resource is not available if there are abnormal instance nodes.
   <br/>4. View field `patch_available` in the result returned by the API for querying instance details and check
   whether a minor version upgrade is supported.
   <br/>5. Perform an upgrade during off-peak hours.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_dds_database_upgrade" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.  
  If omitted, the provider-level region will be used.  
  Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the DDS instance ID.

* `upgrade_mode` - (Optional, String, NonUpdatable) Specifies the upgrade mode.  
  The valid values are as follows:
  + **minimized_interrupt_time**: The upgrade with the shortest interruption time is preferred.
    In this mode, the upgrade has little impact on services.
  + **minimized_upgrade_time**: The upgrade with the shortest upgrade time is preferred.

  The default value is **minimized_interrupt_time**.

* `is_delayed` - (Optional, Bool, NonUpdatable) Specifies whether the instance is automatically upgraded
  during the maintenance window.  
  The valid values are as follows:
  + **true**: The instance will be upgraded during the specified maintenance window.
  + **false**: The instance will be upgraded immediately.

  The default value is **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
