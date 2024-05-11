---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_logstash_configurations"
description: |-
  Use this data source to get the list of CSS logstash configurations.
---

# huaweicloud_css_logstash_configurations

Use this data source to get the list of CSS logstash configurations.

## Example Usage

```hcl
variable "cluster_id" {}
variable "name" {}

data "huaweicloud_css_logstash_configurations" "test" {
  cluster_id = var.cluster_id
  name       = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies ID of the CSS logstash cluster.

* `name` - (Optional, String) Specifies the configuration file name.

* `status` - (Optional, String) Specifies the configuration file content check status.
  The values can be **checking**, **available** and **unavailable**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `confs` - The configuration file list.

  The [confs](#confs_struct) structure is documented below.

<a name="confs_struct"></a>
The `confs` block supports:

* `name` - The configuration file name.

* `status` - The configuration file content check status.

* `conf_content` - The configuration file content.

* `setting` - The configuration file setting information.

  The [setting](#confs_setting_struct) structure is documented below.

* `updated_at` - The update time.

<a name="confs_setting_struct"></a>
The `setting` block supports:

* `workers` - The number of worker threads.
  that is in the **Filters** + **Outputs** stage of the execution pipeline.
  The default value is the number of CPU cores.

* `batch_size` - The maximum number of events.
  This event refers to a single worker thread will collect from inputs before attempting
  to execute its **Filters** and **Outputs**.
  Larger values are generally more efficient but increase memory overhead.

* `batch_delay_ms` - The minimum time for an event to be scheduled to wait.

* `queue_type` - Internal queue model for event buffering.

* `queue_check_point_writes` - The maximum number of events to be written.
  This refers to before forcing a checkpoint when a persistent queue is used.

* `queue_max_bytes_mb` - The total capacity of the persistent queue.
  The unit is megabytes. The disk is guaranteed to have a maximum capacity of this value.
