---
subcategory: "MapReduce Service (MRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_mapreduce_clusters"
description: ""
---

# huaweicloud_mapreduce_clusters

Use this data source to get clusters of MapReduce.

## Example Usage

```hcl
data "huaweicloud_mapreduce_clusters" "test" {
  status = "running"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) The name of cluster.

* `status` - (Optional, String) The status of cluster.  
  The following options are supported:
    + **existing**: Query existing clusters, including all clusters except those in the deleted state
      and the yearly/monthly clusters in the Order processing or preparing state.
    + **history**: Quer historical clusters, including all the deleted clusters, clusters that fail to delete,
      clusters whose VMs fail to delete, and clusters whose database updates fail to delete.
    + **starting**: Query a list of clusters that are being started.
    + **running**: Query a list of running clusters.
    + **terminated**: Query a list of terminated clusters.
    + **failed**: Query a list of failed clusters.
    + **abnormal**: Query a list of abnormal clusters.
    + **terminating**: Query a list of clusters that are being terminated.
    + **frozen**: Query a list of frozen clusters.
    + **scaling-out**: Query a list of clusters that are being scaled out.
    + **scaling-in**: Query a list of clusters that are being scaled in.

* `enterprise_project_id` - (Optional, String) The enterprise project ID used to query clusters in a specified
  enterprise project.
  The default value is **0**, indicating the default enterprise project.

* `tags` - (Optional, String) You can search for a cluster by its tags.  
  If you specify multiple tags, the relationship between them is **AND**.
  The format of the tags parameter is **tags=k1\*v1,k2\*v2,k3\*v3**.
  When the values of some tags are null, the format is **tags=k1,k2,k3\*v3**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `clusters` - The list of clusters.
  The [clusters](#MrsClusters_Clusters) structure is documented below.

<a name="MrsClusters_Clusters"></a>
The `Clusters` block supports:

* `id` - Cluster ID.

* `name` - Cluster name.

* `master_node_num` - Number of Master nodes deployed in a cluster.

* `core_node_num` - Number of Core nodes deployed in a cluster.

* `total_node_num` - Total number of nodes deployed in a cluster.

* `status` - Cluster status.  
  The following options are supported:  
    + **starting**: The cluster is being started.
    + **running**: The cluster is running.
    + **terminated**: The cluster has been terminated.
    + **failed**: The cluster fails.
    + **abnormal**: The cluster is abnormal.
    + **terminating**: The cluster is being terminated.
    + **frozen**: The cluster has been frozen.
    + **scaling-out**: The cluster is being scaled out.
    + **scaling-in**: The cluster is being scaled in.

* `billing_type` - Cluster billing mode.  
  The following options are supported:  
    + **11**: Yearly/Monthly.
    + **12**: Pay-per-use.

* `vpc_id` - VPC ID.

* `subnet_id` - Subnet ID.

* `duration` - Cluster subscription duration.

* `fee` - Cluster creation fee, which is automatically calculated.

* `hadoop_version` - Hadoop version.

* `master_node_size` - Instance specifications of a Master node.

* `core_node_size` - Instance specifications of a Core node.

* `component_list` - Component list.
  The [component_list](#MrsClusters_ClustersComponentList) structure is documented below.

* `external_ip` - External IP address.

* `external_alternate_ip` - Backup external IP address.

* `internal_ip` - Internal IP address.

* `deployment_id` - Cluster deployment ID.

* `description` - Cluster description.

* `order_id` - Cluster creation order ID.

* `master_node_product_id` - Product ID of a Master node.

* `master_node_spec_id` - Specification ID of a Master node.

* `core_node_product_id` - Product ID of a Core node.

* `core_node_spec_id` - Specification ID of a Core node.

* `availability_zone` - The AZ.

* `vnc` - URI for remotely logging in to an ECS.

* `volume_size` - Disk storage space.

* `volume_type` - Disk type.

* `enterprise_project_id` - Enterprise project ID.

* `type` - Cluster type.  
  The following options are supported:  
    + **0**: analysis cluster.
    + **1**: streaming cluster.
    + **2**: hybrid cluster.
    + **3**: custom cluster.
    + **4**: Offline cluster.

* `security_group_id` - Security group ID.

* `slave_security_group_id` - Security group ID of a non-Master node.  
  Currently, one MRS cluster uses only one security group. Therefore, this field has been discarded.

* `stage_desc` - Cluster progress description.  
  The cluster installation progress includes:
    + **Verifying cluster parameters**: Cluster parameters are being verified.
    + **Applying for cluster resources**: Cluster resources are being applied for.
    + **Creating VMs**: The VMs are being created.
    + **Initializing VMs**: The VMs are being initialized.
    + **Installing MRS Manager**: MRS Manager is being installed.
    + **Deploying the cluster**: The cluster is being deployed.
    + **Cluster installation failed**: Failed to install the cluster.

  The cluster scale-out progress includes:
    + **Preparing for scale-out**: Cluster scale-out is being prepared.
    + **Creating VMs**: The VMs are being created.
    + **Initializing VMs**: The VMs are being initialized.
    + **Adding nodes to the cluster**: The nodes are being added to the cluster.
    + **Scale-out failed**: Failed to scale out the cluster.

  The cluster scale-in progress includes:
    + **Preparing for scale-in**: Cluster scale-in is being prepared.
    + **Decommissioning instance**: The instance is being decommissioned.
    + **Deleting VMs**: The VMs are being deleted.
    + **Deleting nodes from the cluster**: The nodes are being deleted from the cluster.
    + **Scale-in failed**: Failed to scale in the cluster.

  If the cluster installation, scale-out, or scale-in fails, stageDesc will display the failure cause.

* `safe_mode` - Running mode of an MRS cluster.  
  The following options are supported:
    + **0**: Normal cluster.
    + **1**: Security cluster.

* `version` - Cluster version.

* `node_public_cert_name` - Name of the key file.

* `master_node_ip` - IP address of a Master node.

* `private_ip_first` - Preferred private IP address.

* `tags` - The tag information.

* `log_collection` - Whether to collect logs when cluster installation fails.  
  The following options are supported:
    + **0**: Do not collect logs.
    + **1**: Collect logs.

* `task_node_groups` - List of Task nodes.
  The [NodeGroup](#MrsClusters_ClustersNodeGroup) structure is documented below.

* `node_groups` - List of Master, Core and Task nodes.
  The [NodeGroup](#MrsClusters_ClustersNodeGroup) structure is documented below.

* `master_data_volume_type` - Data disk storage type of the Master node.  
  Currently, **SATA**, **SAS**, and **SSD** are supported.

* `master_data_volume_size` - Data disk storage space of the Master node  
  To increase data storage capacity, you can add disks at the same time when creating a cluster.  
  Value range: 100 GB to 32,000 GB

* `master_data_volume_count` - Number of data disks of the Master node  
  The value can be set to 1 only.

* `core_data_volume_type` - Data disk storage type of the Core node.  
  Currently, **SATA**, **SAS**, and **SSD** are supported.

* `core_data_volume_size` - Data disk storage space of the Core node.  
  To increase data storage capacity, you can add disks at the same time when creating a cluster.  
  Value range: 100 GB to 32,000 GB

* `core_data_volume_count` - Number of data disks of the Core node.

* `period_type` - Whether the subscription type is yearly or monthly.  
  The following options are supported:  
    + **0**: monthly subscription.
    + **1**: yearly subscription.

* `scale` - Status of node changes  
  If this parameter is left blank, no change operation is performed on a cluster node.  
  The options are as follows:
    + **Scaling-out**: The cluster is being scaled out.
    + **Scaling-in**: The cluster is being scaled in.
    + **scaling-error**: The cluster is in the running state and fails to be scaled in or out or the specifications
      fail to be scaled up for the last time.
    + **scaling-up**: The master node specifications are being scaled up.
    + **scaling_up_first**: The standby master node specifications are being scaled up.
    + **scaled_up_first**: The standby master node specifications have been scaled up.
    + **scaled-up-success**: The master node specifications have been scaled up.

* `eip_id` - Unique ID of the cluster EIP.

* `eip_address` - IPv4 address of the cluster EIP.

* `eipv6_address` - IPv6 address of the cluster EIP.  
  This parameter is not returned when an IPv4 address is used.

* `mrs_ecs_default_agency` - The default agency name bound to the cluster node.

<a name="MrsClusters_ClustersComponentList"></a>
The `component_list` block supports:

* `component_id` - Component ID  
  For example, the component_id of Hadoop is MRS 3.0.2_001, MRS 2.1.0_001, MRS 1.9.2_001, MRS 1.8.10_001.

* `component_name` - Component name.

* `component_version` - Component version.

* `component_desc` - Component description.

<a name="MrsClusters_ClustersNodeGroup"></a>
The `NodeGroup` block supports:

* `group_name` - Node group name.

* `node_num` - Number of nodes in a node group.

* `node_size` - Instance specifications of a node group.

* `node_spec_id` - Instance specification ID of a node group.

* `node_product_id` - Instance product ID of a node group.

* `vm_product_id` - VM product ID of a node group.

* `vm_spec_code` - VM specification code of a node group.

* `root_volume_size` - Root disk storage space of a node group.

* `root_volume_type` - Root disk storage type of a node group.

* `root_volume_product_id` - Root disk product ID of a node group.

* `root_volume_resource_spec_code` - Root disk specification code of a node group.

* `root_volume_resource_type` - System disk product type of a node group.

* `data_volume_type` - Data disk storage type of a node group.  
  The following options are supported:  
    + **SATA**: Common I/O.
    + **SAS**: High I/O.
    + **SSD**: Ultra-high I/O.

* `data_volume_count` - Number of data disks of a node group.

* `data_volume_size` - Data disk storage space of a node group.

* `data_volume_product_id` - Data disk product ID of a node group.

* `data_volume_resource_spec_code` - Data disk specification code of a node group.

* `data_volume_resource_type` - Data disk product type of a node group.
