---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_flavors"
description: ""
---

# huaweicloud_workspace_flavors

Use this data source to get the list of Workspace flavors.

## Example Usage

```hcl
variable "os_type" {}

data "huaweicloud_workspace_flavors" "test" {
  os_type = var.os_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `vcpus` - (Optional, Int) Specifies CPU specification of the Workspace flavors.

* `memory` - (Optional, Int) Specifies the memory size of the Workspace flavors, in GB.

* `os_type` - (Optional, String) Specifies the operating system type of the Workspace flavors.  
  The options are as follows:
  + **Windows**: The operating system type of the Workspace flavor is Windows.
  + **Linux**: The operating system type of the Workspace flavor is Linux.

* `availability_zone` - (Optional, String) Specifies the availability zone to which the Workspace flavors belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavors` - Flavor list.
  The [flavors](#workspace_flavors) structure is documented below.

<a name="workspace_flavors"></a>
The `flavors` block supports:

* `id` - Flavor ID.

* `type` - The type of the Workspace flavor.

* `architecture` - The Workspace flavor architecture, currently supporting arm and x86.
   The valid values are as follows:
  + **arm**: The Workspace flavor architecture is arm.
  + **x86**: The Workspace flavor architecture is x86.

* `vcpus` - CPU specifications of the Workspace flavor.

* `memory` - The Workspece flavor memory size in GB.

* `is_gpu` - The Workspace flavor is a specification of GPU type or not.

* `system_disk_type` - The system disk type of the Workspace flavor.

* `system_disk_size` - The system disk size of the Workspace flavor, in GB.

* `description` - Description of the Workspace flavor.

* `charging_mode` - Periodic package identification of the Workspace flavor.
  The valid values are as follows:
  + **postPaid**: Indicates on-demand billing of the Workspace flavor.

* `status` - The status of the Workspace flavor.  
  The valid values are as follows:
  + **normal**: The status is normal of the Workspace flavor.
