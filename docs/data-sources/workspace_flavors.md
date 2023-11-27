---
subcategory: "Workspace Flavors Service (Workspace)"
---

# huaweicloud_workspace_flavors

Use this data source to get the list of Workspace flavor.

## Example Usage

```hcl
variable "os_type" {}

data "huaweicloud_workspaces_flavors" "test" {
  os_type = var.os_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `vcpus` - (Optional, Int) Specifies CPU specification of the workspace flavors.

* `memory` - (Optional, Int) Specifies the workspace flavors memory size in GB.

* `os_type` - (Optional, String) Specifies the operating system type of the workspace flavors.  
  The options are as follows:
  + **Windows**: The operating system type of the workspace flavor is Windows.
  + **Linux**: The operating system type of the workspace flavor is Linux.

* `availability_zone` - (Optional, String) Specifies the availability zone to which the workspace flavors belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavors` - Flavor list.
  The [flavors](#workspace_flavors) structure is documented below.

<a name="workspace_flavors"></a>
The `flavors` block supports:

* `id` - Flavor ID.

* `type` - The type of the workspace flavor.

* `architecture` - The workspace flavor architecture, currently supporting arm and x86.
   The valid values are as follows:
  + **arm**: The flavor architecture is arm.
  + **x86**: The flavor architecture is x86.

* `vcpus` - CPU specifications of the workspace flavor.

* `memory` - The workspece flavor memory size in GB.

* `is_gpu` - The workspace flavor is a specification of GPU type or not.

* `system_disk_type` - The workspace flavor system disk type.

* `system_disk_size` - The workspace flavor system disk size in GB.

* `description` - Description of the workspace flavor.

* `charging_mode` - Periodic package identification of the workspace flavor.
  The valid values are as follows:
  + **postPaid**: Indicates on-demand billing of the workspace flavor.

* `status` - The status of the workspace flavor.  
  The valid values are as follows:
  + **normal**: The status is normal of the workspace flavor.
