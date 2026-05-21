---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_htap_starrocks_nodes"
description: |-
  Use this data source to query the list of TaurusDB HTAP StarRocks instance nodes within HuaweiCloud.
---

# huaweicloud_taurusdb_htap_starrocks_nodes

Use this data source to query the list of TaurusDB HTAP StarRocks instance nodes within HuaweiCloud.

## Example Usage

```hcl
variable "htap_instance_id" {}

data "huaweicloud_taurusdb_htap_starrocks_nodes" "test" {
  instance_id = var.htap_instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the HTAP StarRocks instance nodes.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the HTAP instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `node_list` - The list of HTAP StarRocks instance nodes.
  The [node_list](#taurusdb_htap_starrocks_node_list_attr) structure is documented below.

<a name="taurusdb_htap_starrocks_node_list_attr"></a>
The `node_list` block supports:

* `node_id` - The ID of the HTAP instance node.

* `node_name` - The name of the HTAP instance node.

* `role` - The role of the HTAP instance node.
  The valid values are as follows:
  + **fe-leader**: Frontend leader node.
  + **fe-follower**: Frontend follower node.
  + **fe-observer**: Frontend observer node.
  + **be**: Backend node.

* `status` - The status of the HTAP instance node.
  The valid values are as follows:
  + **createfail**: The node failed to be created.
  + **creating**: The node is being created.
  + **normal**: The node is running properly.
  + **abnormal**: The node is abnormal
