---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_image_sensitive_information_detail"
description: |-
  Use this data source to get the list of HSS image sensitive information detail within HuaweiCloud.
---

# huaweicloud_hss_image_sensitive_information_detail

Use this data source to get the list of HSS image sensitive information detail within HuaweiCloud.

## Example Usage

```hcl
variable "image_id" {}
variable "image_type" {}

data "huaweicloud_hss_image_sensitive_information_detail" "test" {
  image_id   = var.image_id
  image_type = var.image_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `image_id` - (Required, String) Specifies the image ID.

* `image_type` - (Required, String) Specifies the image type. Valid values are:
  + **private_image**: SWR private image.
  + **shared_image**: SWR shared image.
  + **instance_image**: SWR enterprise edition image.
  + **cicd**: CICD image.
  + **harbor**: Harbor repository image.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `namespace` - (Optional, String) Specifies the organization name.

* `image_name` - (Optional, String) Specifies the image name.

* `image_version` - (Optional, String) Specifies the image version name.

* `file_path` - (Optional, String) Specifies the file path.

* `severity` - (Optional, String) Specifies the threat level. Valid values are:
  + **critical**: Critical.
  + **high**: High risk.
  + **medium**: Medium risk.
  + **low**: Low risk.

* `handle_status` - (Optional, String) Specifies whether it has been handled. Valid values are:
  + **unhandled**: Not handled.
  + **handled**: Handled.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_num` - The total number of image sensitive information records.

* `data_list` - The list of image sensitive information data.

  The [data_list](#image_sensitive_information_detail_data_list) structure is documented below.

<a name="image_sensitive_information_detail_data_list"></a>
The `data_list` block supports:

* `sensitive_info_id` - The sensitive event ID.

* `severity` - The threat level. Valid values are:
  + **critical**: Critical.
  + **high**: High risk.
  + **medium**: Medium risk.
  + **low**: Low risk.

* `name` - The rule name.

* `description` - The rule description.

* `position` - The layer where the sensitive information is located in the image.

* `file_path` - The file path.

* `content` - The sensitive information content.

* `latest_scan_time` - The last scan time, in milliseconds.

* `handle_status` - Whether it has been handled. Valid values are:
  + **unhandled**: Not handled.
  + **handled**: Handled.

* `operate_accept` - The operation type. Valid values are:
  + **ignore**: Ignore.
  + **do_not_ignore**: Do not ignore.
