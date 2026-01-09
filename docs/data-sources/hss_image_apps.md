---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_image_apps"
description: |-
  Use this data source to get the list of HSS image apps within HuaweiCloud.
---

# huaweicloud_hss_image_apps

Use this data source to get the list of HSS image apps within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_image_apps" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `image_type` - (Optional, String) Specifies the image type.  
  The valid values are as follows:
  + **local**: Local image.
  + **registry**: Warehouse image.
  + **cicd**: CICD image.

* `app_name` - (Optional, String) Specifies the app name.

* `is_compliant` - (Optional, String) Specifies whether it is compliant. The value can be **true** or **false**.
  Only supports **false*.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_num` - The total number of image apps.

* `data_list` - The list of image apps.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `apps` block supports:

* `app_name` - The app name.

* `app_type` - The app type.

* `app_version` - The app version.

* `vul_num` - The number of vulnerabilities.

* `app_path` - The app path.

* `layer_digest` - The layer digest.

* `is_compliant` - Whether it is compliant.

* `latest_scan_time` - The last detection time, in milliseconds (ms).

* `image_type` - The warehouse image type.  
  The valid values are as follows:
  + **SwrPrivate**: SWR private image.
  + **SwrShared**: SWR sharing.
  + **SwrEnterprise**: SWR Enterprise.
  + **Harbor**: Harbor Warehouse.
  + **Jfrog**: Jfrog warehouse.
  + **Other**: Other warehouses.

* `namespace` - The namespace.

* `image_name` - The image name.

* `image_version` - The image version.

* `image_id` - The image ID.
