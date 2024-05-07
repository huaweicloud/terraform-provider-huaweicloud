---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_logstash_pipeline"
description: ""
---

# huaweicloud_css_logstash_pipeline

Manages CSS logstash cluster pipeline resource within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}
variable "config_file_names" {}

resource "huaweicloud_css_logstash_pipeline" "test" {
  cluster_id = var.cluster_id
  names      = var.config_file_names
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `cluster_id` - (Required, String, ForceNew) Specifies ID of the CSS logstash cluster.
  Changing this creates a new resource.

* `names` - (Required, List) Specifies the configuration file names of the CSS logstash cluster pipeline.
  Changing this creates a new resource.

* `keep_alive` - (Optional, Bool, ForceNew) Specifies whether keep alive. The value can be **true** and **false**.
  Defaults to **false**. During hot start, the value of keep alive of existing pipelines in the cluster needs to
  be consistent.
  Changing this creates a new resource.

  -> **NOTE:** Keepalive can be enabled for long-running services. Enabling it will configure a daemon process
    on each node. If the Logstash service is faulty, the daemon process will rectify the fault and restart the
    service. Do not enable it for short running services, or your migration tasks may fail due to lack of source data.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `pipelines` - The pipeline list of the CSS logstash cluster.
  The [pipelines](#css_logstash_pipelines) structure is documented below.

<a name="css_logstash_pipelines"></a>
The `pipelines` block supports:

* `name` - The configuration file name of the CSS logstash cluster.

* `keep_alive` - Whether keep alive.

* `events` - The event of the CSS logstash cluster pipeline.
  The [events](#css_logstash_pipelines_events) structure is documented below.

  -> **Note:** Events can only be viewed in real time in the "working" state (manual refresh is required).
    In the "Stopped" state, please go to the output side to view the amount of migrated data.

* `status` - The status of the CSS logstash cluster pipeline.

* `updated_at` - The update time of the CSS logstash cluster pipeline.

<a name="css_logstash_pipelines_events"></a>
The `events` block supports:

* `in` - The number of received data that needs to be processed.

* `filtered` - The number of data to be filtered.

* `out` - The number of output data.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.

## Import

The CSS logstash cluster pipeline can be imported using `cluster_id`, e.g.

```bash
$ terraform import huaweicloud_css_logstash_pipeline.test <cluster_id>
```
