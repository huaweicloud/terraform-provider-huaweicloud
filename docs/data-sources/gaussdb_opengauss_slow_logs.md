---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_slow_logs"
description: |-
  Use this data source to get the list of downloaded slow query log information.
---

# huaweicloud_gaussdb_opengauss_slow_logs

Use this data source to get the list of downloaded slow query log information.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_opengauss_slow_logs" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB OpenGauss instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `list` - Indicates the list of downloaded slow query log information.

  The [list](#list_struct) structure is documented below.

<a name="list_struct"></a>
The `list` block supports:

* `id` - Indicates the slow query log ID.

* `instance_id` - Indicates the instance ID.

* `node_id` - Indicates the node ID.

* `workflow_id` - Indicates the workflow ID.

* `file_name` - Indicates the file name.

* `file_size` - Indicates the file size in bytes.

* `file_link` - Indicates the link for downloading the file.

* `bucket_name` - Indicates the bucket name.

* `created_at` - Indicates the creation time.

* `updated_at` - Indicates the update time.

* `version` - Indicates the version.

* `status` - Indicates the status.

* `message` - Indicates the message.

## Timeouts

This resource provides the following timeouts configuration options:

* `read` - Default is 10 minutes.
