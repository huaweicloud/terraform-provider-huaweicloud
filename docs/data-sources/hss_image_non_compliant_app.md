---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_image_non_compliant_app"
description: |-
  Use this data source to get the list of HSS non-compliant app of the image within HuaweiCloud.
---

# huaweicloud_hss_image_non_compliant_app

Use this data source to get the list of HSS non-compliant app of the image within HuaweiCloud.

## Example Usage

```hcl
variable "image_id" {}
variable "image_type" {}

data "huaweicloud_hss_image_non_compliant_app" "test" {
  image_id   = var.image_id
  image_type = var.image_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `image_id` - (Required, String) Specifies the image ID.

* `image_type` - (Required, String) Specifies the image type.  
  The valid values are:
  + **private_image**: SWR private image.
  + **shared_image**: SWR shared image.
  + **instance_image**: SWR enterprise edition image.
  + **cicd**: CICD image.
  + **harbor**: Harbor warehouse image.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `namespace` - (Optional, String) Specifies the namespace.

* `image_name` - (Optional, String) Specifies the image name.

* `image_version` - (Optional, String) Specifies the image version.

* `app_name` - (Optional, String) Specifies the non-compliant app name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number.

* `data_list` - The list of non-compliant app of the image.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `id` - The ID.

* `app_name` - The non-compliant app name.

* `app_path` - The non-compliant app path.

* `app_version` - The non-compliant app version.

* `layer_digest` - The layer digest.
