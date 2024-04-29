---
subcategory: "Cloud Data Migration (CDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdm_clusters"
description: ""
---

# huaweicloud_cdm_clusters

Use this data source to get clusters of CDM.

## Example Usage

```hcl
variable "cluster_name" {}

data "huaweicloud_cdm_clusters" "test" {
  name = var.cluster_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Cluster name.

* `availability_zone` - (Optional, String) The AZ name.  

* `status` - (Optional, String) Cluster status.  
  Value options are as follows:
    + **100**: creating.
    + **200**: normal.
    + **300**: failed.
    + **303**: failed to be created.
    + **500**: restarting.
    + **800**: frozen.
    + **900**: stopped.
    + **910**: stopping.
    + **920**: starting.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `clusters` - The list of clusters.
  The [clusters](#CdmClusters_Cluster) structure is documented below.

<a name="CdmClusters_Cluster"></a>
The `clusters` block supports:

* `id` - Cluster ID.

* `name` - Cluster name.

* `availability_zone` - The AZ name.  

* `version` - Cluster version.  

* `is_auto_off` - Whether auto shutdown is enabled.

* `status` - Cluster status.
  Value options are as follows:
    + **100**: creating.
    + **200**: normal.
    + **300**: failed.
    + **303**: failed to be created.
    + **500**: restarting.
    + **800**: frozen.
    + **900**: stopped.
    + **910**: stopping.
    + **920**: starting.

* `recent_event` - Number of events.  

* `public_endpoint` - EIP bound to the cluster.  

* `is_frozen` - Whether the cluster is frozen. The value can be 0 (not frozen) or 1 (frozen).  

* `is_failure_remind` - Whether to notifications when a table/file migration job fails or an EIP exception occurs.  

* `instances` - The list of instance nodes.
  The [Instance](#CdmClusters_ClusterInstance) structure is documented below.

<a name="CdmClusters_ClusterInstance"></a>
The `ClusterInstance` block supports:

* `id` - Instance ID.  

* `name` - Instance name.

* `private_ip` - Private IP.

* `public_ip` - Public IP.

* `manage_ip` - Management IP address.

* `traffic_ip` - Traffic IP.

* `role` - Instance role.

* `type` - Instance type.

* `is_frozen` - Whether the node is frozen. The value can be 0 (not frozen) or 1 (frozen).  

* `status` - Instance node status.
  Value options are as follows:
    + **100**: creating.
    + **200**: normal.
    + **300**: failed.
    + **303**: failed to be created.
    + **400**: deleted.
    + **800**: frozen.
