---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_snapshot_setting"
description: "Use this resource to manage snapshot settings of a CSS cluster."
---

# huaweicloud_css_snapshot_setting

Use this resource to manage snapshot settings of a CSS cluster within HuaweiCloud.

## Example Usage

### Basic Snapshot Setting

```hcl
variable "cluster_id" {}
variable "agency" {}
variable "bucket" {}
variable "base_path" {}

resource "huaweicloud_css_snapshot_setting" "test" {
  cluster_id = var.cluster_id
  agency     = var.agency
  bucket     = var.bucket
  base_path  = var.base_path
}
```

### Enable Automatic Snapshot

```hcl
variable "cluster_id" {}
variable "agency" {}
variable "bucket" {}
variable "base_path" {}

resource "huaweicloud_css_snapshot_setting" "test" {
  cluster_id                     = var.cluster_id
  agency                         = var.agency
  bucket                         = var.bucket
  base_path                      = var.base_path
  max_snapshot_bytes_per_seconds = "100mb"
  max_restore_bytes_per_seconds  = "100mb"
  enable                         = "true"
  indices                        = "index1,index2"
  prefix                         = "snapshot"
  period                         = "02:00 GMT+08:00"
  keepday                        = 7
  frequency                      = "DAY"
  delete_auto                    = "false"
}
```

### Disable Automatic Snapshot

```hcl
variable "cluster_id" {}
variable "agency" {}
variable "bucket" {}
variable "base_path" {}

resource "huaweicloud_css_snapshot_setting" "test" {
  cluster_id                     = var.cluster_id
  agency                         = var.agency
  bucket                         = var.bucket
  base_path                      = var.base_path
  max_snapshot_bytes_per_seconds = "100mb"
  max_restore_bytes_per_seconds  = "100mb"
  enable                         = "false"
  delete_auto                    = "true"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the ID of the CSS cluster.

* `agency` - (Required, String) Specifies the agency name. You can create an agency to allow CSS to
  access OBS buckets for storing snapshots.

* `bucket` - (Required, String) Specifies the name of the OBS bucket for storing snapshots.

* `base_path` - (Optional, String) Specifies the storage path of snapshots in the OBS bucket.

* `max_snapshot_bytes_per_seconds` - (Optional, String) Specifies the maximum backup speed per node (bytes per second).
  The following values and formats are allowed:
  + **Number + Unit**: The number ranges from `0` to `9999`. The unit can be `k`, `kb`, `m`, `mb`, `g`, `gb`,
    `t`, `tb`, `p`, `pb`, or `b` (case-insensitive). For example: `40mb`, `100MB`, `1gb`.
  + **0mb**, **0**, **-1**: Indicates no speed limit.
  
  Defaults to **40mb**. If this parameter is left blank, the default value is used.
  
  -> **Note:** Setting this parameter to `0mb`, `0`, or `-1` means there is no speed limit.
  An overly high backup speed may lead to excessive resource usage, which may impact system stability.
  When it is exceeded, flow control is triggered to prevent excessive resource usage and ensure
  system stability. The actual backup speed may not reach the configured value, as it depends on
  many factors, such as OBS performance and disk I/O.

* `max_restore_bytes_per_seconds` - (Optional, String) Specifies the maximum restoration speed per node (bytes per second).
  The following values and formats are allowed:
  + **Number + Unit**: The number ranges from `0` to `9999`. The unit can be `k`, `kb`, `m`, `mb`, `g`, `gb`,
    `t`, `tb`, `p`, `pb`, or `b` (case-insensitive). For example: `40mb`, `100MB`, `1gb`.
  + **0mb**, **0**, **-1**: Indicates no speed limit.
  
  The default value depends on the cluster type and version:
  + For Elasticsearch clusters of version `7.6.2` or earlier, the default value is `40 MB`.
  + For OpenSearch clusters and Elasticsearch clusters later than version `7.6.2`, the default setting is no limit,
    but the recovery speed is still limited by the `indices.recovery.max_bytes_per_sec` parameter.
  
  -> **Note:** Setting this parameter to `0mb`, `0`, or `-1` means there is no speed limit.
  However, for OpenSearch clusters and Elasticsearch clusters later than version `7.6.2`,
  the recovery speed is also limited by the `indices.recovery.max_bytes_per_sec` parameter.
  An overly high recovery speed may lead to excessive resource usage, which may impact system stability.
  When it is exceeded, flow control is triggered to prevent excessive resource usage and ensure system stability.
  The actual recovery speed may not reach the configured value, as it depends on many factors,
  such as OBS performance and disk I/O.

* `enable` - (Optional, String) Specifies whether to enable automatic snapshot creation.
  The valid values are **true** and **false**. Defaults to **false**.

 -> **Note:** Only Elasticsearch and OpenSearch clusters support automatic snapshot creation.
  When set to **false**, the automatic snapshot policy parameters(`indices`, `prefix`, `period`,
  `keepday`, `frequency`) will not take effect.

* `indices` - (Optional, String) Specifies the indexes to be backed up. Multiple indexes are
  separated by commas (,), for example, **index1,index2**. By default, all indexes are backed up.
  This parameter takes effect only when `enable` is set to **true**.

* `prefix` - (Optional, String) Specifies the prefix of the snapshot name.
  This parameter is mandatory when `enable` is set to **true**. And takes effect only when `enable`
  is set to **true**.

* `period` - (Optional, String) Specifies the time at which a snapshot is created.
  The valid format is **HH:mm GMT+XX:XX**. For example, **01:00 GMT+08:00** indicates that snapshots
  are created at 01:00 UTC+08:00 time every day.
  This parameter takes effect only when `enable` is set to **true**.

* `keepday` - (Optional, Int) Specifies the number of days for retaining generated snapshots.
  This parameter is mandatory when `enable` is set to **true**. and takes effect only when `enable` is set to **true**.
  The value ranges from **1** to **90**.

* `frequency` - (Optional, String) Specifies the interval at which snapshots are created.
  The valid values are:
  + **DAY**: Snapshots are created once a day.
  + **WEEK**: Snapshots are created once a week.
  + **SUN**, **MON**, **TUE**, **WED**, **THU**, **FRI**, **SAT**: Execute the task at the specified
  day of every week. For example, **SUN** indicates that the task is executed once every Sunday.
  This parameter takes effect only when the parameter `enable` is set to **true**.

* `delete_auto` - (Optional, String) Specifies whether to automatically delete expired snapshots.
  The valid values are **true** and **false**. Defaults to **false**.
  This parameter takes effect only when `enable` is set to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format of `<cluster_id>`.

* `snapshot_cmk_id` - The CMK ID used for snapshot encryption.

## Import

The CSS snapshot setting can be imported using the `cluster_id`, e.g.

```bash
$ terraform import huaweicloud_css_snapshot_setting.test <cluster_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `delete_auto`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_css_snapshot_setting" "test" {
  ...
  
  lifecycle {
    ignore_changes = [
      delete_auto
    ]
  }
}
```
