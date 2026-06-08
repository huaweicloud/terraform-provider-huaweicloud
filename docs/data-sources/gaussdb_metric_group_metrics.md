---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_metric_group_metrics"
description: |-
  Use this data source to query the metrics of a metric group within HuaweiCloud.
---

# huaweicloud_gaussdb_metric_group_metrics

Use this data source to query the metrics of a metric group within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_gaussdb_metric_group_metrics" "test" {
  group_name = "CPUMEMORY"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the metric group names.
  If omitted, the provider-level region will be used.

* `group_name` - (Required, String) Specifies the metric group name.
  The valid values are as follows:
  + **CPUMEMORY**: CPU/memory.
  + **IOSTORAGE**: IO storage.
  + **NETWORK**: Network.
  + **CONNECTION**: Connection.
  + **TRANSACTION**: Transaction.
  + **LOCK**: Lock.
  + **SYNCSTAT**: Sync status.
  + **PROCESSRESOURCE**: Process resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `metric_names` - The list of metric names in the specified group.
  The [metric_names](#gaussdb_instance_metric_group_name_metric_names) structure is documented below.

<a name="gaussdb_instance_metric_group_name_metric_names"></a>
The `metric_names` block supports:

* `metric` - The metric ID.

* `name` - The metric name.
