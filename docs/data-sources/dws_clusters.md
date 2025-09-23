---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_clusters"
description: |-
  Use this data source to get the list of DWS clusters within HuaweiCloud.
---

# huaweicloud_dws_clusters

Use this data source to get the list of DWS clusters within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dws_clusters" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `clusters` - All clusters that match the filter parameters.

  The [clusters](#clusters_struct) structure is documented below.

<a name="clusters_struct"></a>
The `clusters` block supports:

* `id` - The ID of the cluster.

* `name` - The name of the cluster.

* `node_type` - The flavor of the cluster.

* `number_of_node` - The number of nodes of the cluster.

* `user_name` - Administrator username for logging in to the cluster.

* `vpc_id` - The ID of the VPC corresponding to the cluster.

* `subnet_id` - The subnet ID corresponding to the cluster.

* `security_group_id` - The security group ID corresponding to the cluster.

* `availability_zone` - The availability zone of the cluster.

* `enterprise_project_id` - The enterprise project ID.

* `port` - The service port of the cluster.

* `version` - The version of the cluster.

* `tags` - The key/value pairs to associate with the cluster.

* `public_ip` - The public IP information of the cluster.

  The [public_ip](#clusters_public_ip_struct) structure is documented below.

* `endpoints` - The private network connection information of the cluster.

  The [endpoints](#clusters_endpoints_struct) structure is documented below.

* `public_endpoints` - The public network connection information of the cluster.

  The [public_endpoints](#clusters_public_endpoints_struct) structure is documented below.

* `recent_event` - The number of recent events of the cluster.

* `nodes` - The instance information of the cluster.

  The [nodes](#clusters_nodes_struct) structure is documented below.

* `status` - The current status of the cluster.
  + **ACTIVE**
  + **AVAILABLE**
  + **FAILED**
  + **CREATE_FAILED**
  + **DELETE_FAILED**
  + **DELETED**
  + **FROZEN**

* `sub_status` - The sub-status of the available cluster state.
  + **NORMAL**
  + **READONLY**
  + **REDISTRIBUTING**
  + **REDISTRIBUTION-FAILURE**
  + **UNBALANCED**
  + **UNBALANCED | READONLY**
  + **DEGRADED**
  + **DEGRADED | READONLY**
  + **DEGRADED | UNBALANCED**
  + **UNBALANCED | REDISTRIBUTING**
  + **UNBALANCED | REDISTRIBUTION-FAILURE**
  + **READONLY | REDISTRIBUTION-FAILURE**
  + **UNBALANCED | READONLY | REDISTRIBUTION-FAILURE**
  + **DEGRADED | REDISTRIBUTION-FAILURE**
  + **DEGRADED | UNBALANCED | REDISTRIBUTION-FAILURE**
  + **DEGRADED | UNBALANCED | READONLY | REDISTRIBUTION-FAILURE**
  + **DEGRADED | UNBALANCED | READONLY**

* `task_status` - The management task status of the cluster.

* `created_at` - The creation time of the cluster, in RFC3339 format.

* `updated_at` - The latest update time of the cluster, in RFC3339 format.

<a name="clusters_public_ip_struct"></a>
The `public_ip` block supports:

* `public_bind_type` - The bind type of public IP.
  + **auto_assign**
  + **not_use**
  + **bind_existing**

* `eip_id` - The EIP ID.

<a name="clusters_endpoints_struct"></a>
The `endpoints` block supports:

* `connect_info` - The private network connection information.

* `jdbc_url` - The JDBC URL on the private network.

<a name="clusters_public_endpoints_struct"></a>
The `public_endpoints` block supports:

* `public_connect_info` - The public network connection information.

* `jdbc_url` - The JDBC URL of the public network.

<a name="clusters_nodes_struct"></a>
The `nodes` block supports:

* `id` - The ID of the cluster instance.

* `status` - The status of the cluster instance.
