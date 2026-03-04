---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_resize_flavors"
description: |-
  Use this data source to query available node specifications for a cluster.
---

# huaweicloud_css_resize_flavors

Use this data source to query available node specifications for a cluster.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_css_resize_flavors" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the cluster ID.

* `type` - (Optional, String) Specifies the cluster node type.
  The valid values are as follows:
  + **ess**: Data node.
  + **ess-cold**: Cold data node.
  + **ess-client**: Client node.
  + **ess-master**: Master node.
  + **lgs**: Logstash node.

  If the cluster is Elasticsearch or OpenSearch cluster, the default value is **ess**,
  If the cluster is Logstash cluster, the default value is **lgs**,

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `datastore_id` - The engine type ID.

* `dbname` - The engine name.

* `versions` - The list of versions.
  The [versions](#versions_struct) structure is documented below.

<a name="versions_struct"></a>
The `versions` block supports:

* `id` - The engine version ID.

* `name` - The engine version name.

* `flavors` - The list of node objects.
  The [flavors](#flavors_struct) structure is documented below.

<a name="flavors_struct"></a>
The `flavors` block supports:

* `str_id` - The flavor ID.

* `cpu` - The number of vCPUs available with an instance.

* `ram` - The memory size of an instance, in GB.

* `name` - The flavor name.

* `region` - The regions where the node flavor is available.

* `diskrange` - The disk capacity of an instance, in GB.

* `typename` - The node type.

* `cond_operation_status` - The flavor sales status.
  + **normal**: The flavor is in normal commercial use.
  + **sellout**: The flavor has been sold out.

* `localdisk` - Whether the node uses local disks.

* `edge` - Whether this is a node flavor for edge deployments.
