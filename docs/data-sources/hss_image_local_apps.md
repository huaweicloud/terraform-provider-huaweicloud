---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_image_local_apps"
description: |-
  Use this data source to get the list of HSS image local apps within HuaweiCloud.
---

# huaweicloud_hss_image_local_apps

Use this data source to get the list of HSS image local apps within HuaweiCloud.

## Example Usage

```hcl
variable "image_id" {}

data "huaweicloud_hss_image_local_apps" "test" {
  image_id = var.image_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `image_id` - (Required, String) Specifies the image ID.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `app_name` - (Optional, String) Specifies app name filtering query, supports fuzzy matching.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_num` - The total number of image local apps.

* `data_list` - The list of image local apps.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `app_name` - The app name.

* `app_type` - The app type.

* `app_version` - The app version.

* `vul_num` - The number of vulnerabilities.
