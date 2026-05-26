---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_asp_collect"
description: |-
  Manages a GaussDB ASP collect resource within HuaweiCloud.
---

# huaweicloud_gaussdb_asp_collect

Manages a GaussDB ASP collect resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_gaussdb_asp_collect" "test" {
  instance_id = var.instance_id
  start_time  = "2024-01-01T00:00:00+0800"
  end_time    = "2024-12-31T23:59:59+0800"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the ASP collect. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the GaussDB instance.

* `start_time` - (Required, String, NonUpdatable) Specifies the start time for ASP collect. The format is
  **yyyy-mm-ddThh:mm:ssZ**, where T indicates the start of a certain time, and Z indicates the time zone offset.

* `end_time` - (Required, String, NonUpdatable) Specifies the end time for ASP collect. The format is
  **yyyy-mm-ddThh:mm:ssZ**, where T indicates the start of a certain time, and Z indicates the time zone offset.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the job ID.

* `file_size` - The file size in KB.

* `download_url` - The report download link. The link is valid for 30 minutes.

* `status` - The collection status.  The valid values are as follows:
  + **SUCCESS**: Collection successful.
  + **FAILED**: Collection failed.
  + **EXPORTING**: Collecting.

* `file_name` - The ASP report file name.

* `file_path` - The ASP report file save path.

* `obs_bucket` - The OBS bucket information for storing the ASP report file.
  The [obs_bucket](#obs_bucket_struct) structure is documented below.

<a name="obs_bucket_struct"></a>
The `obs_bucket` block supports:

* `name` - The OBS bucket name.

* `type` - The OBS bucket type.
  The valid values are as follows:
  + **common**: Public bucket.
  + **aps**: Intelligent O&M dedicated bucket.

* `url` - The OBS service access address.

* `port` - The OBS service port number.

* `domain_id` - The tenant ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default timeout is 30 minutes.

## Import

The GaussDB ASP collect can be imported using the `instance_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_gaussdb_asp_collect.test <instance_id>/<id>
```
