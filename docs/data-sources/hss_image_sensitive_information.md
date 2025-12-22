---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_image_sensitive_information"
description: |-
  Use this data source to get the list of HSS image sensitive information within HuaweiCloud.
---

# huaweicloud_hss_image_sensitive_information

Use this data source to get the list of HSS image sensitive information within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_image_sensitive_information" "test" {}
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

* `image_type` - (Optional, String) Specifies the image type. Valid values are:
  + **registry**: Repository image.
  + **cicd**: CICD image.

  Defaults to **registry**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_num` - The total number of image sensitive information records.

* `data_list` - The list of image sensitive information data.

  The [data_list](#image_sensitive_information_data_list) structure is documented below.

<a name="image_sensitive_information_data_list"></a>
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
