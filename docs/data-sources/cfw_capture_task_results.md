---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_capture_task_results"
description: |-
  Use this data source to get the list of CFW capture task results.
---

# huaweicloud_cfw_capture_task_results

Use this data source to get the list of CFW capture task results.

## Example Usage

```hcl
variable "fw_instance_id" {}
variable "task_id" {}

data "huaweicloud_cfw_capture_task_results" "test" {
  fw_instance_id = var.fw_instance_id
  task_id        = var.task_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fw_instance_id` - (Required, String) Specifies the ID of the firewall instance.

* `task_id` - (Required, String) Specifies the capture task ID.

* `type` - (Optional, Int) Specifies whether to set a public IP address for downloading.
  The valid values are as follows:
  + **0**: unlimited;
  + **1**: set a public IP address for downloading. Currently, this feature can be used in the
  following regions: **cn-north-11**, **cn-east-5**, **af-south-1**, **ap-southeast-1**,
  **ap-southeast-3**, **ap-southeast-2**, **ap-southeast-4**, **tr-west-1**, **la-north-2**,
  **sa-brazil-1**, **la-south-2**, **me-east-1**.

* `ip` - (Optional, List) Specifies the public IP address ranges.
  A maximum of three address ranges can be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `captcha` - The captcha.

* `expires` - The expiration time of the download link.

* `file_list` - The capture file list.

  The [file_list](#data_file_list_struct) structure is documented below.

* `request_header` - The request header.

  The [request_header](#data_request_header_struct) structure is documented below.

* `url` - The download link.

<a name="data_file_list_struct"></a>
The `file_list` block supports:

* `file_name` - The file name.

* `url` - The download link.

* `file_path` - The file path.

<a name="data_request_header_struct"></a>
The `request_header` block supports:

* `host` - The host header information.
