---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_workload_queues"
description: |-
  Use this data source to get the list of workload queues.
---

# huaweicloud_dws_workload_queues

Use this data source to get the list of workload queues.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_dws_workload_queues" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the workload queues.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the cluster ID to which the workload queues belong.

* `name` - (Optional, String) Specifies the name of the workload queue.

* `logical_cluster_name` - (Optional, String) Specifies the name of the cluster. Required
  if the cluster is a logical cluster.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `queues` - The list of the workload queues.
  The [queues](#attrblock_queues) structure is documented below.

<a name="attrblock_queues"></a>
The `queues` block supports:

* `configuration` - The configuration information for workload queue.
  The [configuration](#attrblock_queues_configuration) structure is documented below.

* `logical_cluster_name` - The logical cluster name.

* `name` - The name of the workload queue.

* `short_query_concurrency_num` - The concurrency of short queries in the workload queue.

* `short_query_optimize` - Short query acceleration switch.
  + **true**: Support short query acceleration.
  + **false**: Short query acceleration not supported.

<a name="attrblock_queues_configuration"></a>
The `configuration` block supports:

* `resource_name` - The resource name.

* `resource_value` - The resource attribute value.

* `value_unit` - The resource attribute unit.
