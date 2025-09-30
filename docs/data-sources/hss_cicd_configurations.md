---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_cicd_configurations"
description: |-
  Use this data source to get the list of HSS CiCd configurations within HuaweiCloud.
---

# huaweicloud_hss_cicd_configurations

Use this data source to get the list of HSS CiCd configurations within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_cicd_configurations" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cicd_id` - (Optional, String) Specifies the CiCd ID.

* `cicd_name` - (Optional, String) Specifies the CiCd name.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data_list` - The CiCd configuration list.

  The [data_list](#app_structure) structure is documented below.

<a name="app_structure"></a>
The `data_list` block supports:

* `cicd_id` - The CiCd ID.

* `cicd_name` - The CiCd name.

* `associated_images_num` - The number of associated mirror scans.

* `associated_config_num` - The number of associated configuration scans.
