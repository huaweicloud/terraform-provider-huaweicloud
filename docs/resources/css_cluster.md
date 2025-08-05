---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_cluster"
description: |-
  Manages CSS cluster resource within HuaweiCloud
---

# huaweicloud_css_cluster

Manages CSS cluster resource within HuaweiCloud

## Example Usage

### create a cluster

```hcl
variable "availability_zone" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "secgroup_id" {}

resource "huaweicloud_css_cluster" "cluster" {
  name           = "terraform_test_cluster"
  engine_version = "7.10.2"

  ess_node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = 1
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  availability_zone = var.availability_zone
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.secgroup_id
}
```

### create a cluster with ess-data node and master node

```hcl
variable "availability_zone" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "secgroup_id" {}

resource "huaweicloud_css_cluster" "cluster" {
  name           = "terraform_test_cluster"
  engine_version = "7.10.2"

  ess_node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = 1
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  master_node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = 3
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  availability_zone = var.availability_zone
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.secgroup_id
}
```

### create a cluster with ess-data node and cold node use local disk

```hcl
variable "availability_zone" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "secgroup_id" {}

resource "huaweicloud_css_cluster" "cluster" {
  name           = "terraform_test_cluster"
  engine_version = "7.10.2"

  ess_node_config {
    flavor          = "ess.spec-ds.xlarge.8"
    instance_number = 1
  }

  cold_node_config {
    flavor          = "ess.spec-ds.2xlarge.8"
    instance_number = 2
  }

  availability_zone = var.availability_zone
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.secgroup_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the cluster resource. If omitted, the
  provider-level region will be used. Changing this creates a new cluster resource.

* `name` - (Required, String, ForceNew) Specifies the cluster name. It contains `4` to `32` characters.
  Only letters, digits, hyphens (-), and underscores (_) are allowed. The value must start with a letter.
  Changing this parameter will create a new resource.

* `engine_type` - (Optional, String, ForceNew) Specifies the engine type. The valid value can be **elasticsearch** or **opensearch**.
  Defaults to **elasticsearch**. Changing this parameter will create a new resource.

* `engine_version` - (Required, String, NonUpdatable) Specifies the engine version.
  For details, see [Supported Cluster Versions](https://support.huaweicloud.com/intl/en-us/api-css/css_03_0056.html).

* `security_mode` - (Optional, Bool) Specifies whether to enable authentication.
  The value can be **true** or **false**. Authentication is disabled by default.
  + **true:** Authentication is enabled for the cluster.
  + **false:** Authentication is disabled for the cluster.

  -> **Note:** This parameter is supported in clusters 6.5.4 or later. Exercise caution when performing this operation.
  The cluster will be restarted to apply the changes. Services will be interrupted for a while.
  After the cluster is restarted, the authentication mode for invoking the cluster changes.
  You need to adjust the authentication mode accordingly.

* `password` - (Optional, String) Specifies the password of the cluster administrator in security mode.
  This parameter is mandatory only when `security_mode` is set to true.
  The administrator password must meet the following requirements:
  + The password can contain 8 to 32 characters.
  + The password must contain at least 3 of the following character types: uppercase letters, lowercase letters, digits,
    and special characters (~!@#$%^&*()-_=+\\|[{}];:,<.>/?).

* `https_enabled` - (Optional, Bool, ForceNew) Specifies whether to enable HTTPS. Defaults to **false**.
  When `https_enabled` is set to **true**, the `security_mode` needs to be set to **true**.
  Changing this parameter will create a new resource.

* `ess_node_config` - (Required, List) Specifies the config of data node.
  The [ess_node_config](#Css_ess_node_config) structure is documented below.

* `master_node_config` - (Optional, List) Specifies the config of master node.
  The [master_node_config](#Css_ess_master_or_client_node_config) structure is documented below.

* `client_node_config` - (Optional, List) Specifies the config of client node.
  The [client_node_config](#Css_ess_master_or_client_node_config) structure is documented below.

* `cold_node_config` - (Optional, List) Specifies the config of cold data node.
  The [cold_node_config](#Css_ess_node_config) structure is documented below.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID.
  Changing this parameter will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the Subnet ID.
  Changing this parameter will create a new resource.

* `security_group_id` - (Required, String) Specifies the security group ID.

* `availability_zone` - (Required, String, NonUpdatable) Specifies the availability zone name.
  Separate multiple AZs with commas (,), for example, az1,az2. AZs must be unique. The number of nodes must be greater
  than or equal to the number of AZs. If the number of nodes is a multiple of the number of AZs, the nodes are evenly
  distributed to each AZ. If the number of nodes is not a multiple of the number of AZs, the absolute difference
  between node quantity in any two AZs is **1** at most.

* `backup_strategy` - (Optional, List) Specifies the advanced backup policy. Structure is documented below.

* `tags` - (Optional, Map) The key/value pairs to associate with the cluster.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project id of the css cluster, The value **0**
  indicates the default enterprise project.

* `public_access` - (Optional, List) Specifies the public network access information.
  The [public_access](#Css_public_access) structure is documented below.

* `vpcep_endpoint` - (Optional, List) Specifies the VPC endpoint service information.
  The [vpcep_endpoint](#Css_vpcep_endpoint) structure is documented below.

* `kibana_public_access` - (Optional, List) Specifies Kibana public network access information.
  This parameter is valid only when security_mode is set to **true**.
  The [kibana_public_access](#Css_kibana_public_access) structure is documented below.

* `charging_mode` - (Optional, String) Specifies the charging mode of the cluster.
  Valid values are **prePaid** and **postPaid**, defaults to **postPaid**.

* `period_unit` - (Optional, String) Specifies the charging period unit of the instance.
  Valid values are **month** and **year**.

* `period` - (Optional, Int) Specifies the charging period of the instance.
  If `period_unit` is set to **month**, the value ranges from `1` to `9`.
  If `period_unit` is set to **year**, the value ranges from `1` to `9`.

  -> **NOTE:** `charging_mode`, `period_unit`, `period` can only be updated when changing
  from **postPaid** to **prePaid** billing mode.

* `auto_renew` - (Optional, String) Specifies whether auto renew is enabled.
  The valid values are **true** and **false**, defaults to **false**.

<a name="Css_ess_node_config"></a>
The `ess_node_config` and `cold_node_config` block supports:

* `flavor` - (Required, String) Specifies the flavor name.

* `instance_number` - (Required, Int) Specifies the number of cluster instances.
  + When it is `ess_node_config`, The value range is `1` to `200`.
  + When it is `cold_node_config`, The value range is `1` to `32`.

* `type` - (Optional, String) Specifies the instance type.
  The valid values are as follows:
  + **ess**
  + **chinese**
  + **english**
  + **arabic**
  + **tools**
  + **thai**
  + **turkish**
  + **portuguese**
  + **chinese-english**
  + **spanish**

* `volume` - (Optional, List, ForceNew) Specifies the information about the volume. This field should not be specified
  when `flavor` is set to a local dist flavor. But It is required when `flavor` is not a local disk flavor.
  Currently, the following local disk flavors are supported:
  + **ess.spec-i3small**
  + **ess.spec-i3medium**
  + **ess.spec-i3.8xlarge.8**
  + **ess.spec-ds.xlarge.8**
  + **ess.spec-ds.2xlarge.8**
  + **ess.spec-ds.4xlarge.8**

  The [volume](#Css_volume) structure is documented below. Changing this parameter will create a new resource.

* `shrink_node_ids` - (Optional, List) Specifies the node IDs that needs to be scaled down.

<a name="Css_volume"></a>
The `volume` block supports:

* `size` - (Required, Int) Specifies the volume size in **GB**, which must be a multiple of `10`.

* `volume_type` - (Required, String, ForceNew) Specifies the volume type. Value options are as follows:
  + **COMMON:** Common I/O. The SATA disk is used.
  + **HIGH:** High I/O. The SAS disk is used.
  + **ULTRAHIGH:** Ultra-high I/O. The solid-state drive (SSD) is used.

  Changing this parameter will create a new resource.

<a name="Css_ess_master_or_client_node_config"></a>
The `master_node_config` and `client_node_config` block supports:

* `flavor` - (Required, String) Specifies the flavor name.

* `instance_number` - (Required, Int) Specifies the number of cluster instances.
  + When it is `master_node_config`, The value range is `3` to `10`.
  + When it is `client_node_config`, The value range is `1` to `32`.

* `volume` - (Required, List) Specifies the information about the volume.
  The [volume](#Css_master_or_client_volume) structure is documented below.

* `shrink_node_ids` - (Optional, List) Specifies the node IDs that needs to be scaled down.

<a name="Css_master_or_client_volume"></a>
The `volume` block supports:

* `size` - (Required, Int) Specifies the volume size in **GB**, which must be a multiple of `10`.

* `volume_type` - (Required, String) Specifies the volume type. Value options are as follows:
  + **COMMON**: Common I/O. The SATA disk is used.
  + **HIGH**: High I/O. The SAS disk is used.
  + **ULTRAHIGH**: Ultra-high I/O. The solid-state drive (SSD) is used.

<a name="Css_public_access"></a>
The `public_access` block supports:

* `bandwidth` - (Required, Int) Specifies the public network bandwidth.

* `whitelist_enabled` - (Required, Bool) Specifies whether to enable the Kibana access control.

* `whitelist` - (Optional, String) Specifies the whitelist of Kibana access control.
  Separate the whitelisted network segments or IP addresses with commas (,), and each of them must be unique.

<a name="Css_kibana_public_access"></a>
The `kibana_public_access` block supports:

* `bandwidth` - (Required, Int) Specifies the public network bandwidth.

* `whitelist_enabled` - (Required, Bool) Specifies whether to enable the public network access control.

* `whitelist` - (Required, String) Specifies the whitelist of public network access control.
  Separate the whitelisted network segments or IP addresses with commas (,), and each of them must be unique.

<a name="Css_vpcep_endpoint"></a>
The `vpcep_endpoint` block supports:

* `endpoint_with_dns_name` - (Required, Bool) Specifies whether to enable the private domain name.

* `whitelist` - (Optional, List) Specifies the whitelist of access control. The whitelisted account id must be unique.

The `backup_strategy` block supports:

* `start_time` - (Required, String) Specifies the time when a snapshot is automatically created everyday. Snapshots can
  only be created on the hour. The time format is the time followed by the time zone, specifically, **HH:mm z**. In the
  format, **HH:mm** refers to the hour time and z refers to the time zone. For example, "00:00 GMT+08:00"
  and "01:00 GMT+08:00".

* `keep_days` - (Optional, Int) Specifies the number of days to retain the generated snapshots. Snapshots are reserved
  for seven days by default.

* `prefix` - (Optional, String) Specifies the prefix of the snapshot that is automatically created. Defaults to **snapshot**.

* `bucket` - (Optional, String) Specifies the OBS bucket used for index data backup. If there is snapshot data in an OBS
   bucket, only the OBS bucket is used and cannot be changed.

* `backup_path` - (Optional, String) Specifies the storage path of the snapshot in the OBS bucket.

* `agency` - (Optional, String) Specifies the IAM agency used to access OBS.

  -> **NOTE:**  If the `bucket`, `backup_path`, and `agency` parameters are empty at the same time, the system will
  automatically create an OBS bucket and IAM agent, otherwise the configured parameter values will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `endpoint` - The IP address and port number.

* `status` - The cluster status
  + **100:** The operation, such as instance creation, is in progress.
  + **200:** The cluster is available.
  + **303:** The cluster is unavailable.

* `nodes` - List of node objects.
  The [nodes](#Css_nodes_attr) structure is documented below.

* `vpcep_endpoint_id` - The VPC endpoint service ID.

* `vpcep_ip` - The private IP address of VPC endpoint service.

* `public_access/public_ip` - The public IP address.

* `kibana_public_access/public_ip` - The Kibana public IP address.

* `created_at` - Time when a cluster is created. The format is ISO8601: CCYY-MM-DDThh:mm:ss.

* `updated_at` - Time when a cluster is updated. The format is ISO8601: CCYY-MM-DDThh:mm:ss.

* `bandwidth_resource_id` - The resource ID of bandwidth.

* `is_period` - Whether a cluster is billed on the yearly/monthly mode.

* `backup_available` - Whether the snapshot function is enabled.

* `disk_encrypted` - Whether disks are encrypted.

<a name="Css_nodes_attr"></a>
The `nodes` block supports:

* `id` - Instance ID.

* `name` - Instance name.

* `type` - Node type. The options are as follows:
  + **ess-master:** Indicates a master node.
  + **ess-client:** Indicates a client node.
  + **ess-cold:** Indicates a cold data node.
  + **ess indicates:** Indicates a data node.

* `availability_zone` - The availability zone where the instance resides.

* `status` - Instance status.

* `spec_code` - Instance specification code.

* `ip` - Instance IP address.

* `resource_id` - The resource ID of this instance.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `update` - Default is 60 minutes.
* `delete` - Default is 60 minutes.

## Import

The CSS cluster can be imported by `id`, e.g.

```bash
terraform import huaweicloud_css_cluster.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `ess_node_config.0.type`, `cold_node_config.0.type`.
It is generally recommended running `terraform plan` after importing a cluster.
You can then decide if changes should be applied to the cluster, or the resource definition should be updated to
align with the cluster. Also you can ignore changes as below.

```hcl
resource "huaweicloud_css_cluster" "test" {
  ...

  lifecycle {
    ignore_changes = [
      ess_node_config.0.type, cold_node_config.0.type,
    ]
  }
}
```
