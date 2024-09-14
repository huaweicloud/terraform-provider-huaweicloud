---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_logstash_configuration"
description: ""
---

# huaweicloud_css_logstash_configuration

Manages a CSS logstah configuration resource within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}
variable "logstash_conf_name" {}
variable "conf_content" {}

resource "huaweicloud_css_logstash_configuration" "test"  {
  cluster_id   = var.cluster_id
  name         = var.logstash_conf_name
  conf_content = var.conf_content

  setting {
    queue_type = "memory"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `cluster_id` - (Required, String, ForceNew) Specifies ID of the CSS logstash cluster.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the configuration file name of the CSS logstash cluster.
  Changing this creates a new resource.

* `conf_content` - (Required, String) Specifies the configuration file content of the CSS logstash cluster.

* `setting` - (Required, List) Specifies configuration file setting information of the CSS logstash cluster.
  The [setting](#Css_logstash_configuration_setting) structure is documented below.

* `sensitive_words` - (Optional, List) Specifies the input list of sensitive strings that need to be hidden.
  After configuring the hidden string list, all strings in the list will be hidden as `***` in the returned
  configuration content (the list supports a maximum of `20` items, and the maximum length of a single string
  is 512 bytes).

  -> **NOTE:** When this field is used, the configuration file content will also trigger update changes when the
    update operation is performed again, restoring the hidden content. If you import resources, you need to manually
    restore sensitive characters hidden in the configuration content.

<a name="Css_logstash_configuration_setting"></a>
The `setting` block supports:

* `workers` - (Optional, Int) Specifies the number of worker threads in the **Filters** + **Outputs** stage of
  the execution pipeline. The default value is the number of CPU cores.

* `batch_size` - (Optional, Int) Specifies the maximum number of events a single worker thread will collect
  from inputs before attempting to execute its **Filters** and **Outputs**. Larger values ​​are generally more
  efficient but increase memory overhead. Default is `125`.

* `batch_delay_ms` - (Optional, Int) Specifies the minimum time in the unit of milliseconds for each event to be
  waited for by pipeline scheduling.

* `queue_type` - (Required, String) Specifies internal queue model for event buffering.
  + **memory:** a traditional memory-based queue.
  + **persisted:** a disk-based ACKed persistence queue.

* `queue_check_point_writes` - (Optional, Int) Specifies the maximum number of events to be written before forcing
  a checkpoint when using a persistent queue, default is `1,024`.

* `queue_max_bytes_mb` - (Optional, Int) Specifies the total capacity of the persistent queue in megabytes (MB) when
  using a persistent queue. Make sure the disk is larger than this value. The default value is `1,024`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The configuration file content check status.

* `updated_at` - The update time of configuration file.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.

## Import

The CSS logstash configuration can be imported using `cluster_id` and `name` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_css_logstash_configuration.test <cluster_id>/<name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `sensitive_words`.
It is generally recommended running `terraform plan` after importing a CSS logstash configuration.
You can then decide if changes should be applied to the CSS logstash configuration, or the resource definition should
be updated to align with the CSS logstash configuration. Also you can ignore changes as below.

```hcl
resource "huaweicloud_css_logstash_configuration" "test" {
    ...

  lifecycle {
    ignore_changes = [
      sensitive_words,
    ]
  }
}
```
