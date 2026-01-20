---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_image_baseline_statistic"
description: |-
  Use this data source to get the HSS image baseline statistic within HuaweiCloud.
---

# huaweicloud_hss_image_baseline_statistic

Use this data source to get the HSS image baseline statistic within HuaweiCloud.

## Example Usage

```hcl
variable "image_type" {}

data "huaweicloud_hss_image_baseline_statistic" "test" {
  image_type = var.image_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `image_type` - (Required, String) Specifies the type of the image.  
  The valid values are as follows:
  + **private_image**: Private image repository.
  + **shared_image**: Shared image repository.
  + **local_image**: Local image.
  + **instance_image**: Enterprise image.
  + **registry**: Registry image.
  + **local**: Local image, used to query global data.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `namespace` - (Optional, String) Specifies the organization name. When there is no image related information, it means
  to query all images.

* `image_name` - (Optional, String) Specifies the name of the image.

* `image_version` - (Optional, String) Specifies the version of the image.

* `instance_id` - (Optional, String) Specifies the enterprise warehouse instance ID, SWR shared version does not require
  this parameter.

* `image_id` - (Optional, String) Specifies the ID of the image.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `image_weak_pwd` - The weak password detection.

* `pwd_policy` - The password complexity strategy detection.

* `security_check` - The server configuration detection.
