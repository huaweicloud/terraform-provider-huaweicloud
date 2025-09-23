---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_cluster"
description: |-
  Manages a GaussDB(DWS) cluster resource within HuaweiCloud.  
---

# huaweicloud_dws_cluster

Manages a GaussDB(DWS) cluster resource within HuaweiCloud.  

## Example Usage

```hcl
variable "availability_zone" {}
variable "dws_cluster_name" {}
variable "dws_cluster_version" {}
variable "user_name" {}
variable "user_pwd" {}
variable "vpc_id" {}
variable "network_id" {}

resource "huaweicloud_networking_secgroup" "secgroup" {
  name        = "sg_dws"
  description = "terraform security group"
}

resource "huaweicloud_dws_cluster" "cluster" {
  name              = var.dws_cluster_name
  version           = var.dws_cluster_version
  node_type         = "dws.m3.xlarge"
  number_of_node    = 3
  number_of_cn      = 3
  availability_zone = var.availability_zone
  user_name         = var.user_name
  user_pwd          = var.user_pwd
  vpc_id            = var.vpc_id
  network_id        = var.network_id
  security_group_id = huaweicloud_networking_secgroup.secgroup.id

  volume {
    type     = "SSD"
    capacity = 300
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Cluster name, which must be unique and contains 4 to 64 characters, which
  consist of letters, digits, hyphens(-), or underscores(_) only and must start with a letter.
  Changing this creates a new cluster resource.

* `node_type` - (Required, String, ForceNew) The flavor of the cluster.  
 [For details](https://support.huaweicloud.com/intl/en-us/productdesc-dws/dws_01_00018.html).
  Changing this parameter will create a new resource.

* `number_of_node` - (Required, Int) Number of nodes in a cluster.  
 The value ranges from 3 to 256 in cluster mode. The value of stream warehouse(stand-alone mode) is 1.

* `user_name` - (Required, String, ForceNew) Administrator username for logging in to a data warehouse cluster.  
 The administrator username must: Consist of lowercase letters, digits, or underscores.
 Start with a lowercase letter or an underscore.
  Changing this parameter will create a new resource.

* `user_pwd` - (Required, String) Administrator password for logging in to a data warehouse cluster.
  A password contains 8 to 32 characters, which consist of letters, digits,
  and special characters(~!@#%^&*()-_=+|[{}];:,<.>/?).
  It cannot be the same as the username or the username written in reverse order.

* `vpc_id` - (Required, String, ForceNew) The VPC ID.
  Changing this parameter will create a new resource.

* `network_id` - (Required, String, ForceNew) The subnet ID.
  Changing this parameter will create a new resource.

* `security_group_id` - (Required, String) Specifies the security group ID of the cluster.
  Changing this parameter will create a new resource.

* `availability_zone` - (Required, String, ForceNew) The availability zone in which to create the cluster instance.
  If there are multiple available zones, separate by commas, e.g. **cn-north-4a,cn-north-4b,cn-north-4g**.
  Currently, multi-AZ clusters only support selecting `3` AZs. Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional, String) The enterprise project ID.

* `number_of_cn` - (Required, Int, ForceNew) The number of CN.  
  The value ranges from 2 to **number_of_node**, the maximum value is 20.
  This parameter must be used together with `version`.
  Changing this parameter will create a new resource.

* `version` - (Required, String, ForceNew) The cluster version.
  [For details](https://support.huaweicloud.com/intl/en-us/bulletin-dws/dws_12_0000.html).
  Changing this parameter will create a new resource.

* `volume` - (Optional, List, ForceNew) The information about the volume.
  Changing this parameter will create a new resource.
  For local disks, this parameter can not be specified.

  The [Volume](#DwsCluster_Volume) structure is documented below.

* `port` - (Optional, Int, ForceNew) Service port of a cluster (8000 to 10000). The default value is 8000.  
  Changing this parameter will create a new resource.

* `tags` - (Optional, Map) The key/value pairs to associate with the cluster.

* `dss_pool_id` - (Optional, String, ForceNew) Dedicated storage pool ID.
  Changing this parameter will create a new resource.

* `kms_key_id` - (Optional, String, ForceNew) The KMS key ID.
  Changing this parameter will create a new resource.

* `public_ip` - (Optional, List, ForceNew) The information about public IP.  

  Changing this parameter will create a new resource.

  The [PublicIp](#DwsCluster_PublicIp) structure is documented below.

* `keep_last_manual_snapshot` - (Optional, Int) The number of latest manual snapshots that need to be
  retained when deleting the cluster.

* `logical_cluster_enable` - (Optional, Bool) Specifies whether to enable logical cluster. The switch needs to be turned
  on before creating a logical cluster.

* `elb_id` - (Optional, String) Specifies the ID of the ELB load balancer.

* `lts_enable` - (Optional, Bool) Specifies whether to enable LTS. The default value is **false**.

* `description` - (Optional, String) Specifies the description of the cluster.

* `force_backup` - (Optional, Bool) Specified whether to automatically execute snapshot when shrinking the number of nodes.
  The default value is **true**.
  This parameter is required and available only when scaling-in the `number_of_node` parameter value.

<a name="DwsCluster_PublicIp"></a>
The `PublicIp` block supports:

* `public_bind_type` - (Optional, String) The bind type of public IP.  
  The valid value are **auto_assign**, **not_use**, and **bind_existing**. Defaults to **not_use**.

* `eip_id` - (Optional, String) The EIP ID.  

<a name="DwsCluster_Volume"></a>
The `Volume` block supports:

* `type` - (Optional, String) The volume type. Value options are as follows:
  + **SSD**: Ultra-high I/O. The solid-state drive (SSD) is used.
  + **SAS**: High I/O. The SAS disk is used.
  + **SATA**: Common I/O. The SATA disk is used.

* `capacity` - (Optional, String) The capacity size, in GB.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `endpoints` - Private network connection information about the cluster.
  The [Endpoint](#DwsCluster_Endpoint) structure is documented below.

* `public_endpoints` - Public network connection information about the cluster.
  The [PublicEndpoint](#DwsCluster_PublicEndpoint) structure is documented below.

* `recent_event` - The recent event number.

* `status` - The cluster status.  
  The valid values are **CREATING**, **AVAILABLE**, **ACTIVE**, **FAILED**, **CREATE_FAILED**,
  **DELETING**, **DELETE_FAILED**, **DELETED**, and **FROZEN**.

* `sub_status` - Sub-status of clusters in the AVAILABLE state.  
  The value can be one of the following:
    + READONLY
    + REDISTRIBUTING
    + REDISTRIBUTION-FAILURE
    + UNBALANCED
    + UNBALANCED | READONLY
    + DEGRADED
    + DEGRADED | READONLY
    + DEGRADED | UNBALANCED
    + UNBALANCED | REDISTRIBUTING
    + UNBALANCED | REDISTRIBUTION-FAILURE
    + READONLY | REDISTRIBUTION-FAILURE
    + UNBALANCED | READONLY | REDISTRIBUTION-FAILURE
    + DEGRADED | REDISTRIBUTION-FAILURE
    + DEGRADED | UNBALANCED | REDISTRIBUTION-FAILURE
    + DEGRADED | UNBALANCED | READONLY | REDISTRIBUTION-FAILURE
    + DEGRADED | UNBALANCED | READONLY

* `task_status` - Cluster management task.  
  The value can be one of the following:
    + UNFREEZING
    + FREEZING
    + RESTORING
    + SNAPSHOTTING
    + GROWING
    + REBOOTING
    + SETTING_CONFIGURATION
    + CONFIGURING_EXT_DATASOURCE
    + DELETING_EXT_DATASOURCE
    + REBOOT_FAILURE
    + RESIZE_FAILURE

* `created` - The creation time of the cluster.  
  Format: ISO8601: **YYYY-MM-DDThh:mm:ssZ**.

* `updated` - The updated time of the cluster.  
  Format: ISO8601: **YYYY-MM-DDThh:mm:ssZ**.

* `private_ip` - List of private network IP addresses.  

* `maintain_window` - Cluster maintenance window.
  The [MaintainWindow](#DwsCluster_MaintainWindow) structure is documented below.

* `elb` - The ELB information bound to the cluster.
  The [elb](#DwsCluster_elb) structure is documented below.

<a name="DwsCluster_Endpoint"></a>
The `Endpoint` block supports:

* `connect_info` - Private network connection information.

* `jdbc_url` - JDBC URL. Format: jdbc:postgresql://<connect_info>/<YOUR_DATABASE_NAME>  

<a name="DwsCluster_PublicEndpoint"></a>
The `PublicEndpoint` block supports:

* `public_connect_info` - Public network connection information.

* `jdbc_url` - JDBC URL. Format: jdbc:postgresql://<public_connect_info>/<YOUR_DATABASE_NAME>  

<a name="DwsCluster_MaintainWindow"></a>
The `MaintainWindow` block supports:

* `day` - Maintenance time in each week in the unit of day.  
  The valid values are **Mon**, **Tue**, **Wed**, **Thu**, **Fri**,
  **Sat**, and **Sun**.

* `start_time` - Maintenance start time in HH:mm format. The time zone is GMT+0.

* `end_time` - Maintenance end time in HH:mm format. The time zone is GMT+0.

<a name="DwsCluster_elb"></a>
The `elb` block supports:

* `name` - The name of the ELB load balancer.

* `id` - The ID of the ELB load balancer.

* `public_ip` - The public IP address of the ELB load balancer.

* `private_ip` - The private IP address of the ELB load balancer.

* `private_endpoint` - The private endpoint of the ELB load balancer.

* `vpc_id` - The ID of VPC to which the ELB load balancer belongs.

* `private_ip_v6` - The IPv6 address of the ELB load balancer.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `update` - Default is 60 minutes.
* `delete` - Default is 60 minutes.

## Import

The resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dws_cluster.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `user_pwd`, `number_of_cn`, `kms_key_id`,
`volume`, `dss_pool_id`, `logical_cluster_enable`, `lts_enable`, `force_backup`.
It is generally recommended running `terraform plan` after importing a cluster.
You can then decide if changes should be applied to the cluster, or the resource definition
should be updated to align with the cluster. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dws_cluster" "test" {
  ...

  lifecycle {
    ignore_changes = [
      user_pwd, number_of_cn, kms_key_id, volume, dss_pool_id, logical_cluster_enable, lts_enable, `force_backup`,
    ]
  }
}
```
