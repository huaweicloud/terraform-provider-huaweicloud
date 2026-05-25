---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_wdr_snapshot_collect"
description: |-
  Manages a GaussDB WDR snapshot collect resource within HuaweiCloud.
---

# huaweicloud_gaussdb_wdr_snapshot_collect

Manages a GaussDB WDR snapshot collect resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "start_time" {}
variable "end_time" {}

resource "huaweicloud_gaussdb_wdr_snapshot_collect" "test" {
  instance_id = var.instance_id
  start_time  = var.start_time
  end_time    = var.end_time
  wdr_type    = ["cluster"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the GaussDB instance ID.

* `start_time` - (Required, String, NonUpdatable) Specifies the snapshot start time. The format is **yyyy-mm-ddThh:mm:ssZ**,
  where T indicates the start of a certain time, and Z indicates the time zone offset.

* `end_time` - (Required, String, NonUpdatable) Specifies the snapshot end time. The format is **yyyy-mm-ddThh:mm:ssZ**,
  where T indicates the start of a certain time, and Z indicates the time zone offset.

* `wdr_type` - (Required, List, NonUpdatable) Specifies the collection type, which can be **instance-level** or
  **component-level**. For the instance level, enter **cluster**. For the component level, enter the component ID.
  The two collection types cannot be specified at the same time.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the job ID.

* `file_size` - The file size in KB.

* `wdr_type_attr` - The collection type. The value can be:
  + **cluster**: instance-level
  + **component**: component-level
  + **pdb**: tenant-level

* `job_create_time` - The job creation time.

* `start_snapshot_id` - The first comparison snapshot ID.

* `end_snapshot_id` - The second comparison snapshot ID.

* `download_url` - The report download link. The link is valid for 30 minutes.

* `status` - The collection status.
  The valid values are as follows:
  + **SUCCESS**: Collection successful.
  + **FAILED**: Collection failed.
  + **EXPORTING**: Collecting.

* `notes` - The remarks. When the collection type is component level, the content includes the component ID.

* `error_msg` - The error message.

* `file_name` - The WDR report temporary file name.

* `file_path` - The WDR report temporary file save path.

* `obs_bucket` - The OBS bucket information for storing the WDR report temporary file.
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

* `create` - Default is 30 minute.

## Import

The GaussDB WDR snapshot collect can be imported using the `instance_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_gaussdb_wdr_snapshot_collect.test <instance_id>/<id>
```
