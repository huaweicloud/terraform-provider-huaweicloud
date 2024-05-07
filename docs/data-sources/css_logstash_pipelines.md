---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_logstash_pipelines"
description: |-
  Use this data source to get the list of CSS logstash pipelines.
---

# huaweicloud_css_logstash_pipelines

Use this data source to get the list of CSS logstash pipelines.

## Example Usage

```hcl
variable "cluster_id" {}
variable "config_file_name" {}

data "huaweicloud_css_logstash_pipelines" "test" {
  cluster_id = var.cluster_id
  name       = var.config_file_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies ID of the CSS logstash cluster.

* `name` - (Optional, String) Specifies the configuration file names of the CSS logstash cluster pipeline.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `pipelines` - The pipeline list of the CSS logstash cluster.

  The [pipelines](#pipelines_struct) structure is documented below.

<a name="pipelines_struct"></a>
The `pipelines` block supports:

* `update_at` - The update time of the CSS logstash cluster pipeline.

* `name` - The configuration file name of the CSS logstash cluster.

* `status` - The status of the CSS logstash cluster pipeline.

* `keep_alive` - Whether keep alive.

* `events` - The event of the CSS logstash cluster pipeline.
