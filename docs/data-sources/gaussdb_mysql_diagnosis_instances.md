---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_diagnosis_instances"
description: |-
  Use this data source to get the abnormal instance information by a specific metric.
---

# huaweicloud_gaussdb_mysql_diagnosis_instances

Use this data source to get the abnormal instance information by a specific metric.

## Example Usage

```hcl
variable "metric_name" {}

data "huaweicloud_gaussdb_mysql_diagnosis_instances" "test" {
  metric_name = var.metric_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `metric_name` - (Required, String) Specifies the metric name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instance_infos` - Indicates the information about the abnormal instances.

  The [instance_infos](#instance_infos_struct) structure is documented below.

<a name="instance_infos_struct"></a>
The `instance_infos` block supports:

* `instance_id` - Indicates the instance ID.

* `master_node_id` - Indicates the primary node ID.
