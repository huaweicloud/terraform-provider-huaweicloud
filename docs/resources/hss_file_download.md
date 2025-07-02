---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_file_download"
description: |-
  Manages an HSS file download resource within HuaweiCloud.
---

# huaweicloud_hss_file_download

Manages an HSS file download resource within HuaweiCloud.

-> This resource is only a one-time action resource for HSS file download. Deleting this resource will not clear the
  corresponding request record, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "file_id" {}

resource "huaweicloud_hss_file_download" "test" {
  file_id = var.file_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `file_id` - (Required, String, NonUpdatable) Specifies the file ID.  
  You can obtain the `file_id` using the `huaweicloud_hss_vulnerability_information_export` resource.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID to which the hosts
  belong.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the hosts under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `export_file_name` - (Optional, String, NonUpdatable) Specifies the file name that can save data.
  Defaults to **hss-export-{{file_id}}.zip**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
