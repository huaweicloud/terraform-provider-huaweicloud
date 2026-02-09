---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_webtamper_image_options"
description: |-
  Use this data source to get the list of HSS web tamper image options within HuaweiCloud.
---

# huaweicloud_hss_webtamper_image_options

Use this data source to get the list of HSS web tamper image options within HuaweiCloud.

## Example Usage

```hcl
variable "image_type" {}

data "huaweicloud_hss_webtamper_image_options" "test" {
  image_type = var.image_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `image_type` - (Required, String) Specifies the type of the image.  
  The valid values are as follows:
  + **registry**
  + **local**

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `registry_type` - (Optional, String) Specifies the image warehouse type for the warehouse image.

* `image_namespace` - (Optional, String) Specifies the organization name of the warehouse image.

* `registry_name` - (Optional, String) Specifies the image warehouse name for the specified warehouse image.

* `image_name` - (Optional, String) Specifies the image name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_num` - The total number.

* `data_list` - The data list.

The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `image_name` - The image name.

* `image_full_name` - The image full name.

* `image_id` - The image ID.

* `image_version_list` - The image version list.

* `image_namespace` - The organization name of the warehouse image.

* `registry_name` - The image warehouse name for the specified warehouse image.

* `registry_type` - The image warehouse type for the warehouse image.  
  The valid values are as follows:
  + **SwrPrivate**: SWR private repository.
  + **SwrShared**: SWR shared repository.
  + **SwrEnterprise**: SWR enterprise repository.
  + **Harbor**: Harbor repository.
  + **Jfrog**: JFrog repository.
  + **Other**: Other repository.
