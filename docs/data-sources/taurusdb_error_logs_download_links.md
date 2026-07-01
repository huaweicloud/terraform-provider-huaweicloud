---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_error_logs_download_links"
description: |-
  Use this data source to query the download links of error logs of TaurusDB within HuaweiCloud.
---

# huaweicloud_taurusdb_error_logs_download_links

Use this data source to query the download links of error logs of TaurusDB within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}

data "huaweicloud_taurusdb_error_logs_download_links" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the TaurusDB instance.

* `node_id` - (Required, String) Specifies the ID of the instance node.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `list` - Indicates the list of the error log download links.

  The [list](#list_struct) structure is documented below.

<a name="list_struct"></a>
The `list` block supports:

* `job_id` - Indicates the task ID.

* `file_name` - Indicates the file name.

* `status` - Indicates the status.

* `file_size` - Indicates the file size.

* `file_link` - Indicates the link for downloading the file.

* `create_at` - Indicates the creation time.

* `updated_at` - Indicates the update time.

## Timeouts

This data source provides the following timeouts configuration options:

* `Read` - Default is 10 minutes.
