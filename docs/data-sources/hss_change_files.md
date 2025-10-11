---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_change_files"
description: |-
  Use this data source to get the list of HSS change files within HuaweiCloud.
---

# huaweicloud_hss_change_files

Use this data source to get the list of HSS change files within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_change_files" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `host_name` - (Optional, String) Specifies the server name. Character length `1` - `256` bits.

* `begin_time` - (Optional, String) Specifies the start time as a 13-digit timestamp. Minimum value `0`, maximum value
  `9,223,372,036,854,775,807`.

* `end_time` - (Optional, String) Specifies the end time as a 13-digit timestamp. Minimum value `0`, maximum value
  `9,223,372,036,854,775,807`.

* `file_name` - (Optional, String) Specifies the file name.

* `file_path` - (Optional, String) Specifies the file path.

* `change_type` - (Optional, String) Specifies the change type. Valid values are:
  + **all**: All.
  + **registry**: Registry.
  + **file**: File.

* `change_category` - (Optional, String) Specifies the change category. Valid values are:
  + **all**: All.
  + **modify**: Modify.
  + **add**: Add.
  + **delete**: Delete.

* `status` - (Optional, String) Specifies the status. Valid values are:
  + **all**: All.
  + **trust**: Trusted.
  + **untrust**: Untrusted.
  + **unknown**: Unknown.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data_list` - The list of change file information.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `id` - The ID.

* `file_name` - The file name.

* `file_path` - The file path.

* `status` - The status.

* `host_name` - The server name.

* `change_type` - The change type.

* `change_category` - The change category.

* `after_change` - The modified hash.

* `before_change` - The hash.

* `latest_time` - The last change time.
