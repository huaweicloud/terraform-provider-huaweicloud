---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_event_delete_isolated_file"
description: |-
  Using this resource to delete HSS isolated file within HuaweiCloud."
---

# huaweicloud_hss_event_delete_isolated_file

Using this resource to delete HSS isolated file within HuaweiCloud.

-> This resource is only a one-time action resource to delete HSS isolated file. Deleting this resource will not
  clear the corresponding isolated file deletion record, but will only remove the resource information from the
  tf state file.

## Example Usage

```hcl
variable "host_id" {}
variable "file_hash" {}
variable "file_path" {}
variable "file_attr" {}

resource "huaweicloud_hss_event_delete_isolated_file" "test" {
  data_list {
    host_id   = var.host_id
    file_hash = var.file_hash
    file_path = var.file_path
    file_attr = var.file_attr
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted, the
  provider-level region will be used.

* `data_list` - (Required, List, NonUpdatable) Specifies the list of files to be deleted.
  The [data_list](#data_list_struct) structure is documented below.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the hosts under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `host_id` - (Required, String, NonUpdatable) Specifies the host ID.

* `file_hash` - (Required, String, NonUpdatable) Specifies the file hash.

* `file_path` - (Required, String, NonUpdatable) Specifies the file path.

* `file_attr` - (Required, String, NonUpdatable) Specifies the file attribute.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
