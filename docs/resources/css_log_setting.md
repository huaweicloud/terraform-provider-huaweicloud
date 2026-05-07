---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_log_setting"
description: "Use this resource to manage log backup or log ingestion of a CSS cluster within HuaweiCloud."
---

# huaweicloud_css_log_setting

Use this resource to manage log backup or log ingestion of a CSS cluster within HuaweiCloud.

## Example Usage

### Manage Log Backup

```hcl
variable "cluster_id" {}
variable "agency" {}
variable "base_path" {}
variable "bucket" {}
variable "period" {}

resource "huaweicloud_css_log_setting" "test" {
  cluster_id = var.cluster_id
  agency     = var.agency
  base_path  = var.base_path
  bucket     = var.bucket
  period     = var.period
}
```

### Manage Log Ingestion

```hcl
variable "cluster_id" {}
variable "index_prefix" {}
variable "keep_days" {}
variable "target_cluster_id" {}

resource "huaweicloud_css_log_setting" "test" {
  cluster_id        = var.cluster_id
  action            = "real_time_log_collect"
  index_prefix      = var.index_prefix
  keep_days         = var.keep_days
  target_cluster_id = var.target_cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `cluster_id` - (Required, String, ForceNew) Specifies ID of the cluster whose log function you want to enable.
  Changing this creates a new resource.

* `action` - (Optional, String, ForceNew) Specifies the log setting action type. Changing this creates a new resource.
  Valid values are as follows:
  + **base_log_collect**: Enables log backup.
  + **real_time_log_collect**: Enable log ingestion.
  
  Defaults to **base_log_collect**.

  -> **Note**: Log ingestion is only supported in the scenarios as follows:
  The Elasticsearch cluster version is `7.10.2` and the image version is no earlier than `7.10.2_24.2.0_×.×.×`.
  The OpenSearch cluster version is `1.3.6` or `2.19.0` and the image version is no earlier than `×.×.×_24.2.0_×.×.×`.

* `agency` - (Optional, String) Specifies the agency name. You can create an agency to allow CSS to
  call other cloud services. This parameter is mandatory when `action` is set to **base_log_collect**.

* `base_path` - (Optional, String) Specifies the storage path of backed up logs in the OBS bucket.
  This parameter is mandatory when `action` is set to **base_log_collect**.

* `bucket` - (Optional, String) Specifies the name of the OBS bucket for storing logs.
  This parameter is mandatory when `action` is set to **base_log_collect**.

* `period` - (Optional, String) Specifies the start time of automatic backup.
  The value must be in the format of **HH:MM GMT±HH:MM**, e.g., **01:00 GMT+08:00**.
  This parameter is optional only when `action` is set to **base_log_collect**.

* `index_prefix` - (Optional, String) Specifies the index prefix for real-time log ingestion.
  This parameter is mandatory when `action` is set to **real_time_log_collect**.

* `keep_days` - (Optional, Int) Specifies the number of days to keep real-time logs. The value ranges from 1 to 3650.
  This parameter is mandatory when `action` is set to **real_time_log_collect**.

* `target_cluster_id` - (Optional, String) Specifies the ID of the target cluster for real-time log ingestion.
  This parameter is mandatory when `action` is set to **real_time_log_collect**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `log_switch` - Whether the log backup is enabled.
  + **true**: Log backup is enabled.
  + **false**: Log backup is disabled.

* `updated_at` - The update time of the log backup, in RFC3339 format.

* `auto_enabled` - Whether log automatic backup is enabled.
   + **true**: Automatic backup is enabled.
   + **false**: Automatic backup is disabled.

* `status` - The status of log ingestion task. The values are as follows:
   + **100**: A real-time log ingestion task is being created.
   + **150**: Real-time log ingestion is available.
   + **200**: Real-time log ingestion task is activated.
   + **300**: Real-time log ingestion failed.
   + **302**: The real-time log ingestion task failed to be deleted.
   + **303**: The real-time log ingestion task failed to be created.
   + **304**: The real-time log ingestion task is being disabled.
   + **400**: The real-time log ingestion task is disabled.

* `log_ingestion_create_at` - The creation time of real-time log ingestion task, in RFC3339 format.

* `log_ingestion_update_at` - The update time of real-time log ingestion task, in RFC3339 format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

### Import Log Backup

The CSS log backup setting can be imported using `cluster_id` (defaults to **base_log_collect**)
or explicitly specifying the action, e.g.

```bash
$ terraform import huaweicloud_css_log_setting.test <cluster_id>
```

Or:

```bash
$ terraform import huaweicloud_css_log_setting.test <cluster_id>/base_log_collect
```

### Import Log Ingestion

The CSS log ingestion setting can be imported using `cluster_id` and **real_time_log_collect**, e.g.

```bash
$ terraform import huaweicloud_css_log_setting.test <cluster_id>/real_time_log_collect
```
