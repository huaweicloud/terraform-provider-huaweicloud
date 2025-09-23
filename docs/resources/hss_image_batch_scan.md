---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_image_batch_scan"
description: |-
  Using this resource to batch scan images in image repositories within HuaweiCloud.
---

# huaweicloud_hss_image_batch_scan

Using this resource to batch scan images in image repositories within HuaweiCloud.

-> This resource is only a one-time action resource to batch scan images in image repositories. Deleting this resource
  will not clear the corresponding scan records, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "image_type" {}
variable "repo_type" {}

resource "huaweicloud_hss_image_batch_scan" "test" {
  image_type  = var.image_type
  repo_type   = var.repo_type
  operate_all = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `image_type` - (Required, String, NonUpdatable) Specifies the image type.  
  The valid values are as follows:
  + **private_image**: Private image repository.
  + **shared_image**: Shared image repository.

* `repo_type` - (Optional, String, NonUpdatable) Specifies the repository type. Currently only supports **SWR**.

* `image_info_list` - (Optional, List, NonUpdatable) Specifies the list of images to scan.  
  This parameter is required when `operate_all` is left blank or set to **false**.

  The [image_info_list](#image_info_list_struct) structure is documented below.

* `operate_all` - (Optional, Bool, NonUpdatable) Specifies whether to scan all images. The valid values are **true**
  and **false**.  
  If this parameter is left blank or set to **false**, the `image_info_list` is required.

* `namespace` - (Optional, String, NonUpdatable) Specifies the organization name.

* `image_name` - (Optional, String, NonUpdatable) Specifies the image name.

* `image_version` - (Optional, String, NonUpdatable) Specifies the image version.

* `scan_status` - (Optional, String, NonUpdatable) Specifies the scan status.  
  The valid values are as follows:
  + **unscan**: Not scanned.
  + **success**: Scan completed.
  + **scanning**: Scanning.
  + **failed**: Scan failed.
  + **download_failed**: Download failed.
  + **image_oversized**: Large image size.

* `latest_version` - (Optional, Bool, NonUpdatable) Specifies whether to only focus on the latest version image.

* `image_size` - (Optional, Int, NonUpdatable) Specifies the image size.

* `start_latest_update_time` - (Optional, Int, NonUpdatable) Specifies the start time of the creation date
  in milliseconds.

* `end_latest_update_time` - (Optional, Int, NonUpdatable) Specifies the end time of the creation date in milliseconds.

* `start_latest_scan_time` - (Optional, Int, NonUpdatable) Specifies the start time of the last scan completion date
  in milliseconds.

* `end_latest_scan_time` - (Optional, Int, NonUpdatable) Specifies the end time of the last scan completion date
  in milliseconds.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the hosts under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

<a name="image_info_list_struct"></a>
The `image_info_list` block supports:

* `namespace` - (Required, String, NonUpdatable) Specifies the namespace.

* `image_name` - (Required, String, NonUpdatable) Specifies the image name.

* `image_version` - (Required, String, NonUpdatable) Specifies the image version.

* `instance_id` - (Optional, String, NonUpdatable) Specifies the enterprise instance ID.

* `instance_url` - (Optional, String, NonUpdatable) Specifies the URL for downloading enterprise images.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
