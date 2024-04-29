---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_clusters"
description: |-
  Use this data source to get the list of CSS clusters.
---

# huaweicloud_css_clusters

Use this data source to get the list of CSS clusters.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_css_clusters" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Optional, String) Specifies the cluster ID.

* `name` - (Optional, String) Specifies the cluster name.

* `engine_type` - (Optional, String) Specifies the engine type. The values can be **elasticsearch** and **logstash**.

* `engine_version` - (Optional, String) Specifies the engine version.
  [For details](https://support.huaweicloud.com/intl/en-us/bulletin-css/css_05_0001.html)

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `clusters` - The list of cluster objects.

  The [clusters](#clusters_struct) structure is documented below.

<a name="clusters_struct"></a>
The `clusters` block supports:

* `id` - The cluster ID.

* `name` - The cluster name.

* `security_group_id` - The security group ID.

* `bandwidth_size` - The public network bandwidth. The unit is Mbit/s.

* `actions` - The current behavior of a cluster.
  + **REBOOTING** indicates that the cluster is being restarted.
  + **GROWING** indicates that the cluster is being scaled.
  + **RESTORING** indicates that the cluster is being restored.
  + **SNAPSHOTTING** indicates that a snapshot is being created.

* `period` - Whether a cluster is billed on the yearly/monthly mode.
  + **true**: The cluster is billed on the yearly/monthly mode.
  + **false**: The cluster is billed on the pay-per-use mode.

* `instances` - The list of node objects.

  The [instances](#clusters_instances_struct) structure is documented below.

* `public_ip` - The public IP address information.

* `status` - The cluster status.
  + **100**: Creating.
  + **200**: Available.
  + **303**: Unavailable, for example, due to a creation failure.

* `subnet_id` - The subnet ID.

* `backup_available` - Whether the snapshot function is enabled.
  + **true**: The snapshot function is enabled.
  + **false**: The snapshot function is disabled.

* `enterprise_project_id` - The ID of the enterprise project that a cluster belongs to.
  If the user of the cluster does not enable the enterprise project,
  the setting of this parameter is not returned.

* `public_kibana_resp` - The kibana public network access information.

  The [public_kibana_resp](#clusters_public_kibana_resp_struct) structure is documented below.

* `vpc_id` - The ID of a VPC.

* `datastore` - The cluster data store.

  The [datastore](#clusters_datastore_struct) structure is documented below.

* `endpoint` - The IP address and port number of the user used to access the VPC.

* `https_enable` - The communication encryption status.
  + **false**: Communication encryption is not enabled.
  + **true**: Communication encryption is enabled.

* `authority_enable` - Whether to enable authentication.
  + **true**: Authentication is enabled for the cluster.
  + **false**: Authentication is not enabled for the cluster.

* `disk_encrypted` - Whether disks are encrypted.
  + **true**: Disks are encrypted.
  + **false**: Disks are not encrypted.

* `elb_white_list` - The EIP whitelist.

  The [elb_white_list](#clusters_elb_white_list_struct) structure is documented below.

* `updated_at` - The last modification time of a cluster.

* `created_at` - The cluster creation time.
  The returned cluster list is sorted by creation time in descending order.
  The latest cluster is displayed at the top.

* `bandwidth_resource_id` - The resource id for ES public network access.

<a name="clusters_instances_struct"></a>
The `instances` block supports:

* `spec_code` - The node specifications.

* `az_code` - The AZ of a node.

* `ip` - The instance IP address.

* `volume` - The instance volume.

  The [volume](#instances_volume_struct) structure is documented below.

* `status` - The node status.
  + **100**: Creating.
  + **200**: Available.
  + **303**: Unavailable, for example, due to a creation failure.

* `type` - The type of the current node.

* `id` - The cluster instance ID.

* `name` - The cluster instance name.

<a name="instances_volume_struct"></a>
The `volume` block supports:

* `type` - The instance volume type.

* `size` - The instance volume size.

<a name="clusters_public_kibana_resp_struct"></a>
The `public_kibana_resp` block supports:

* `eip_size` - The bandwidth range. The unit is Mbit/s.

* `elb_white_list_resp` - The elb white list of the cluster public kibana.

  The [elb_white_list_resp](#public_kibana_resp_elb_white_list_resp_struct) structure is documented below.

* `public_kibana_ip` - The IP address for accessing kibana.

* `bandwidth_resource_id` - The resource id for ES public network access.

<a name="public_kibana_resp_elb_white_list_resp_struct"></a>
The `elb_white_list_resp` block supports:

* `enable_white_list` - Whether the kibana access control is enabled.
  + **true**: Access control is enabled.
  + **false**: Access control is disabled.

* `white_list` - Whitelist of public network for accessing kibana.

<a name="clusters_datastore_struct"></a>
The `datastore` block supports:

* `type` - The engine type.

* `version` - The version of the CSS cluster engine.

<a name="clusters_elb_white_list_struct"></a>
The `elb_white_list` block supports:

* `enable_white_list` - Whether the public network access control is enabled.
  + **true**: Public network access control is enabled.
  + **false**: Public network access control is disabled.

* `white_list` - Whitelist for public network access.
