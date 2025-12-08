---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_image_files_statistics"
description: |-
  Use this data source to get the statistics of HSS image files within HuaweiCloud.
---

# huaweicloud_hss_image_files_statistics

Use this data source to get the statistics of HSS image files within HuaweiCloud.

## Example Usage

```hcl
variable "image_id" {}
variable "image_type" {}

data "huaweicloud_hss_image_files_statistics" "test" {
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
  + **private_image**: Private image repository
  + **shared_image**: Shared image repository
  + **local_image**: Local image
  + **instance_image**: Enterprise image
  + **registry**: Registry image
  + **local**: Local image, used to query global data

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `namespace` - (Optional, String) Specifies the organization name.

* `image_name` - (Optional, String) Specifies the image name.

* `tag_name` - (Optional, String) Specifies the image version name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_files_num` - The total number of image files.

* `total_files_size` - The total size of image files.
