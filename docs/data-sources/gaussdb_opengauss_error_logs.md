---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_error_logs"
description: |-
  Use this data source to get the link for downloading error logs.
---

# huaweicloud_gaussdb_opengauss_error_logs

Use this data source to get the link for downloading error logs.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_opengauss_error_logs" "test" {
  instance_id = var.instance_id
  start_time  = "2025-01-20T19:41:14+0800"
  end_time    = "2025-01-20T20:41:14+0800"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB OpenGauss instance.

* `start_time` - (Required, String) Specifies the start time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `end_time` - (Required, String) Specifies the end time in the **yyyy-mm-ddThh:mm:ssZ** format.
  Only error logs generated in the past 30 days can be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `log_files` - Indicates the log files.

  The [log_files](#log_files_struct) structure is documented below.

<a name="log_files_struct"></a>
The `log_files` block supports:

* `file_link` - Indicates the link for downloading the log file.

* `file_name` - Indicates the log file name.

* `file_size` - Indicates the log file size in KB.

* `status` - Indicates the log collection status.

* `start_time` - Indicates the log start time.

* `end_time` - Indicates the log end time.
